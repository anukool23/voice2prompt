package main

import (
	"context"
	_ "embed"
	"os"
	"path/filepath"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"voice2prompt/internal/autostart"
	"voice2prompt/internal/config"
	"voice2prompt/internal/engine"
	"voice2prompt/internal/hotkey"
	"voice2prompt/internal/inject"
	"voice2prompt/internal/tray"
)

//go:embed tray-icon.png
var trayIconPNG []byte

// App is the Wails backend. Its exported methods are callable from the frontend as
// window.go.main.App.<Method>().
type App struct {
	ctx context.Context

	mu      sync.Mutex
	eng     *engine.Engine
	hk      *hotkey.Carbon
	fn      *hotkey.FnController
	running bool
	history []Utterance
}

// Utterance is a UI-facing (JSON-serializable) view of an engine.Result.
type Utterance struct {
	Raw       string  `json:"raw"`
	Cleaned   string  `json:"cleaned"`
	App       string  `json:"app"`
	Method    string  `json:"method"`
	Command   string  `json:"command"`
	AudioSecs float64 `json:"audioSecs"`
	InferMS   int64   `json:"inferMS"`
	CleanMS   int64   `json:"cleanMS"`
	TotalMS   int64   `json:"totalMS"`
	Error     string  `json:"error"`
}

func toUtterance(r engine.Result) Utterance {
	u := Utterance{
		Raw: r.Raw, Cleaned: r.Cleaned, App: r.App, Method: string(r.Method), Command: r.Command,
		AudioSecs: r.AudioSecs, InferMS: r.InferMS, CleanMS: r.CleanMS, TotalMS: r.TotalMS,
	}
	if r.Err != nil {
		u.Error = r.Err.Error()
	}
	return u
}

// NewApp constructs the backend.
func NewApp() *App { return &App{} }

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.startTray()
}

func (a *App) shutdown(ctx context.Context) { a.StopEngine() }

// startTray installs the menu-bar icon with Open / Start-Stop / Quit actions.
func (a *App) startTray() {
	iconPath := filepath.Join(os.TempDir(), "voice2prompt-tray.png")
	if err := os.WriteFile(iconPath, trayIconPNG, 0o644); err != nil {
		return // tray is best-effort; the window still works
	}
	tray.Start(iconPath,
		func() { runtime.WindowShow(a.ctx) },
		func() { go a.trayToggle() },
		func() { runtime.Quit(a.ctx) },
	)
}

func (a *App) trayToggle() {
	if a.EngineRunning() {
		a.StopEngine()
	} else {
		_ = a.StartEngine()
	}
}

// --- Settings ---------------------------------------------------------------

// GetSettings returns the persisted settings (or defaults).
func (a *App) GetSettings() config.Settings {
	s, _ := config.Load()
	return s
}

// SaveSettings persists settings. If the engine is running it is restarted so the
// new hotkey/model/cleanup choices take effect.
func (a *App) SaveSettings(s config.Settings) error {
	if err := config.Save(s); err != nil {
		return err
	}
	a.mu.Lock()
	wasRunning := a.running
	a.mu.Unlock()
	if wasRunning {
		a.StopEngine()
		return a.StartEngine()
	}
	return nil
}

// CompleteOnboarding marks the first-run flow done.
func (a *App) CompleteOnboarding() error {
	s, _ := config.Load()
	s.OnboardingComplete = true
	return config.Save(s)
}

// --- Permissions ------------------------------------------------------------

// AccessibilityTrusted reports whether text injection permission is granted.
func (a *App) AccessibilityTrusted() bool { return inject.Trusted() }

// RequestAccessibility shows the macOS Accessibility prompt.
func (a *App) RequestAccessibility() { inject.PromptTrust() }

// MicStatus reports microphone permission ("authorized"/"denied"/…).
func (a *App) MicStatus() string { return inject.MicStatus() }

// RequestMicrophone triggers the microphone permission prompt.
func (a *App) RequestMicrophone() { inject.RequestMic() }

// InputMonitoringStatus reports Input Monitoring permission (needed for the Fn trigger).
func (a *App) InputMonitoringStatus() string { return inject.InputMonitoringStatus() }

// RequestInputMonitoring triggers the Input Monitoring permission prompt.
func (a *App) RequestInputMonitoring() { inject.RequestInputMonitoring() }

// LaunchAtLoginEnabled reports whether the app starts at login.
func (a *App) LaunchAtLoginEnabled() bool { return autostart.Enabled() }

// SetLaunchAtLogin toggles launch-at-login.
func (a *App) SetLaunchAtLogin(on bool) error { return autostart.SetEnabled(on) }

// --- Engine control ---------------------------------------------------------

// EngineRunning reports whether dictation is active.
func (a *App) EngineRunning() bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.running
}

// StartEngine brings up the STT/LLM engines and registers the push-to-talk hotkey.
func (a *App) StartEngine() error {
	a.mu.Lock()
	if a.running {
		a.mu.Unlock()
		return nil
	}
	a.mu.Unlock()

	cfg, _ := config.Load()
	eng := engine.New(cfg)
	eng.SetOnResult(a.onResult)
	if err := eng.Start(); err != nil {
		return err
	}

	onDown := func() { _ = eng.StartCapture() }

	var hk *hotkey.Carbon
	var fn *hotkey.FnController
	if cfg.Trigger == "fn" {
		fn = hotkey.NewFn(onDown, eng.StopCapture)
		if err := fn.Register(); err != nil {
			eng.Close()
			return err
		}
	} else {
		hk = hotkey.NewCarbon()
		if err := hk.Register(cfg.Hotkey, onDown, eng.StopCapture); err != nil {
			eng.Close()
			return err
		}
	}

	a.mu.Lock()
	a.eng, a.hk, a.fn, a.running = eng, hk, fn, true
	a.mu.Unlock()

	runtime.EventsEmit(a.ctx, "engine:state", true)
	return nil
}

// StopEngine unregisters the hotkey and shuts the engines down.
func (a *App) StopEngine() {
	a.mu.Lock()
	hk, fn, eng := a.hk, a.fn, a.eng
	a.eng, a.hk, a.fn, a.running = nil, nil, nil, false
	a.mu.Unlock()

	if hk != nil {
		hk.Unregister()
	}
	if fn != nil {
		fn.Unregister()
	}
	if eng != nil {
		eng.Close()
	}
	if a.ctx != nil {
		runtime.EventsEmit(a.ctx, "engine:state", false)
	}
}

// History returns recent utterances, newest first.
func (a *App) History() []Utterance {
	a.mu.Lock()
	defer a.mu.Unlock()
	out := make([]Utterance, len(a.history))
	copy(out, a.history)
	return out
}

func (a *App) onResult(r engine.Result) {
	u := toUtterance(r)
	a.mu.Lock()
	a.history = append([]Utterance{u}, a.history...)
	if len(a.history) > 25 {
		a.history = a.history[:25]
	}
	a.mu.Unlock()
	if a.ctx != nil {
		runtime.EventsEmit(a.ctx, "utterance", u)
	}
}
