package player

import (
	"context"
	"errors"

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
	Mode      string
	Files     int
	Scenarios int
	Steps     int
}

type Runner interface {
	Execute(ctx context.Context, plan ExecutionPlan) (ExecutionResult, error)
}

func NewRunner(dryRun bool) Runner {
	if dryRun {
		return DryRunner{}
	}
	return BrowserRunner{}
}

type DryRunner struct{}

func (DryRunner) Execute(_ context.Context, plan ExecutionPlan) (ExecutionResult, error) {
	files, scenarios, steps := summarizePlan(plan)
	return ExecutionResult{
		Mode:      "dry-run",
		Files:     files,
		Scenarios: scenarios,
		Steps:     steps,
	}, nil
}

type BrowserRunner struct{}

func (BrowserRunner) Execute(_ context.Context, _ ExecutionPlan) (ExecutionResult, error) {
	return ExecutionResult{}, ErrBrowserExecutionNotImplemented
}

func summarizePlan(plan ExecutionPlan) (files int, scenarios int, steps int) {
	files = len(plan.Features)
	for _, input := range plan.Features {
		steps += len(input.Feature.Background)
		scenarios += len(input.Feature.Scenarios)
		for _, scenario := range input.Feature.Scenarios {
			steps += len(scenario.Steps)
		}
	}
	return files, scenarios, steps
}
