package recorder

import "testing"

func TestFormatRecordedGherkinLine(t *testing.T) {
	got := FormatRecordedGherkinLine(`нажимаю "#btn"`, "Допустим")
	want := "\tДопустим нажимаю \"#btn\""
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestFormatRecordedGherkinLineStripsExistingKeyword(t *testing.T) {
	got := FormatRecordedGherkinLine("И нажимаю \"#btn\"", "Допустим")
	want := "\tДопустим нажимаю \"#btn\""
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestFormatRecordedStepsAsGherkin(t *testing.T) {
	lines := FormatRecordedStepsAsGherkin([]RecordedStep{
		{Action: "goto", Value: "https://example.com"},
		{Action: "click", Selector: "#btn"},
	})
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	if !stringsHasPrefix(lines[0], "\tДопустим ") {
		t.Fatalf("first line = %q", lines[0])
	}
	if !stringsHasPrefix(lines[1], "\tИ ") {
		t.Fatalf("second line = %q", lines[1])
	}
}

func stringsHasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}
