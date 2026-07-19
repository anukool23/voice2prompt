# Voice2Prompt

An open-source, on-device voice-dictation tool for macOS + Windows (a Wispr Flow
equivalent). Hold a hotkey, speak, and the transcribed text is pasted into whatever
field you're focused on — no network calls, no per-word API cost.

**Stack:** Go engine + whisper.cpp (STT sidecar) + Ollama (LLM cleanup, Phase 2) +
Wails/React UI (Phase 3). See
[`flow-clone-implementation-plan.md`](flow-clone-implementation-plan.md) for the full plan.

## Status

| Phase | What | State |
|---|---|---|
| **0** | Feasibility spike: hotkey → record → whisper → paste, measure latency | ✅ working |
| **1** | macOS Accessibility-API text injection (AX caret insert + clipboard fallback) | ✅ working |
| **2** | Local LLM cleanup via Ollama (punctuation, fillers, context-aware tone) | ✅ working |
| **3** | Wails settings/onboarding app (permissions, hotkey, language, dictionary) | 🟡 builds, needs live test |
| 4 | Windows port (UI Automation) | — (needs a Windows machine) |
| **5** | Command mode — editing keystrokes + LLM rewrite of selection | 🟡 builds, needs live test |
| **6** | Packaging: self-contained `.app`, menu-bar tray, launch-at-login | 🟡 builds, needs live test |

## Architecture

The engine is Go. Speech-to-text runs as a local `whisper-server` sidecar that the Go
app talks to over HTTP (localhost) — this keeps the model resident so each utterance
skips the model load, and mirrors how the Ollama LLM will run in Phase 2. No cgo linking
against whisper is required.

```
hotkey (golang.design/x/hotkey)
  → mic capture (malgo → 16 kHz mono WAV)
    → POST /inference        →  whisper-server (whisper.cpp, Metal)   [STT]
      → POST /api/chat       →  ollama (qwen2.5:3b)                    [cleanup]
        → inject: AX caret-insert  ─(fallback)→  clipboard + ⌘V
```

Both the STT and LLM engines are local HTTP sidecars on `127.0.0.1` — nothing leaves
the machine at runtime.

### LLM cleanup (Phase 2)

`internal/llm` sends the raw transcript to a local Ollama model with a strict
post-processing prompt: fix punctuation/capitalization, drop filler words, apply light
context-aware tone (casual for Slack/Discord/Messages, professional for Mail/Docs), and
**never** answer or act on the dictated content. Cleanup is optional — `--raw` skips it,
and if Ollama is unavailable the app degrades to raw transcription instead of failing.

```sh
# tune/verify cleanup headlessly (no mic needed)
./bin/voice2prompt --clean-test "um so like what time is the meeting tomorrow"

# pick a different model / disable cleanup
./bin/voice2prompt --llm-model llama3.2:3b
./bin/voice2prompt --raw
```

### Text injection (Phase 1)

`internal/inject` inserts the transcript with a two-tier strategy:

1. **Accessibility API** (`AXUIElementSetAttributeValue` on the focused element's
   `AXSelectedText`). Inserts at the caret, preserves the clipboard, and can read the
   focused field's current text (`FocusedValue()`) as context for the Phase 2 LLM.
2. **Clipboard + ⌘V fallback**, used automatically when AX insertion is rejected. It
   saves and restores the previous clipboard text (best-effort; non-text clipboard
   contents like images aren't preserved).

The CLI prints which path was used (`via accessibility` / `via clipboard`).

**Known limitations** (the "ongoing QA" the plan calls out): native Cocoa fields accept
AX insertion reliably; many browser web inputs and Electron apps (Slack, Notion, VS Code)
reject it and take the clipboard fallback. Both paths need Accessibility permission.

## Phase 0 — feasibility spike

### Prerequisites

- Go (`brew install go`)
- Xcode Command Line Tools (for cgo: `xcode-select --install`)
- `cmake` (`brew install cmake`) — builds the whisper.cpp sidecar
- A whisper ggml model in `models/` (default `ggml-base.en.bin`):

  ```sh
  mkdir -p models
  curl -L -o models/ggml-base.en.bin \
    https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-base.en.bin
  ```

### One-time setup

```sh
bash scripts/setup-whisper.sh   # clones + builds whisper.cpp with Metal (STT)
bash scripts/setup-ollama.sh    # installs Ollama + pulls qwen2.5:3b (LLM cleanup)
```

The Ollama step is optional — without it, the app runs with raw transcription (`--raw`).

### Run

```sh
# interactive push-to-talk
go run ./cmd/voice2prompt
# or the built binary
go build -o bin/voice2prompt ./cmd/voice2prompt && ./bin/voice2prompt

# headless latency benchmark against a 16 kHz mono WAV
go run ./cmd/voice2prompt --bench path/to/audio.wav

# custom model
go run ./cmd/voice2prompt models/ggml-small.en.bin
```

Then **hold `Ctrl+Option+Space`**, speak, and **release**. The transcript is pasted at
your cursor and a timing line is printed:

```
📝 "hello world this is a test"
   audio 2.1s | infer 190ms | TOTAL 260ms | via accessibility  ✅ under 800ms
```

### macOS permissions

Grant to your **terminal app** (or, once packaged, to Voice2Prompt.app):

- **Accessibility** (System Settings → Privacy & Security → Accessibility): required for
  the simulated ⌘V paste. **Restart the terminal after granting** — TCC changes only take
  effect on app relaunch.
- **Microphone**: the OS prompts on first capture.

The hotkey needs **no** permission: `golang.design/x/hotkey` is pinned to **v0.4.0**,
which uses Carbon's `RegisterEventHotKey` (later versions switched to a `CGEventTap`,
which would require Input Monitoring). Trade-off: Carbon can't bind modifier-only
"hold `fn` to talk" chords — if we want that UX later we'll revisit the event-tap path
behind the Phase 3 onboarding flow.

Phase 3 builds a proper onboarding flow around the Accessibility/Microphone prompts.

## Phase 3 — desktop app (Wails)

A native settings/onboarding app that runs the dictation engine in-process. The CLI
and app share `internal/engine`; both read/write the same settings via `internal/config`
(JSON at `~/Library/Application Support/Voice2Prompt/config.json`).

Features: engine start/stop, Accessibility permission status + request, push-to-talk
hotkey picker, language selection, whisper/LLM model config, a custom dictionary editor,
first-run onboarding, and a live transcript feed.

```sh
cd desktop
wails dev      # hot-reload dev mode
wails build    # produces build/bin/Voice2Prompt.app
open build/bin/Voice2Prompt.app
```

## Triggers: hotkey chord or Fn key

Settings → Trigger picks how dictation activates:

- **Fn / 🌐 key** (default) — Wispr-Flow-style:
  - **Hold Fn** → push-to-talk (record while held, transcribe on release)
  - **Double-tap Fn** → lock into hands-free continuous recording; **tap Fn again** to stop
  - Implemented with a listen-only `CGEventTap` on the Fn modifier flag
    ([internal/hotkey/fn_darwin.go](internal/hotkey/fn_darwin.go)), so it needs **Input
    Monitoring** permission. Also set System Settings → Keyboard → "Press 🌐 key to" =
    **Do Nothing** so macOS doesn't open the emoji picker.
- **Hotkey chord** — e.g. `Ctrl+Option+Space`, hold to talk. Uses Carbon
  `RegisterEventHotKey`; needs no extra permission. (The CLI always uses this.)

## Phase 5 — command mode

When enabled (Settings → Voice commands), an utterance that matches a command is
executed instead of typed. `internal/command` does fast, rule-based intent parsing
(with guards so ordinary dictation like "select all the files" isn't misfired); the
engine executes it.

- **Editing** (instant, via CGEvent keystrokes): "select all", "copy/cut/paste",
  "undo" / "scratch that", "redo", "new line", "new paragraph", "delete last word",
  "delete line".
- **Rewrite** (via the LLM on your current selection): "make this formal", "make it
  casual", "shorten this", "make it longer", "fix the grammar", "make it bullets". Reads
  the selection (AX `kAXSelectedText`), rewrites it, and replaces it in place.

Everything else is dictated normally. Rewrite needs the cleanup LLM enabled and a
non-empty text selection.

## Phase 6 — packaging (macOS)

`scripts/package-mac.sh` builds a **self-contained** `Voice2Prompt.app`:

```sh
bash scripts/package-mac.sh
open desktop/build/bin/Voice2Prompt.app
```

- **Bundled engine** — `whisper-server`, its dylibs, and the whisper model are copied into
  `Contents/Resources`, with the dylib rpath rewritten to `@loader_path`, so the app runs
  from anywhere (no repo-relative paths). `internal/resource` resolves paths against the
  bundle at runtime.
- **Menu-bar tray** — `internal/tray` adds an `NSStatusItem` (Open / Start-Stop / Quit) via
  cgo/Objective-C on the main queue, coexisting with the Wails run loop. The window
  hides-on-close so the app lives in the menu bar.
- **Launch at login** — `internal/autostart` installs a per-user LaunchAgent plist (no
  Automation permission needed); toggled from Settings.
- The build is **ad-hoc signed**. Developer ID signing + notarization (for distribution
  without Gatekeeper warnings) still needs an Apple Developer account.

**Design notes**
- **Single process, one set of permissions** — the app hosts the engine directly, so you
  grant Accessibility to Voice2Prompt.app once (no child-process permission mess).
- **Hotkey/main-thread coexistence** — Wails owns the macOS main run loop, so the app
  can't use `golang.design/x/hotkey`. Instead `internal/hotkey/carbon_darwin.go` installs
  a Carbon `RegisterEventHotKey` handler on the app's existing event target (no extra
  loop, no Input Monitoring permission).
- **Plain HTML/CSS/JS frontend** (no npm/bundler) — keeps the build dependency-free; swap
  in React/Vite later if desired.

## Layout

```
go.mod                              # module: voice2prompt (shared by CLI + app)
cmd/voice2prompt/main.go         # CLI: interactive + --bench + --clean-test
desktop/                            # Wails app (main.go, app.go, wails.json, frontend/)
internal/engine/engine.go           # orchestration: record→transcribe→clean→dict→inject
internal/config/                    # settings model + persistence + dictionary
internal/audio/capture.go           # malgo capture → 16 kHz mono WAV
internal/stt/whisper.go             # whisper-server process mgr + HTTP client
internal/llm/                       # Ollama cleanup client + prompt/few-shot
internal/inject/inject_darwin.go    # AX caret-insert + clipboard fallback (cgo)
internal/hotkey/                    # golang.design (CLI) + Carbon cgo (Wails app)
scripts/setup-whisper.sh            # build the whisper.cpp sidecar
scripts/setup-ollama.sh             # install Ollama + pull model
models/                             # whisper models (gitignored)
third_party/whisper.cpp/            # built sidecar (gitignored)
```
