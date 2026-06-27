package player

import (
	"context"
	"errors"
	"testing"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
)

func TestDryRunnerExecute(t *testing.T) {
	runner := DryRunner{}
	plan := ExecutionPlan{
		Features: []FeatureInput{
			{
				Path: "a.feature",
				Feature: &gherkin.Feature{
					Title: "A",
					Background: []gherkin.Step{
						{Keyword: "Допустим", Text: "шаг"},
					},
					Scenarios: []gherkin.Scenario{
						{
							Title: "S1",
							Steps: []gherkin.Step{
								{Keyword: "Когда", Text: "делаю 1"},
								{Keyword: "Тогда", Text: "вижу 1"},
							},
						},
					},
				},
			},
		},
	}

	result, err := runner.Execute(context.Background(), plan)
	if err != nil {
		t.Fatalf("DryRunner.Execute returned error: %v", err)
	}
	if result.Mode != "dry-run" || result.Files != 1 || result.Scenarios != 1 || result.Steps != 3 {
		t.Fatalf("unexpected run result: %+v", result)
	}
	if len(result.ScenarioResults) != 1 || result.ScenarioResults[0].Status != "skipped" {
		t.Fatalf("unexpected scenario dry-run results: %+v", result.ScenarioResults)
	}
}

func TestBrowserRunnerNotImplemented(t *testing.T) {
	runner := BrowserRunner{Executor: StubBrowserExecutor{}}
	plan := ExecutionPlan{
		Features: []FeatureInput{
			{
				Path: "a.feature",
				Feature: &gherkin.Feature{
					Title: "A",
					Scenarios: []gherkin.Scenario{
						{Title: "S1", Steps: []gherkin.Step{{Keyword: "Когда", Text: "x"}}},
					},
				},
			},
		},
	}
	_, err := runner.Execute(context.Background(), plan)
	if !errors.Is(err, ErrBrowserExecutionNotImplemented) {
		t.Fatalf("expected not-implemented error, got: %v", err)
	}
}

func TestBrowserRunnerWithExecutor(t *testing.T) {
	runner := BrowserRunner{
		Executor: fakeExecutor{},
	}
	plan := ExecutionPlan{
		Features: []FeatureInput{
			{
				Path: "a.feature",
				Feature: &gherkin.Feature{
					Title: "A",
					Scenarios: []gherkin.Scenario{
						{Title: "S1", Steps: []gherkin.Step{{Keyword: "Когда", Text: "x"}}},
						{Title: "S2", Steps: []gherkin.Step{{Keyword: "Тогда", Text: "y"}}},
					},
				},
			},
		},
	}

	result, err := runner.Execute(context.Background(), plan)
	if err != nil {
		t.Fatalf("BrowserRunner.Execute returned error: %v", err)
	}
	if result.Mode != "browser" || result.Scenarios != 2 || len(result.ScenarioResults) != 2 {
		t.Fatalf("unexpected browser result: %+v", result)
	}
	if result.ScenarioResults[0].Status != "passed" {
		t.Fatalf("unexpected scenario status: %+v", result.ScenarioResults[0])
	}
}

func TestBrowserRunnerNilExecutor(t *testing.T) {
	runner := BrowserRunner{}
	_, err := runner.Execute(context.Background(), ExecutionPlan{})
	if err == nil {
		t.Fatal("expected nil executor error")
	}
}

type fakeExecutor struct{}

func (fakeExecutor) ExecuteScenario(_ context.Context, input ScenarioInput) (ScenarioResult, error) {
	return ScenarioResult{
		FeaturePath: input.FeaturePath,
		Scenario:    input.Scenario.Title,
		Status:      "passed",
	}, nil
}
