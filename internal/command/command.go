// Package command detects whether a dictation utterance is actually a voice
// command (e.g. "select all", "delete last word", "make this formal") rather than
// text to type. Detection is pure and rule-based so it's fast and testable; the
// engine executes the returned Intent.
package command

import (
	"regexp"
	"strings"
)

// Kind is the type of command detected.
type Kind int

const (
	None Kind = iota // not a command — dictate the text normally
	SelectAll
	Undo
	Redo
	Copy
	Cut
	Paste
	NewLine
	NewParagraph
	DeleteWord
	DeleteLine
	Rewrite // uses Style; operates on the current selection via the LLM
)

// Intent is the parsed command.
type Intent struct {
	Kind  Kind
	Style string // for Rewrite: "formal", "casual", "shorter", "longer", "grammar", "bullets"
	Name  string // human-readable label for logging/UI
}

// IsCommand reports whether this intent is an actual command.
func (i Intent) IsCommand() bool { return i.Kind != None }

var nonWord = regexp.MustCompile(`[^\w\s]`)

// normalize lowercases, strips punctuation, and collapses whitespace.
func normalize(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = nonWord.ReplaceAllString(s, "")
	return strings.Join(strings.Fields(s), " ")
}

// exact maps a fully-normalized phrase to its command.
var exact = map[string]Intent{
	"select all": {Kind: SelectAll, Name: "select all"},

	"undo":         {Kind: Undo, Name: "undo"},
	"undo that":    {Kind: Undo, Name: "undo"},
	"scratch that": {Kind: Undo, Name: "undo"},
	"redo":         {Kind: Redo, Name: "redo"},
	"redo that":    {Kind: Redo, Name: "redo"},

	"copy":       {Kind: Copy, Name: "copy"},
	"copy that":  {Kind: Copy, Name: "copy"},
	"copy this":  {Kind: Copy, Name: "copy"},
	"cut":        {Kind: Cut, Name: "cut"},
	"cut that":   {Kind: Cut, Name: "cut"},
	"paste":      {Kind: Paste, Name: "paste"},
	"paste that": {Kind: Paste, Name: "paste"},

	"new line":      {Kind: NewLine, Name: "new line"},
	"newline":       {Kind: NewLine, Name: "new line"},
	"new paragraph": {Kind: NewParagraph, Name: "new paragraph"},

	"delete word":          {Kind: DeleteWord, Name: "delete word"},
	"delete last word":     {Kind: DeleteWord, Name: "delete word"},
	"delete the last word": {Kind: DeleteWord, Name: "delete word"},
	"delete that":          {Kind: DeleteWord, Name: "delete word"},
	"delete this":          {Kind: DeleteWord, Name: "delete word"},
	"delete line":          {Kind: DeleteLine, Name: "delete line"},
	"delete the line":      {Kind: DeleteLine, Name: "delete line"},
	"delete last line":     {Kind: DeleteLine, Name: "delete line"},
}

// rewritePrefixes gate rewrite detection so ordinary dictation that merely mentions
// "formal" etc. isn't misread as a command.
var rewritePrefixes = []string{
	"make this", "make it", "make that",
	"fix the", "fix ", "shorten", "expand", "rewrite", "turn this into", "turn that into",
}

// Parse classifies a transcript. Returns Intent{Kind: None} for ordinary dictation.
func Parse(transcript string) Intent {
	n := normalize(transcript)
	if n == "" {
		return Intent{Kind: None}
	}
	if i, ok := exact[n]; ok {
		return i
	}
	// Rewrite commands: must look like a command and be short.
	if len(strings.Fields(n)) <= 7 && hasRewritePrefix(n) {
		if style, ok := rewriteStyle(n); ok {
			return Intent{Kind: Rewrite, Style: style, Name: "rewrite: " + style}
		}
	}
	return Intent{Kind: None}
}

func hasRewritePrefix(n string) bool {
	for _, p := range rewritePrefixes {
		if strings.HasPrefix(n, p) {
			return true
		}
	}
	return false
}

func rewriteStyle(n string) (string, bool) {
	switch {
	case strings.Contains(n, "formal"):
		return "formal", true
	case strings.Contains(n, "casual") || strings.Contains(n, "friendly"):
		return "casual", true
	case strings.Contains(n, "shorter") || strings.Contains(n, "shorten") || strings.Contains(n, "concise"):
		return "shorter", true
	case strings.Contains(n, "longer") || strings.Contains(n, "expand"):
		return "longer", true
	case strings.Contains(n, "grammar") || strings.Contains(n, "spelling"):
		return "grammar", true
	case strings.Contains(n, "bullet"):
		return "bullets", true
	}
	return "", false
}
