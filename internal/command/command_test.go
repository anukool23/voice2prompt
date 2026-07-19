package command

import "testing"

func TestParseCommands(t *testing.T) {
	cases := []struct {
		in    string
		kind  Kind
		style string
	}{
		// editing commands (with varied casing/punctuation)
		{"Select all.", SelectAll, ""},
		{"scratch that", Undo, ""},
		{"Undo that!", Undo, ""},
		{"copy this", Copy, ""},
		{"new paragraph", NewParagraph, ""},
		{"delete last word", DeleteWord, ""},
		{"delete the line", DeleteLine, ""},
		// rewrite commands
		{"make this formal", Rewrite, "formal"},
		{"make it more casual", Rewrite, "casual"},
		{"shorten this", Rewrite, "shorter"},
		{"fix the grammar", Rewrite, "grammar"},
		// NOT commands — ordinary dictation
		{"select all the files in the folder and zip them", None, ""},
		{"this needs to be a formal letter to the board", None, ""},
		{"let's copy the approach from last quarter", None, ""},
		{"please delete my account when you get a chance", None, ""},
		{"", None, ""},
	}
	for _, c := range cases {
		got := Parse(c.in)
		if got.Kind != c.kind {
			t.Errorf("Parse(%q).Kind = %v, want %v", c.in, got.Kind, c.kind)
		}
		if got.Kind == Rewrite && got.Style != c.style {
			t.Errorf("Parse(%q).Style = %q, want %q", c.in, got.Style, c.style)
		}
	}
}
