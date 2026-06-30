package player

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
)

type countingExecutor struct {
	calls atomic.Int32
}

func (c *countingExecutor) ExecuteScenario(_ context.Context, input ScenarioInput) (ScenarioResult, error) {
	c.calls.Add(1)
	return ScenarioResult{
		Scenario: input.ScenarioName,
		Status:   "passed",
	}, nil
}

func TestBrowserRunnerParallelWorkers(t *testing.T) {
	executor := &countingExecutor{}
	plan := ExecutionPlan{
		Cases: []RunCase{
			{Name: "A"},
			{Name: "B"},
			{Name: "C"},
		},
	}
	runner := BrowserRunner{Executor: executor, ParallelWorkers: 2}
	result, err := runner.Execute(context.Background(), plan)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	if got := int(executor.calls.Load()); got != 3 {
		t.Fatalf("expected 3 scenario executions, got %d", got)
	}
	if len(result.ScenarioResults) != 3 {
		t.Fatalf("expected 3 scenario results, got %d", len(result.ScenarioResults))
	}
	for i, scenarioResult := range result.ScenarioResults {
		want := string(rune('A' + i))
		if scenarioResult.Scenario != want {
			t.Fatalf("result[%d].Scenario = %q, want %q", i, scenarioResult.Scenario, want)
		}
	}
}

func TestBrowserRunnerParallelPreservesAllResults(t *testing.T) {
	executor := &namedParallelExecutor{}
	plan := ExecutionPlan{
		Cases: make([]RunCase, 16),
	}
	for i := range plan.Cases {
		plan.Cases[i] = RunCase{Name: fmt.Sprintf("scenario-%02d", i)}
	}
	runner := BrowserRunner{Executor: executor, ParallelWorkers: 4}
	result, err := runner.Execute(context.Background(), plan)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	if len(result.ScenarioResults) != len(plan.Cases) {
		t.Fatalf("expected %d results, got %d", len(plan.Cases), len(result.ScenarioResults))
	}
	seen := make(map[string]struct{}, len(plan.Cases))
	for i, scenarioResult := range result.ScenarioResults {
		if scenarioResult.Scenario == "" {
			t.Fatalf("result[%d] has empty scenario name", i)
		}
		seen[scenarioResult.Scenario] = struct{}{}
	}
	if len(seen) != len(plan.Cases) {
		t.Fatalf("expected %d unique scenario names, got %d", len(plan.Cases), len(seen))
	}
}

type namedParallelExecutor struct{}

func (namedParallelExecutor) ExecuteScenario(_ context.Context, input ScenarioInput) (ScenarioResult, error) {
	return ScenarioResult{
		Scenario: input.ScenarioName,
		Status:   "passed",
	}, nil
}

func TestExecutorMaxLoopIterations(t *testing.T) {
	exec := NewStepExecutor(ExecutorOptions{MaxLoopIterations: 42})
	if got := exec.maxLoopIterations(); got != 42 {
		t.Fatalf("expected 42, got %d", got)
	}
	exec = NewStepExecutor(ExecutorOptions{})
	if got := exec.maxLoopIterations(); got != DefaultMaxLoopIterations {
		t.Fatalf("expected default %d, got %d", DefaultMaxLoopIterations, got)
	}
}
