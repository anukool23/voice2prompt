// Package config is Voice2Prompt's persisted settings, shared by the CLI and the
// Wails app. Settings live in a JSON file under the user config dir
// (~/Library/Application Support/Voice2Prompt/config.json on macOS).
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Settings is the full user-configurable state.
type Settings struct {
	// Trigger selects the activation mechanism: "chord" (a hotkey combo, no extra
	// permission) or "fn" (hold the Fn/🌐 key; double-tap to lock — needs Input Monitoring).
	Trigger string `json:"trigger"`
	// Hotkey is the push-to-talk chord, e.g. "Ctrl+Option+Space" (used when Trigger=="chord").
	Hotkey string `json:"hotkey"`
	// Language is the whisper language code ("en", "auto", …).
	Language string `json:"language"`
	// WhisperModel is the path to the ggml model file.
	WhisperModel string `json:"whisperModel"`
	// CleanupEnabled toggles the LLM cleanup layer.
	CleanupEnabled bool `json:"cleanupEnabled"`
	// CommandsEnabled toggles voice command mode (select all, delete word, make this formal…).
	CommandsEnabled bool `json:"commandsEnabled"`
	// LLMModel is the Ollama model used for cleanup.
	LLMModel string `json:"llmModel"`
	// Dictionary maps (case-insensitive) spoken/misheard forms to their correction,
	// applied to the transcript before injection — e.g. "github" → "GitHub".
	Dictionary map[string]string `json:"dictionary"`
	// OnboardingComplete gates the first-run walkthrough.
	OnboardingComplete bool `json:"onboardingComplete"`
}

// Defaults returns the built-in settings used on first run.
func Defaults() Settings {
	return Settings{
		Trigger:            "fn",
		Hotkey:             "Ctrl+Option+Space",
		Language:           "en",
		WhisperModel:       "models/ggml-base.en.bin",
		CleanupEnabled:     true,
		CommandsEnabled:    true,
		LLMModel:           "qwen2.5:3b",
		Dictionary:         map[string]string{},
		OnboardingComplete: false,
	}
}

// Dir returns the Voice2Prompt config directory, creating it if needed.
func Dir() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(base, "Voice2Prompt")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return dir, nil
}

// Path returns the config file path.
func Path() (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.json"), nil
}

// Load reads settings from disk, falling back to defaults for a missing file or
// any missing fields.
func Load() (Settings, error) {
	s := Defaults()
	path, err := Path()
	if err != nil {
		return s, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return s, nil // first run — defaults
		}
		return s, err
	}
	// Unmarshal over the defaults so absent keys keep their default value.
	if err := json.Unmarshal(data, &s); err != nil {
		return Defaults(), fmt.Errorf("config %s is corrupt: %w", path, err)
	}
	if s.Dictionary == nil {
		s.Dictionary = map[string]string{}
	}
	return s, nil
}

// Save writes settings to disk atomically (write temp, then rename).
func Save(s Settings) error {
	path, err := Path()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
