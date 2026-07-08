package gui

import "testing"

func TestAnalyzeEditorContent_BundlesStepsIssuesAndHints(t *testing.T) {
	text := `Функционал: Demo
  Сценарий: S
    Допустим открыт "https://example.com"
    И открыт "https://example.org"
    И битый шаг
`
	dto := AnalyzeEditorContent(text, true)
	if len(dto.Steps) != 3 {
		t.Fatalf("expected 3 steps, got %d", len(dto.Steps))
	}
	if len(dto.Issues) == 0 {
		t.Fatal("expected validation issues")
	}
	if len(dto.Hints) == 0 {
		t.Fatal("expected scenario hints")
	}
}

func TestAnalyzeEditorContent_SkipsHintsWhenDisabled(t *testing.T) {
	text := `Функционал: Demo
  Сценарий: S
    Допустим открыт "https://example.com"
`
	dto := AnalyzeEditorContent(text, false)
	if len(dto.Hints) != 0 {
		t.Fatalf("expected no hints, got %d", len(dto.Hints))
	}
}
