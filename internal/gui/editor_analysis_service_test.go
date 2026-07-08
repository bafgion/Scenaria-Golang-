package gui

import "testing"

func TestEditorAnalysisServiceAnalyzeEditorContent(t *testing.T) {
	svc := NewEditorAnalysisService()
	dto := svc.AnalyzeEditorContent(`Функционал: Demo
Сценарий: A
  Допустим открыт "https://example.com"
  И битый шаг
`, true)
	if len(dto.Steps) != 2 {
		t.Fatalf("expected 2 steps, got %d", len(dto.Steps))
	}
	if len(dto.Issues) == 0 {
		t.Fatal("expected validation issues")
	}
}

func TestEditorAnalysisServiceResolveRunRange(t *testing.T) {
	svc := NewEditorAnalysisService()
	text := `Функционал: Demo
Сценарий: A
  Допустим открыт "https://example.com"
  И кликаю "#go"
`
	from, err := svc.ResolveRunFromLine(text, 3)
	if err != nil {
		t.Fatalf("ResolveRunFromLine: %v", err)
	}
	if from.Scenario == "" {
		t.Fatal("expected scenario name")
	}
	to, err := svc.ResolveRunToLine(text, 4)
	if err != nil {
		t.Fatalf("ResolveRunToLine: %v", err)
	}
	if to.Scenario == "" {
		t.Fatal("expected scenario name")
	}
}
