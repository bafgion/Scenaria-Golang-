package recorder

import "testing"

func TestAssignRecordedStepKeywords(t *testing.T) {
	lines := []string{
		`открыт "https://example.com"`,
		`нажимаю "#login"`,
		`ввожу "user@example.com" в "#email"`,
	}
	steps := AssignRecordedStepKeywords(lines, 0)
	if len(steps) != 3 {
		t.Fatalf("got %d steps", len(steps))
	}
	if steps[0].Keyword != "Допустим" || steps[0].Text != `открыт "https://example.com"` {
		t.Fatalf("first: %+v", steps[0])
	}
	if steps[1].Keyword != "И" || steps[2].Keyword != "И" {
		t.Fatalf("continuations: %+v", steps[1:])
	}
}

func TestAssignRecordedStepKeywordsAppendsAfterExisting(t *testing.T) {
	steps := AssignRecordedStepKeywords([]string{`нажимаю "#ok"`}, 2)
	if len(steps) != 1 || steps[0].Keyword != "И" {
		t.Fatalf("got %+v", steps)
	}
}

func TestAssignRecordedStepKeywordsPreservesExplicit(t *testing.T) {
	steps := AssignRecordedStepKeywords([]string{`Тогда вижу "h1"`}, 0)
	if len(steps) != 1 || steps[0].Keyword != "Тогда" {
		t.Fatalf("got %+v", steps[0])
	}
}
