package gherkin

import (
	"path/filepath"
	"testing"
)

func TestParseFeature_Success(t *testing.T) {
	content := `
@smoke @auth
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
	if len(feature.Tags) != 2 {
		t.Fatalf("unexpected feature tags count: %d", len(feature.Tags))
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

func TestParseFeature_OutlineExamplesDocStringAndTable(t *testing.T) {
	content := `
Функционал: Поиск
@outline
Структура сценария: Ищу товар
  Когда открываю страницу "<url>"
  И отправляю json
  """
  {
    "q": "<query>"
  }
  """
  Тогда вижу результаты
  | поле | значение |
  | q    | <query>  |

Примеры:
  | url      | query |
  | /search  | milk  |
`

	feature, err := ParseFeature(content)
	if err != nil {
		t.Fatalf("ParseFeature returned error: %v", err)
	}
	if len(feature.Scenarios) != 1 {
		t.Fatalf("unexpected scenario count: %d", len(feature.Scenarios))
	}
	scenario := feature.Scenarios[0]
	if !scenario.IsOutline {
		t.Fatal("scenario must be marked as outline")
	}
	if len(scenario.Tags) != 1 || scenario.Tags[0] != "@outline" {
		t.Fatalf("unexpected scenario tags: %#v", scenario.Tags)
	}
	if len(scenario.Examples) != 1 {
		t.Fatalf("unexpected examples count: %d", len(scenario.Examples))
	}
	if gotRows := len(scenario.Examples[0].Rows); gotRows != 2 {
		t.Fatalf("unexpected examples row count: %d", gotRows)
	}
	if gotDoc := scenario.Steps[1].DocString; gotDoc == "" {
		t.Fatal("expected docstring on second step")
	}
	if gotTableRows := len(scenario.Steps[2].Table); gotTableRows != 2 {
		t.Fatalf("unexpected step table row count: %d", gotTableRows)
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
		Line:  2,
		Scenarios: []Scenario{
			{Title: "S1", Line: 3},
			{
				Title:     "S2",
				Line:      8,
				IsOutline: true,
				Steps:     []Step{{Keyword: "Когда", Text: "test"}},
				Examples: []Example{
					{Line: 11, Rows: [][]string{{"a"}}},
				},
			},
		},
	}

	issues := ValidateFeature(feature)
	if len(issues) != 2 {
		t.Fatalf("expected two issues, got %d", len(issues))
	}
}

func TestParseFeature_UnterminatedDocString(t *testing.T) {
	content := `
Функционал: Demo
Сценарий: Broken
  Когда шаг
  """
  text
`

	if _, err := ParseFeature(content); err == nil {
		t.Fatal("expected parse error for unterminated docstring")
	}
}

func TestParseFeatureFile_FromTestdata(t *testing.T) {
	path := filepath.Join("testdata", "outline.feature")
	feature, err := ParseFeatureFile(path)
	if err != nil {
		t.Fatalf("ParseFeatureFile returned error: %v", err)
	}
	if feature.Title != "Каталог" {
		t.Fatalf("unexpected feature title: %q", feature.Title)
	}
	if len(feature.Scenarios) != 1 || len(feature.Scenarios[0].Examples) != 1 {
		t.Fatalf("unexpected outline parse result: %+v", feature.Scenarios)
	}
}

func TestParseFeatureFile_RejectExamplesForRegularScenario(t *testing.T) {
	path := filepath.Join("testdata", "invalid.feature")
	if _, err := ParseFeatureFile(path); err == nil {
		t.Fatal("expected parse error for invalid examples usage")
	}
}
