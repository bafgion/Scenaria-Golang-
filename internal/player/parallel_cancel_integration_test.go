//go:build integration

package player

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/mxschmitt/playwright-go"
)

type failFastPlaywrightExecutor struct {
	inner     *PlaywrightExecutor
	slowSteps []gherkin.Step
	completed atomic.Int32
}

func (e *failFastPlaywrightExecutor) ExecuteScenario(ctx context.Context, input ScenarioInput) (ScenarioResult, error) {
	if input.ScenarioName == "fail-fast" {
		return ScenarioResult{
			FeaturePath: input.FeaturePath,
			Scenario:    input.ScenarioName,
			Status:      "failed",
			Message:     "intentional failure",
		}, nil
	}
	result, err := e.inner.ExecuteScenario(ctx, ScenarioInput{
		FeaturePath:  input.FeaturePath,
		ScenarioName: input.ScenarioName,
		Steps:        e.slowSteps,
	})
	if err == nil && result.Status == "passed" {
		e.completed.Add(1)
	}
	return result, err
}

func TestParallelRunnerCancelsPlaywrightWorkers(t *testing.T) {
	if err := playwright.Install(); err != nil {
		t.Fatalf("install playwright: %v", err)
	}

	slowSteps := []gherkin.Step{
		{Line: 1, Text: `Допустим открыт "about:blank"`},
		{Line: 2, Text: `Когда жду 12 с`},
	}
	inner := NewPlaywrightExecutor(PlaywrightExecutorOptions{
		BrowserName: "chromium",
		Headless:    true,
	})
	executor := &failFastPlaywrightExecutor{
		inner:     inner,
		slowSteps: slowSteps,
	}

	plan := ExecutionPlan{
		Cases: []RunCase{
			{Name: "fail-fast"},
			{Name: "slow-1"},
			{Name: "slow-2"},
			{Name: "slow-3"},
		},
	}

	runner := BrowserRunner{Executor: executor, ParallelWorkers: 2}
	start := time.Now()
	_, err := runner.Execute(context.Background(), plan)
	elapsed := time.Since(start)
	if err == nil {
		t.Fatal("expected parallel run error")
	}
	if completed := int(executor.completed.Load()); completed > 0 {
		t.Fatalf("expected no slow scenario to finish, completed %d", completed)
	}
	if elapsed > 10*time.Second {
		t.Fatalf("parallel cancellation took too long: %v", elapsed)
	}
}
