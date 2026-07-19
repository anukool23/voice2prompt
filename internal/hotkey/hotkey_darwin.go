//go:build darwin

// Package hotkey centralizes the push-to-talk key binding.
package hotkey

import gh "golang.design/x/hotkey"

// Chord is the human-readable description of the push-to-talk binding.
const Chord = "Ctrl+Option+Space"

// PushToTalk returns the configured push-to-talk hotkey.
//
// Carbon (which golang.design/x/hotkey uses on macOS) can't register a
// modifier-only hotkey, so we pair modifiers with a real key.
func PushToTalk() *gh.Hotkey {
	return gh.New([]gh.Modifier{gh.ModCtrl, gh.ModOption}, gh.KeySpace)
}
