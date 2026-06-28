package gherkin

import "testing"

func TestExpandScenarioOutline(t *testing.T) {
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
				Examples: []Example{
					{
						Rows: [][]string{
							{"url", "title"},
							{"/catalog", "Items"},
							{"/offers", "Offers"},
						},
					},
				},
			},
		},
	}

	expanded := ExpandFeature(feature)
	if len(expanded) != 2 {
		t.Fatalf("expected 2 runnable scenarios, got %d", len(expanded))
	}
	if expanded[0].Title != "Поиск по строке — /catalog" {
		t.Fatalf("unexpected expanded title: %q", expanded[0].Title)
	}
	if expanded[0].Steps[1].Text != `вижу "Items"` {
		t.Fatalf("unexpected substituted step: %q", expanded[0].Steps[1].Text)
	}
}

func TestExpandScenarioWithBackground(t *testing.T) {
	feature := &Feature{
		Background: []Step{{Keyword: "Допустим", Text: `открыт "https://example.com"`}},
		Scenarios: []Scenario{
			{
				Title: "Check",
				Steps: []Step{{Keyword: "Тогда", Text: `вижу "h1"`}},
			},
		},
	}

	expanded := ExpandFeature(feature)
	if len(expanded) != 1 || len(expanded[0].Steps) != 2 {
		t.Fatalf("expected background + scenario steps, got %+v", expanded)
	}
}

func TestFeatureHasTag(t *testing.T) {
	feature := &Feature{
		Tags: []string{"@api"},
		Scenarios: []Scenario{
			{Title: "A", Tags: []string{"@smoke"}},
		},
	}
	if !FeatureHasTag(feature, "smoke") {
		t.Fatal("expected scenario tag to match")
	}
	if !FeatureHasTag(feature, "@api") {
		t.Fatal("expected feature tag to match")
	}
	if FeatureHasTag(feature, "slow") {
		t.Fatal("unexpected tag match")
	}
}
