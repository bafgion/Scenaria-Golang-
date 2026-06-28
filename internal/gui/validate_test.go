package gui

import "testing"

func TestValidateFeatureContent(t *testing.T) {
	issues := ValidateFeatureContent("Функционал: X\nСценарий: Y\n  Когда нажимаю \"#ok\"\n  Когда битый шаг\n")
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d: %#v", len(issues), issues)
	}
	if issues[0].Line != 4 {
		t.Fatalf("unexpected line: %d", issues[0].Line)
	}
}
