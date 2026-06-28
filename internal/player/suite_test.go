package player

import (
	"testing"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
)

func TestBuildExecutionPlan_ExpandsOutline(t *testing.T) {
	feature := &gherkin.Feature{
		Title: "Каталог",
		Tags:  []string{"@api"},
		Scenarios: []gherkin.Scenario{
			{
				Title:     "Поиск по строке",
				Tags:      []string{"@smoke"},
				IsOutline: true,
				Steps: []gherkin.Step{
					{Keyword: "Когда", Text: `открываю "<url>"`},
					{Keyword: "Тогда", Text: `вижу "<title>"`},
				},
				Examples: []gherkin.Example{
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

	plan := BuildExecutionPlan([]FeatureInput{{Path: "catalog.feature", Feature: feature}}, "", nil)
	if len(plan.Cases) != 2 {
		t.Fatalf("expected 2 expanded scenarios, got %d", len(plan.Cases))
	}
	if plan.Cases[0].Steps[0].Text != `открываю "/catalog"` {
		t.Fatalf("unexpected first step: %q", plan.Cases[0].Steps[0].Text)
	}
	if plan.Cases[1].Steps[1].Text != `вижу "Offers"` {
		t.Fatalf("unexpected second scenario step: %q", plan.Cases[1].Steps[1].Text)
	}
}

func TestBuildExecutionPlan_TagFilter(t *testing.T) {
	feature := &gherkin.Feature{
		Title: "Demo",
		Scenarios: []gherkin.Scenario{
			{Title: "Smoke", Tags: []string{"@smoke"}, Steps: []gherkin.Step{{Keyword: "Когда", Text: "шаг"}}},
			{Title: "Other", Tags: []string{"@slow"}, Steps: []gherkin.Step{{Keyword: "Когда", Text: "шаг"}}},
		},
	}

	plan := BuildExecutionPlan([]FeatureInput{{Path: "demo.feature", Feature: feature}}, "smoke", nil)
	if len(plan.Cases) != 1 || plan.Cases[0].Name != "Smoke" {
		t.Fatalf("unexpected tag-filtered plan: %+v", plan.Cases)
	}
}

func TestBuildExecutionPlan_SkipsUntaggedFile(t *testing.T) {
	feature := &gherkin.Feature{
		Title: "Demo",
		Scenarios: []gherkin.Scenario{
			{Title: "Only", Steps: []gherkin.Step{{Keyword: "Когда", Text: "шаг"}}},
		},
	}

	plan := BuildExecutionPlan([]FeatureInput{{Path: "demo.feature", Feature: feature}}, "smoke", nil)
	if len(plan.Cases) != 0 {
		t.Fatalf("expected empty plan for missing tag, got %+v", plan.Cases)
	}
}
