package player

import (
	"os"
	"path/filepath"
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

func TestBuildExecutionPlan_TestClientOverride(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, ".scenaria", "test_clients")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	payload := `{"name":"OverrideUser","base_url":"https://override.example","cookies":[],"local_storage":{}}`
	if err := os.WriteFile(filepath.Join(dir, "OverrideUser.json"), []byte(payload), 0o644); err != nil {
		t.Fatal(err)
	}

	feature := &gherkin.Feature{
		Title:      "Demo",
		TestClient: "FeatureUser",
		Scenarios: []gherkin.Scenario{
			{Title: "One", Steps: []gherkin.Step{{Keyword: "Когда", Text: "шаг"}}},
		},
	}
	featurePath := filepath.Join(root, "demo.feature")
	if err := os.WriteFile(featurePath, []byte("Функционал: x\nСценарий: y\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	plan := BuildExecutionPlanWithTestClient([]FeatureInput{{Path: featurePath, Feature: feature}}, "", nil, "OverrideUser")
	if len(plan.Cases) != 1 {
		t.Fatalf("expected one case, got %d", len(plan.Cases))
	}
	if plan.Cases[0].TestClient == nil || plan.Cases[0].TestClient.Name != "OverrideUser" {
		t.Fatalf("unexpected test client: %+v", plan.Cases[0].TestClient)
	}
}
