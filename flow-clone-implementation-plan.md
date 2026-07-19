# PromptVoice — Implementation Plan

Target: Voice2Prompt for macOS + Windows, built entirely on open-source components.

## 1. Full Open-Source Stack

| Layer | Component | License | Notes |
|---|---|---|---|
| App shell / runtime | Tauri | MIT / Apache 2.0 | Rust core + WebView UI, ~10x lighter than Electron |
| Core language | Go | — | Hotkey, audio, OS text injection |
| Global hotkey | `global-hotkey` crate | MIT / Apache 2.0 | Cross-platform push-to-talk capture |
| Audio capture | `cpal` crate | MIT / Apache 2.0 | Cross-platform mic streaming |
| Text injection (macOS) | Accessibility API via `accessibility-rs` / `objc2` | MIT / Apache 2.0 | Reads focused element, injects text |
| Text injection (Windows) | UI Automation via `windows-rs` | MIT / Apache 2.0 | Same role on Windows |
| Speech-to-text | whisper.cpp (via `whisper-rs`) | MIT | 99+ languages, runs fully on-device |
| Speech-to-text (alt, faster) | Voxtral Realtime (Mistral, 4B) | Apache 2.0 | Lower latency/WER than Whisper, but only 13 languages — evaluate as a swap-in later |
| Text cleanup / formatting LLM | Phi-3.5 Mini or Qwen2.5 (1.5B–3B, quantized) | MIT / Apache 2.0 | Runs locally via `llama.cpp` / Ollama |
| LLM inference runtime | llama.cpp / Ollama | MIT | Serves the cleanup model locally |
| Settings/onboarding UI | React + TypeScript | MIT | Rendered inside Tauri's webview |
| Backend (accounts, sync, dictionary) | Node.js + TypeScript, PostgreSQL | MIT / PostgreSQL License | Only needed if you want cross-device sync |
| Auth (if backend used) | Ory Kratos or self-rolled JWT | Apache 2.0 | Self-hostable, no third-party dependency |
| CI/CD | GitHub Actions or self-hosted Woodpecker CI | Apache 2.0 | Build/sign/release pipeline |

**One honest caveat:** if you ever want to charge money, the payment processor (Stripe, Paddle, etc.) will not be open source — that's an unavoidable exception, not a gap in the plan. Everything else above is fully open.

## 2. Architecture (end-to-end loop)

Hotkey press → `cpal` streams mic audio → buffered chunks fed to whisper.cpp/whisper-rs in streaming mode → partial transcript → sent with active-app context (frontmost app name, window title, selected text if any) to the local LLM (Phi-3.5/Qwen) with a cleanup prompt → cleaned text streamed back → injected into the focused field via Accessibility API (macOS) or UI Automation (Windows).

Everything above runs **on-device** — no network call is required end-to-end, which also solves privacy positioning and avoids ongoing API costs.

## 3. Phased Build Plan

**Phase 0 — Feasibility spike (1–2 weeks)**
Bare-bones CLI: global hotkey → record → whisper.cpp transcribe → paste raw text via keyboard simulation (`enigo`). Goal: confirm end-to-end latency is under ~800ms on target hardware (CPU-only and GPU-assisted). This is the make-or-break checkpoint — if local Whisper is too slow on average hardware, decide now whether to fall back to a smaller/distilled model or accept GPU-only support first.

**Phase 1 — macOS native integration (3–4 weeks)**
Replace raw paste with proper Accessibility API integration: detect focused element, read surrounding context, inject text correctly across different app types (native fields, web inputs in browsers, Electron apps like Slack/Notion). This is historically the hardest and buggiest part — budget the most time here.

**Phase 2 — LLM cleanup layer (2–3 weeks)**
Wire in local LLM inference (llama.cpp/Ollama) with a cleanup/formatting prompt. Tune prompt + model choice for latency vs. quality trade-off. Add app-context-aware formatting (casual for chat apps, formal for email/docs).

**Phase 3 — Product surface (2–3 weeks)**
Tauri-based settings UI: onboarding, permissions flow (mic + accessibility access), custom dictionary, hotkey configuration, language selection.

**Phase 4 — Windows port (2–3 weeks)**
Swap the macOS-specific Rust modules (Accessibility API → UI Automation), reuse everything else (STT, LLM, UI). Should be materially faster than Phase 1 since the hard problems are already solved once.

**Phase 5 — Command mode + dictionary sync (3+ weeks)**
Voice commands ("delete last sentence", "make this formal") via intent classification on top of the transcript. Optional backend for cross-device dictionary/settings sync if you want it.

**Phase 6 — Hardening & packaging**
Code signing, notarization (macOS), auto-update via Tauri's built-in updater, telemetry (opt-in, self-hosted if you want to stay fully open-source end to end).

## 4. Key Risks

Local LLM + local STT running simultaneously is CPU/GPU-intensive — needs real testing on low-end hardware (a 2019 MacBook Air, not just an M-series dev machine) or you'll ship something great by developers, unusable for most users. Accessibility permission prompts on macOS are known to confuse users during onboarding — plan for a dedicated, well-tested permissions flow. Cross-app text injection edge cases (rich text editors, contenteditable web fields, terminal apps) will surface bugs continuously — budget ongoing QA time here rather than treating Phase 1 as "done."

## 5. Suggested Team Shape

One Rust engineer owning core engine + OS integration (this is the critical path). One engineer on ML integration (STT/LLM tuning, quantization, latency optimization). One engineer on frontend/product (Tauri UI, onboarding, settings). Total: roughly 4–5 months to a usable macOS v1, another 4–6 weeks to Windows parity.
