package recorder

import "testing"

func TestEscapeStepText(t *testing.T) {
	got := escapeStepText(`say "hi"` + "\n" + `line2`)
	want := `say \"hi\"\nline2`
	if got != want {
		t.Fatalf("escapeStepText = %q, want %q", got, want)
	}
}
