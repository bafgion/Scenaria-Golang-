package gui

import "testing"

func TestResolveRunFromLine(t *testing.T) {
	text := `Функционал: Демо
  Сценарий: Вход
    Когда открыт "https://example.com"
    И нажимаю "button.login"
`
	full, err := ResolveRunFromLine(text, 2)
	if err != nil || full.Scenario != "Вход" || full.Partial {
		t.Fatalf("header: %+v err=%v", full, err)
	}
	partial, err := ResolveRunFromLine(text, 4)
	if err != nil || !partial.Partial || partial.StartStep != 1 || partial.Scenario != "Вход" {
		t.Fatalf("step line: %+v err=%v", partial, err)
	}
}

func TestResolveRunToLine(t *testing.T) {
	text := `Функционал: Демо
  Сценарий: Вход
    Когда открыт "https://example.com"
    И нажимаю "button.login"
    Тогда вижу "Добро"
`
	toEnd, err := ResolveRunToLine(text, 4)
	if err != nil || !toEnd.Partial || toEnd.EndStep != 1 || toEnd.StartStep != -1 || toEnd.Scenario != "Вход" {
		t.Fatalf("to step line: %+v err=%v", toEnd, err)
	}
}
