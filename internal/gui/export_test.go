package gui

import "testing"

func TestPreviewExport(t *testing.T) {
	text := `Feature: Demo
  Scenario: S1
    Допустим открыт "https://example.com"
    Когда нажимаю "Войти"
`
	preview := (&Service{}).PreviewExport(text)
	if preview.StepCount != 2 {
		t.Fatalf("step count = %d, want 2", preview.StepCount)
	}
	if preview.ScenarioTitle != "Demo" {
		t.Fatalf("title = %q, want Demo", preview.ScenarioTitle)
	}
	if len(preview.Issues) != 0 {
		t.Fatalf("expected no issues, got %v", preview.Issues)
	}
}

func TestFeatureTitleFromText(t *testing.T) {
	if got := featureTitleFromText("Feature: Login flow\n"); got != "Login flow" {
		t.Fatalf("got %q", got)
	}
}
