//go:build darwin

package hotkey

import (
	"errors"
	"fmt"
	"strings"
)

var errRegister = errors.New("hotkey: RegisterEventHotKey failed (chord may already be taken)")

// keyCodes maps key names to macOS ANSI virtual key codes (kVK_*).
var keyCodes = map[string]uint{
	"space": 49, "return": 36, "enter": 36, "tab": 48, "escape": 53, "esc": 53,
	"a": 0, "s": 1, "d": 2, "f": 3, "h": 4, "g": 5, "z": 6, "x": 7, "c": 8, "v": 9,
	"b": 11, "q": 12, "w": 13, "e": 14, "r": 15, "y": 16, "t": 17,
	"o": 31, "u": 32, "i": 34, "p": 35, "l": 37, "j": 38, "k": 40, "n": 45, "m": 46,
	"1": 18, "2": 19, "3": 20, "4": 21, "5": 23, "6": 22, "7": 26, "8": 28, "9": 25, "0": 29,
	"f1": 122, "f2": 120, "f3": 99, "f4": 118, "f5": 96, "f6": 97, "f7": 98, "f8": 100,
	"f9": 101, "f10": 109, "f11": 103, "f12": 111, "f13": 105, "f14": 107, "f15": 113,
	"f16": 106, "f17": 64, "f18": 79, "f19": 80, "f20": 90,
}

// ParseChord turns e.g. "Ctrl+Option+Space" into a virtual key code and Carbon
// modifier mask. Modifiers: Ctrl/Control, Cmd/Command/Super/Meta, Opt/Option/Alt, Shift.
func ParseChord(chord string) (keyCode, mods uint, err error) {
	if strings.TrimSpace(chord) == "" {
		return 0, 0, errors.New("hotkey is empty")
	}
	var key string
	for _, p := range strings.Split(chord, "+") {
		switch t := strings.TrimSpace(strings.ToLower(p)); t {
		case "ctrl", "control":
			mods |= modControl
		case "cmd", "command", "super", "meta":
			mods |= modCmd
		case "opt", "option", "alt":
			mods |= modOption
		case "shift":
			mods |= modShift
		case "":
			// tolerate stray separators
		default:
			key = t
		}
	}
	if key == "" {
		return 0, 0, fmt.Errorf("hotkey %q has no non-modifier key", chord)
	}
	kc, ok := keyCodes[key]
	if !ok {
		return 0, 0, fmt.Errorf("unsupported key %q in hotkey %q", key, chord)
	}
	return kc, mods, nil
}
