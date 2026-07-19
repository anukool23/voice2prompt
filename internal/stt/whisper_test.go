package stt

import "testing"

func TestStripNonSpeech(t *testing.T) {
	cases := []struct{ in, want string }{
		{"[BLANK_AUDIO]", ""},
		{"[ Silence ]", ""},
		{"(music)", ""},
		{"[ Inaudible ]", ""},
		{"♪♪♪", ""},
		{"hello world", "hello world"},
		{"hello [BLANK_AUDIO] world", "hello  world"},            // inner marker removed (spacing normalized later)
		{"the array is [0] indexed", "the array is [0] indexed"}, // real brackets kept
	}
	for _, c := range cases {
		if got := stripNonSpeech(c.in); got != c.want {
			t.Errorf("stripNonSpeech(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}
