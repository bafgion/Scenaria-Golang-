package player

import (
	"context"
	"errors"
	"fmt"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
)

var ErrBrowserExecutionNotImplemented = errors.New("browser execution engine is not implemented yet")

type FeatureInput struct {
	Path    string
	Feature *gherkin.Feature
}

type ExecutionPlan struct {
	Features []FeatureInput
}

type ExecutionResult struct {
	Mode            string
	Files           int
	Scenarios       int
	Steps           int
	ScenarioResults []ScenarioResult
}

type ScenarioResult struct {
	FeaturePath string
	Scenario    string
	Status      string
	Message     string
}

type ScenarioInput struct {
	FeaturePath string
	Feature     *gherkin.Feature
	Scenario    *gherkin.Scenario
}

type BrowserExecutor interface {
	ExecuteScenario(ctx context.Context, input ScenarioInput) (ScenarioResult, error)
}

type Runner interface {
	Execute(ctx context.Context, plan ExecutionPlan) (ExecutionResult, error)
}

func NewRunner(dryRun bool) Runner {
	if dryRun {
		return DryRunner{}
	}
	return BrowserRunner{
		Executor: StubBrowserExecutor{},
	}
}

type DryRunner struct{}

func (DryRunner) Execute(_ context.Context, plan ExecutionPlan) (ExecutionResult, error) {
	files, scenarios, steps, scenarioResults := summarizePlan(plan)
	return ExecutionResult{
		Mode:            "dry-run",
		Files:           files,
		Scenarios:       scenarios,
		Steps:           steps,
		ScenarioResults: scenarioResults,
	}, nil
}

type BrowserRunner struct {
	Executor BrowserExecutor
}

func (r BrowserRunner) Execute(ctx context.Context, plan ExecutionPlan) (ExecutionResult, error) {
	if r.Executor == nil {
		return ExecutionResult{}, fmt.Errorf("browser runner: executor is nil")
	}
	files, scenarios, steps, _ := summarizePlan(plan)
	result := ExecutionResult{
		Mode:      "browser",
		Files:     files,
		Scenarios: scenarios,
		Steps:     steps,
	}

	for _, featureInput := range plan.Features {
		for i := range featureInput.Feature.Scenarios {
			scenario := &featureInput.Feature.Scenarios[i]
			runResult, err := r.Executor.ExecuteScenario(ctx, ScenarioInput{
				FeaturePath: featureInput.Path,
				Feature:     featureInput.Feature,
				Scenario:    scenario,
			})
			if err != nil {
				return ExecutionResult{}, err
			}
			result.ScenarioResults = append(result.ScenarioResults, runResult)
		}
	}
	return result, nil
}

type StubBrowserExecutor struct{}

func (StubBrowserExecutor) ExecuteScenario(_ context.Context, _ ScenarioInput) (ScenarioResult, error) {
	return ScenarioResult{}, ErrBrowserExecutionNotImplemented
}

func summarizePlan(plan ExecutionPlan) (files int, scenarios int, steps int, results []ScenarioResult) {
	files = len(plan.Features)
	results = make([]ScenarioResult, 0)
	for _, input := range plan.Features {
		steps += len(input.Feature.Background)
		scenarios += len(input.Feature.Scenarios)
		for _, scenario := range input.Feature.Scenarios {
			steps += len(scenario.Steps)
			results = append(results, ScenarioResult{
				FeaturePath: input.Path,
				Scenario:    scenario.Title,
				Status:      "skipped",
				Message:     "dry-run mode",
			})
		}
	}
	return files, scenarios, steps, results
}
