// Package engine is the dictation orchestration — record → transcribe → clean →
// dictionary → inject — decoupled from any particular hotkey source or UI. Both the
// CLI and the Wails app construct an Engine and drive it via StartCapture/StopCapture.
package engine

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"voice2prompt/internal/audio"
	"voice2prompt/internal/command"
	"voice2prompt/internal/config"
	"voice2prompt/internal/inject"
	"voice2prompt/internal/llm"
	"voice2prompt/internal/stt"
)

// Result is the outcome of one utterance, for display/logging.
type Result struct {
	Raw       string
	Cleaned   string // == Raw when cleanup is off/failed
	App       string
	Method    inject.Method
	Command   string // non-empty if the utterance was executed as a voice command
	AudioSecs float64
	InferMS   int64
	CleanMS   int64
	TotalMS   int64
	Err       error
}

// Engine holds the running STT/LLM engines and the current recording state.
type Engine struct {
	cfg config.Settings

	stt *stt.Server
	llm *llm.Client // nil when cleanup disabled/unavailable

	mu  sync.Mutex
	rec *audio.Recorder

	onResult func(Result)
}

// New creates an engine for the given settings (not yet started).
func New(cfg config.Settings) *Engine {
	return &Engine{cfg: cfg}
}

// SetOnResult registers a callback invoked (on a background goroutine) after each
// utterance is processed.
func (e *Engine) SetOnResult(fn func(Result)) { e.onResult = fn }

// CleanupActive reports whether the LLM cleanup layer is live.
func (e *Engine) CleanupActive() bool { return e.llm != nil }

// Start brings up the STT sidecar and (if enabled) the LLM cleanup client.
func (e *Engine) Start() error {
	srv, err := stt.Start(stt.Config{
		ModelPath: e.cfg.WhisperModel,
		Language:  e.cfg.Language,
	})
	if err != nil {
		return err
	}
	e.stt = srv

	if e.cfg.CleanupEnabled {
		client, err := llm.Start(llm.Config{Model: e.cfg.LLMModel})
		if err != nil {
			// Cleanup is optional — degrade to raw transcription.
			fmt.Fprintf(os.Stderr, "cleanup unavailable, using raw transcription: %v\n", err)
		} else {
			e.llm = client
		}
	}
	return nil
}

// StartCapture begins recording (call on hotkey down). Idempotent.
func (e *Engine) StartCapture() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.rec != nil {
		return nil
	}
	r, err := audio.Start()
	if err != nil {
		return err
	}
	e.rec = r
	return nil
}

// Recording reports whether a capture is in progress.
func (e *Engine) Recording() bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.rec != nil
}

// StopCapture ends recording (call on hotkey up) and processes the utterance on a
// background goroutine so it never blocks the caller's (UI/main) thread.
func (e *Engine) StopCapture() {
	e.mu.Lock()
	r := e.rec
	e.rec = nil
	e.mu.Unlock()
	if r == nil {
		return
	}
	go e.process(r)
}

func (e *Engine) process(rec *audio.Recorder) {
	t0 := time.Now()
	wav := rec.Stop()

	audioBytes := len(wav) - 44
	minBytes := audio.SampleRate * audio.BitsPerSample / 8 / 5 // 0.2s
	if audioBytes < minBytes {
		return
	}
	audioSecs := float64(audioBytes) / float64(audio.SampleRate*audio.Channels*audio.BitsPerSample/8)

	app := inject.FocusedApp()

	res := Result{App: app, AudioSecs: audioSecs}

	tr, err := e.stt.Transcribe(wav)
	if err != nil {
		res.Err = fmt.Errorf("transcription: %w", err)
		e.emit(res)
		return
	}
	res.Raw = tr.Text
	res.InferMS = tr.Latency.Milliseconds()
	if res.Raw == "" {
		e.emit(res) // no speech
		return
	}

	// Command mode: if the utterance is a voice command, execute it instead of typing.
	if e.cfg.CommandsEnabled {
		if intent := command.Parse(res.Raw); intent.IsCommand() {
			e.execCommand(intent, &res)
			res.TotalMS = time.Since(t0).Milliseconds()
			e.emit(res)
			return
		}
	}

	// LLM cleanup (optional), then the user dictionary always has the final say.
	text := res.Raw
	if e.llm != nil {
		if cr, err := e.llm.Clean(text, app); err != nil {
			fmt.Fprintf(os.Stderr, "cleanup error (using raw): %v\n", err)
		} else {
			text = cr.Text
			res.CleanMS = cr.Latency.Milliseconds()
		}
	}
	text = config.ApplyDictionary(e.cfg.Dictionary, text)
	res.Cleaned = text

	method, err := inject.Paste(text)
	if err != nil {
		res.Err = fmt.Errorf("inject: %w", err)
	}
	res.Method = method
	res.TotalMS = time.Since(t0).Milliseconds()
	e.emit(res)
}

// execCommand runs a detected voice command and records it in res.
func (e *Engine) execCommand(intent command.Intent, res *Result) {
	res.Command = intent.Name
	switch intent.Kind {
	case command.SelectAll:
		inject.SelectAll()
	case command.Undo:
		inject.Undo()
	case command.Redo:
		inject.Redo()
	case command.Copy:
		inject.Copy()
	case command.Cut:
		inject.Cut()
	case command.Paste:
		inject.PasteKey()
	case command.NewLine:
		inject.NewLine()
	case command.NewParagraph:
		inject.NewParagraph()
	case command.DeleteWord:
		inject.DeleteWord()
	case command.DeleteLine:
		inject.DeleteLine()
	case command.Rewrite:
		e.execRewrite(intent.Style, res)
	}
}

// execRewrite reads the current selection, rewrites it in the requested style via
// the LLM, and replaces it. Requires the LLM and a non-empty selection.
func (e *Engine) execRewrite(style string, res *Result) {
	if e.llm == nil {
		res.Err = fmt.Errorf("rewrite needs the cleanup LLM enabled")
		return
	}
	sel := inject.SelectedText()
	if strings.TrimSpace(sel) == "" {
		res.Err = fmt.Errorf("select some text first to rewrite it")
		return
	}
	cr, err := e.llm.Rewrite(sel, style)
	if err != nil {
		res.Err = fmt.Errorf("rewrite: %w", err)
		return
	}
	res.CleanMS = cr.Latency.Milliseconds()
	res.Raw = sel
	res.Cleaned = cr.Text
	method, err := inject.Paste(cr.Text) // AX setSelectedText replaces the selection
	if err != nil {
		res.Err = fmt.Errorf("inject: %w", err)
	}
	res.Method = method
}

func (e *Engine) emit(r Result) {
	if e.onResult != nil {
		e.onResult(r)
	}
}

// Close shuts down the engines.
func (e *Engine) Close() {
	e.mu.Lock()
	if e.rec != nil {
		e.rec.Stop()
		e.rec = nil
	}
	e.mu.Unlock()
	if e.stt != nil {
		e.stt.Close()
	}
	if e.llm != nil {
		e.llm.Close()
	}
}
