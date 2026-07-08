package gui

import "testing"

func TestAnalyzeEditorSteps_ClassifiesKeywordsAndSkipsStructure(t *testing.T) {
	text := `# comment
@smoke
Функционал: Demo
Сценарий: S
	Допустим открыт "https://example.com"
	И битый шаг
`
	steps := analyzeEditorSteps(text)
	if len(steps) != 2 {
		t.Fatalf("expected 2 analyzed steps, got %d", len(steps))
	}
	if steps[0].Keyword != "Допустим" || steps[0].StepText != `открыт "https://example.com"` {
		t.Fatalf("unexpected first step: %+v", steps[0])
	}
	if steps[1].Keyword != "И" || steps[1].ParseErr == nil {
		t.Fatalf("expected second step parse error, got %+v", steps[1])
	}
}
