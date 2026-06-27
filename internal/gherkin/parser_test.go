package gherkin

import "testing"

func TestParseFeature_Success(t *testing.T) {
	content := `
Функционал: Авторизация
Контекст:
  Допустим открыт сайт

Сценарий: Успешный вход
  Когда ввожу логин и пароль
  Тогда вижу главную страницу
`

	feature, err := ParseFeature(content)
	if err != nil {
		t.Fatalf("ParseFeature returned error: %v", err)
	}

	if feature.Title != "Авторизация" {
		t.Fatalf("unexpected feature title: %q", feature.Title)
	}
	if len(feature.Background) != 1 {
		t.Fatalf("unexpected background step count: %d", len(feature.Background))
	}
	if len(feature.Scenarios) != 1 {
		t.Fatalf("unexpected scenario count: %d", len(feature.Scenarios))
	}
	if len(feature.Scenarios[0].Steps) != 2 {
		t.Fatalf("unexpected step count: %d", len(feature.Scenarios[0].Steps))
	}
}

func TestParseFeature_ErrorForStepOutsideScenario(t *testing.T) {
	content := `
Функционал: Авторизация
Когда шаг не в сценарии
`

	_, err := ParseFeature(content)
	if err == nil {
		t.Fatal("expected parse error, got nil")
	}
}

func TestValidateFeature(t *testing.T) {
	feature := &Feature{
		Title: "Demo",
		Scenarios: []Scenario{
			{Title: "S1"},
			{Title: "S2", Steps: []Step{{Keyword: "Когда", Text: "test"}}},
		},
	}

	issues := ValidateFeature(feature)
	if len(issues) == 0 {
		t.Fatal("expected at least one issue")
	}
}
