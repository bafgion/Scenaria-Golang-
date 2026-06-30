package gherkin

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExpandScenarioFromParamsFile(t *testing.T) {
	dir := t.TempDir()
	featurePath := filepath.Join(dir, "search.feature")
	paramsPath := filepath.Join(dir, "search.params.json")
	if err := os.WriteFile(paramsPath, []byte(`{
  "scenarios": {
    "Поиск по строке": [
      {"url": "/catalog", "title": "Items"},
      {"url": "/offers", "title": "Offers"}
    ]
  }
}`), 0o644); err != nil {
		t.Fatal(err)
	}

	feature := &Feature{
		Title: "Каталог",
		Scenarios: []Scenario{
			{
				Title:     "Поиск по строке",
				IsOutline: true,
				Steps: []Step{
					{Keyword: "Когда", Text: `открываю "<url>"`},
					{Keyword: "Тогда", Text: `вижу "<title>"`},
				},
			},
		},
	}

	expanded := ExpandFeatureAtPath(feature, featurePath)
	if len(expanded) != 2 {
		t.Fatalf("expected 2 runnable scenarios, got %d", len(expanded))
	}
	if expanded[0].Title != "Поиск по строке — /catalog" {
		t.Fatalf("unexpected title: %q", expanded[0].Title)
	}
	if expanded[1].Steps[1].Text != `вижу "Offers"` {
		t.Fatalf("unexpected substituted step: %q", expanded[1].Steps[1].Text)
	}
}

func TestInlineExamplesTakePrecedenceOverParamsFile(t *testing.T) {
	dir := t.TempDir()
	featurePath := filepath.Join(dir, "search.feature")
	paramsPath := filepath.Join(dir, "search.params.json")
	_ = os.WriteFile(paramsPath, []byte(`{"scenarios":{"Поиск":[{"url":"/from-json"}]}}`), 0o644)

	feature := &Feature{
		Scenarios: []Scenario{
			{
				Title:     "Поиск",
				IsOutline: true,
				Steps:     []Step{{Keyword: "Когда", Text: `открываю "<url>"`}},
				Examples: []Example{{
					Rows: [][]string{{"url"}, {"/inline"}},
				}},
			},
		},
	}
	expanded := ExpandFeatureAtPath(feature, featurePath)
	if len(expanded) != 1 || expanded[0].Steps[0].Text != `открываю "/inline"` {
		t.Fatalf("expected inline example, got %+v", expanded)
	}
}
