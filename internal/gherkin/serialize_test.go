package gherkin

import (
	"path/filepath"
	"testing"
)

func TestSerializeFeature_RoundTrip(t *testing.T) {
	feature := &Feature{
		Title: "Каталог",
		Tags:  []string{"@api"},
		Background: []Step{
			{Keyword: "Допустим", Text: `я подключаю TestClient "DemoUser"`},
		},
		Scenarios: []Scenario{
			{
				Title:     "Поиск",
				Tags:      []string{"@smoke"},
				IsOutline: true,
				Steps: []Step{
					{Keyword: "Когда", Text: `открываю "<url>"`},
					{
						Keyword:   "И",
						Text:      "отправляю payload",
						DocString: "{\n  \"query\": \"<q>\"\n}",
					},
					{
						Keyword: "Тогда",
						Text:    "вижу результат",
						Table: [][]string{
							{"поле", "значение"},
							{"q", "<q>"},
						},
					},
				},
				Examples: []Example{
					{
						Rows: [][]string{
							{"url", "q"},
							{"/catalog", "milk"},
						},
					},
				},
			},
		},
	}

	serialized, err := SerializeFeature(feature)
	if err != nil {
		t.Fatalf("SerializeFeature returned error: %v", err)
	}
	parsed, err := ParseFeature(serialized)
	if err != nil {
		t.Fatalf("ParseFeature returned error after serialize: %v", err)
	}

	if parsed.Title != feature.Title {
		t.Fatalf("unexpected title after round-trip: %q", parsed.Title)
	}
	if len(parsed.Tags) != 1 || parsed.Tags[0] != "@api" {
		t.Fatalf("unexpected tags after round-trip: %#v", parsed.Tags)
	}
	if len(parsed.Scenarios) != 1 || !parsed.Scenarios[0].IsOutline {
		t.Fatalf("unexpected scenario shape after round-trip: %+v", parsed.Scenarios)
	}
	if len(parsed.Scenarios[0].Examples) != 1 {
		t.Fatalf("unexpected examples after round-trip: %+v", parsed.Scenarios[0].Examples)
	}
}

func TestSaveFeatureFile(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "saved.feature")

	feature := &Feature{
		Title: "Smoke",
		Scenarios: []Scenario{
			{
				Title: "Simple",
				Steps: []Step{{Keyword: "Когда", Text: "выполняю шаг"}},
			},
		},
	}
	if err := SaveFeatureFile(path, feature); err != nil {
		t.Fatalf("SaveFeatureFile returned error: %v", err)
	}
	loaded, err := ParseFeatureFile(path)
	if err != nil {
		t.Fatalf("ParseFeatureFile returned error: %v", err)
	}
	if loaded.Title != "Smoke" {
		t.Fatalf("unexpected loaded feature: %+v", loaded)
	}
}
