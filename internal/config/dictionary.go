package config

import (
	"regexp"
	"sort"
)

// ApplyDictionary replaces whole-word, case-insensitive occurrences of each key
// with its value (e.g. "github" → "GitHub"). Longer keys are applied first so
// multi-word entries take precedence over their sub-words.
func ApplyDictionary(dict map[string]string, text string) string {
	if len(dict) == 0 || text == "" {
		return text
	}

	keys := make([]string, 0, len(dict))
	for k := range dict {
		if k != "" {
			keys = append(keys, k)
		}
	}
	sort.Slice(keys, func(i, j int) bool { return len(keys[i]) > len(keys[j]) })

	for _, k := range keys {
		re := regexp.MustCompile(`(?i)\b` + regexp.QuoteMeta(k) + `\b`)
		text = re.ReplaceAllString(text, dict[k])
	}
	return text
}
