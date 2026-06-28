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
	files, scenarios, steps, scenarioResults := SummarizePlan(plan)
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
	files, scenarios, steps, _ := SummarizePlan(plan)
	result := ExecutionResult{
		Mode:      "browser",
		Files:     files,
		Scenarios: scenarios,
		Steps:     steps,
	}

	for _, runCase := range plan.Cases {
		runResult, err := r.Executor.ExecuteScenario(ctx, ScenarioInput{
			FeaturePath:  runCase.FeaturePath,
			ScenarioName: runCase.Name,
			Steps:        runCase.Steps,
			TestClient:   runCase.TestClient,
			Variables:    runCase.Variables,
			ProjectRoot:  runCase.ProjectRoot,
		})
		if err != nil {
			return ExecutionResult{}, err
		}
		result.ScenarioResults = append(result.ScenarioResults, runResult)
	}
	return result, nil
}

type StubBrowserExecutor struct{}

func (StubBrowserExecutor) ExecuteScenario(_ context.Context, _ ScenarioInput) (ScenarioResult, error) {
	return ScenarioResult{}, ErrBrowserExecutionNotImplemented
}
