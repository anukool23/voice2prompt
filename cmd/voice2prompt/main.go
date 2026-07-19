// Voice2Prompt — CLI (Phases 0–2 engine; shares internal/engine with the Wails app).
//
// Push-to-talk voice dictation, fully on-device:
//
//	hold hotkey → record mic → whisper transcribe → LLM cleanup → dictionary → inject.
//
// STT (whisper.cpp) and cleanup (Ollama) run as local sidecars over loopback HTTP.
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.design/x/hotkey/mainthread"

	"voice2prompt/internal/config"
	"voice2prompt/internal/engine"
	"voice2prompt/internal/hotkey"
	"voice2prompt/internal/inject"
	"voice2prompt/internal/llm"
	"voice2prompt/internal/stt"
)

type args struct {
	benchWav  string
	cleanTest string
	raw       bool
	llmModel  string
	model     string
	modelSet  bool
	rawSet    bool
	llmSet    bool
}

func parseArgs() args {
	var a args
	rest := os.Args[1:]
	for i := 0; i < len(rest); i++ {
		switch rest[i] {
		case "--bench":
			if i+1 < len(rest) {
				a.benchWav = rest[i+1]
				i++
			}
		case "--clean-test":
			if i+1 < len(rest) {
				a.cleanTest = rest[i+1]
				i++
			}
		case "--raw":
			a.raw, a.rawSet = true, true
		case "--llm-model":
			if i+1 < len(rest) {
				a.llmModel, a.llmSet = rest[i+1], true
				i++
			}
		default:
			a.model, a.modelSet = rest[i], true
		}
	}
	return a
}

// applyOverrides layers CLI flags / env on top of the persisted config.
func applyOverrides(cfg config.Settings, a args) config.Settings {
	if env := os.Getenv("PROMPTVOICE_MODEL"); env != "" {
		cfg.WhisperModel = env
	}
	if a.modelSet {
		cfg.WhisperModel = a.model
	}
	if a.rawSet {
		cfg.CleanupEnabled = !a.raw
	}
	if a.llmSet {
		cfg.LLMModel = a.llmModel
	}
	return cfg
}

func main() {
	a := parseArgs()

	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: %v (using defaults)\n", err)
	}
	cfg = applyOverrides(cfg, a)

	// Headless cleanup test (no mic/whisper needed).
	if a.cleanTest != "" {
		if err := runCleanTest(cfg.LLMModel, a.cleanTest); err != nil {
			fatal(err)
		}
		return
	}

	// Headless STT latency benchmark.
	if a.benchWav != "" {
		srv, err := stt.Start(stt.Config{ModelPath: cfg.WhisperModel, Language: cfg.Language})
		if err != nil {
			fatal(err)
		}
		defer srv.Close()
		if err := runBench(srv, a.benchWav); err != nil {
			fatal(err)
		}
		return
	}

	runInteractive(cfg)
}

func runInteractive(cfg config.Settings) {
	fmt.Println("Voice2Prompt (Go) — on-device voice dictation")
	fmt.Printf("  stt     : %s (lang %s)\n", cfg.WhisperModel, cfg.Language)
	if cfg.CleanupEnabled {
		fmt.Printf("  cleanup : Ollama %s\n", cfg.LLMModel)
	} else {
		fmt.Println("  cleanup : disabled")
	}
	fmt.Print("  starting engine… ")

	eng := engine.New(cfg)
	if err := eng.Start(); err != nil {
		fmt.Println("failed")
		fatal(err)
	}
	eng.SetOnResult(printResult)
	fmt.Println("ready")

	// Clean up on Ctrl+C.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sig
		fmt.Println("\nshutting down…")
		eng.Close()
		os.Exit(0)
	}()

	// The hotkey event loop must run on the macOS main thread.
	mainthread.Init(func() { hotkeyLoop(eng) })
}

func hotkeyLoop(eng *engine.Engine) {
	die := func(err error) {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		eng.Close()
		os.Exit(1)
	}

	hk := hotkey.PushToTalk()
	if err := hk.Register(); err != nil {
		die(fmt.Errorf("failed to register hotkey %s: %w", hotkey.Chord, err))
	}
	defer hk.Unregister()

	if !inject.Trusted() {
		fmt.Println("\n⚠️  Accessibility permission not granted (needed to insert text).")
		fmt.Println("   System Settings → Privacy & Security → Accessibility → enable your terminal, then relaunch.")
		inject.PromptTrust()
	}

	fmt.Printf("\nReady. Hold  %s  to talk, release to transcribe & paste.  (Ctrl+C to quit)\n", hotkey.Chord)

	for {
		select {
		case <-hk.Keydown():
			if err := eng.StartCapture(); err != nil {
				fmt.Fprintf(os.Stderr, "record error: %v\n", err)
			} else if eng.Recording() {
				fmt.Println("🎙️  recording…")
			}
		case <-hk.Keyup():
			eng.StopCapture()
		}
	}
}

func printResult(r engine.Result) {
	if r.Err != nil {
		fmt.Fprintf(os.Stderr, "… %v\n", r.Err)
		return
	}
	if r.Raw == "" {
		fmt.Println("… no speech detected.")
		return
	}
	if r.Command != "" {
		detail := ""
		if r.Cleaned != "" && r.Cleaned != r.Raw {
			detail = fmt.Sprintf("  →  %q", r.Cleaned)
		}
		fmt.Printf("⌘ command: %s%s\n", r.Command, detail)
		fmt.Printf("   heard %q | %dms | %s\n", r.Raw, r.TotalMS, appLabel(r.App))
		return
	}
	if r.Cleaned != r.Raw && r.Cleaned != "" {
		fmt.Printf("📝 raw:     %q\n   cleaned: %q\n", r.Raw, r.Cleaned)
	} else {
		fmt.Printf("📝 %q\n", r.Raw)
	}
	clean := ""
	if r.CleanMS > 0 {
		clean = fmt.Sprintf(" | clean %dms", r.CleanMS)
	}
	fmt.Printf("   audio %.1fs | infer %dms%s | TOTAL %dms | %s via %s  %s\n",
		r.AudioSecs, r.InferMS, clean, r.TotalMS, appLabel(r.App), r.Method, budget(r.TotalMS))
}

// runCleanTest exercises the LLM cleanup layer headlessly for verification/tuning.
func runCleanTest(model, text string) error {
	fmt.Printf("Cleanup test — Ollama model %s\n", model)
	client, err := llm.Start(llm.Config{Model: model})
	if err != nil {
		return err
	}
	defer client.Close()

	for _, appName := range []string{"", "Slack", "Mail"} {
		res, err := client.Clean(text, appName)
		if err != nil {
			return err
		}
		ctx := appName
		if ctx == "" {
			ctx = "(neutral)"
		}
		fmt.Printf("\n  context: %-8s  %dms\n  raw:     %q\n  cleaned: %q\n",
			ctx, res.Latency.Milliseconds(), text, res.Text)
	}
	return nil
}

func runBench(srv *stt.Server, wavPath string) error {
	fmt.Printf("Benchmark: %s\n", wavPath)
	wav, err := os.ReadFile(wavPath)
	if err != nil {
		return fmt.Errorf("failed to read WAV: %w", err)
	}
	audioSecs := float64(len(wav)-44) / float64(16000*2)
	fmt.Printf("  audio length: %.2fs\n\n", audioSecs)

	const runs = 4
	var warm []time.Duration
	for i := 0; i < runs; i++ {
		res, err := srv.Transcribe(wav)
		if err != nil {
			return err
		}
		tag := "warm"
		extra := ""
		if i == 0 {
			tag = "cold"
			extra = fmt.Sprintf("  →  %q", res.Text)
		} else {
			warm = append(warm, res.Latency)
		}
		fmt.Printf("  run %d [%s]: %dms (rtf %.3f)%s\n",
			i, tag, res.Latency.Milliseconds(), res.Latency.Seconds()/audioSecs, extra)
	}
	if len(warm) > 0 {
		var sum time.Duration
		best := warm[0]
		for _, d := range warm {
			sum += d
			if d < best {
				best = d
			}
		}
		fmt.Printf("\n  warm round-trip: avg %dms, best %dms\n",
			(sum / time.Duration(len(warm))).Milliseconds(), best.Milliseconds())
	}
	return nil
}

func appLabel(name string) string {
	if name == "" {
		return "?"
	}
	return name
}

func budget(totalMS int64) string {
	if totalMS <= 800 {
		return "✅ under 800ms"
	}
	return "⚠️ over 800ms budget"
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}
