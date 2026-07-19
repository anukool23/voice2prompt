package config

import "testing"

func TestApplyDictionary(t *testing.T) {
	dict := map[string]string{
		"github":          "GitHub",
		"kubernetes":      "Kubernetes",
		"voice to prompt": "Voice2Prompt",
	}
	cases := []struct{ in, want string }{
		{"i pushed to github today", "i pushed to GitHub today"},
		{"GITHUB and Github both match", "GitHub and GitHub both match"},
		{"deploy on kubernetes", "deploy on Kubernetes"},
		{"using voice to prompt now", "using Voice2Prompt now"},
		{"githubbed is not a word match", "githubbed is not a word match"}, // \b guards partials
		{"nothing to change here", "nothing to change here"},
	}
	for _, c := range cases {
		if got := ApplyDictionary(dict, c.in); got != c.want {
			t.Errorf("ApplyDictionary(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}
