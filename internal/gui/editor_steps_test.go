package gui

import "testing"

func TestParseEditorSteps(t *testing.T) {
	text := `Feature: Demo
  Scenario: S1
    Допустим открыт "https://example.com"
    Когда нажимаю "Войти"
    Тогда проверяю текст "OK" в "#status"
`
	rows := ParseEditorSteps(text)
	if len(rows) != 3 {
		t.Fatalf("expected 3 steps, got %d", len(rows))
	}
	if rows[0].Action != "Открыть" || rows[0].Value != "https://example.com" {
		t.Fatalf("goto step: %+v", rows[0])
	}
	if rows[1].Action != "Нажать" || rows[1].Element != "Войти" {
		t.Fatalf("click step: %+v", rows[1])
	}
	if rows[2].Action != "Проверить" || rows[2].Value != "OK" || rows[2].Element != "#status" {
		t.Fatalf("assert step: %+v", rows[2])
	}
}
