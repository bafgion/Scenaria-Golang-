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
		Cases: []RunCase{
			{
				FeaturePath: "a.feature",
				Name:        "S1",
				StartStep:   -1,
				EndStep:     -1,
				Steps: []gherkin.Step{
					{Keyword: "Допустим", Text: "шаг"},
					{Keyword: "Когда", Text: "делаю 1"},
					{Keyword: "Тогда", Text: "вижу 1"},
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

func TestSummarizePlanPartialStepRange(t *testing.T) {
	plan := ExecutionPlan{
		Cases: []RunCase{{
			FeaturePath: "a.feature",
			Name:        "S1",
			StartStep:   1,
			EndStep:     -1,
			Steps: []gherkin.Step{
				{Line: 1, Text: "шаг 1"},
				{Line: 2, Text: "шаг 2"},
				{Line: 3, Text: "шаг 3"},
			},
		}},
	}
	_, _, steps, _ := SummarizePlan(plan)
	if steps != 2 {
		t.Fatalf("expected 2 steps in partial range, got %d", steps)
	}
}

func TestBrowserRunnerNotImplemented(t *testing.T) {
	runner := BrowserRunner{Executor: StubBrowserExecutor{}}
	plan := ExecutionPlan{
		Cases: []RunCase{
			{
				FeaturePath: "a.feature",
				Name:        "S1",
				StartStep:   -1,
				EndStep:     -1,
				Steps:       []gherkin.Step{{Keyword: "Когда", Text: "x"}},
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
		Cases: []RunCase{
			{FeaturePath: "a.feature", Name: "S1", StartStep: -1, EndStep: -1, Steps: []gherkin.Step{{Keyword: "Когда", Text: "x"}}},
			{FeaturePath: "a.feature", Name: "S2", StartStep: -1, EndStep: -1, Steps: []gherkin.Step{{Keyword: "Тогда", Text: "y"}}},
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
		Scenario:    input.ScenarioName,
		Status:      "passed",
	}, nil
}
