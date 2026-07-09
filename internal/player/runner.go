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

type ExecutionResult struct {
	RunID           string
	Mode            string
	Files           int
	Scenarios       int
	Steps           int
	ScenarioResults []ScenarioResult
}

type ScenarioResult struct {
	RunID         string
	CaseID        string
	FeaturePath   string
	Scenario      string
	ExampleIndex  int
	Status        string
	Message       string
	FailedStep    *int
	DurationMS    int64
	StepRecords   []StepRecord
	ScreenshotPNG []byte
	TraceZIP      []byte
	VideoWebM     []byte
	// Artifact paths are preferred over in-memory payloads to avoid retaining large blobs in heap.
	ScreenshotPath string
	TraceZIPPath   string
	VideoWebMPath  string
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

func (DryRunner) Execute(ctx context.Context, plan ExecutionPlan) (ExecutionResult, error) {
	files, scenarios, steps, scenarioResults := SummarizePlan(plan)
	result := ExecutionResult{
		Mode:            "dry-run",
		Files:           files,
		Scenarios:       scenarios,
		Steps:           steps,
		ScenarioResults: scenarioResults,
	}
	stampRunID(&result, RunIDFromContext(ctx))
	return result, nil
}

type StubBrowserExecutor struct{}

func (StubBrowserExecutor) ExecuteScenario(_ context.Context, _ ScenarioInput) (ScenarioResult, error) {
	return ScenarioResult{}, ErrBrowserExecutionNotImplemented
}

func failedStepIndex(index int) *int {
	return &index
}
