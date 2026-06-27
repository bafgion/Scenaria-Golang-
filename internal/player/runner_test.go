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
}

func TestBrowserRunnerNotImplemented(t *testing.T) {
	runner := BrowserRunner{}
	_, err := runner.Execute(context.Background(), ExecutionPlan{})
	if !errors.Is(err, ErrBrowserExecutionNotImplemented) {
		t.Fatalf("expected not-implemented error, got: %v", err)
	}
}
