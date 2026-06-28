package player

import (
	"context"
	"testing"
)

type countingExecutor struct {
	calls int
}

func (c *countingExecutor) ExecuteScenario(_ context.Context, _ ScenarioInput) (ScenarioResult, error) {
	c.calls++
	return ScenarioResult{Status: "passed"}, nil
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
	if _, err := runner.Execute(context.Background(), plan); err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	if executor.calls != 3 {
		t.Fatalf("expected 3 scenario executions, got %d", executor.calls)
	}
}

func TestSetMaxLoopIterations(t *testing.T) {
	old := MaxLoopIterations
	defer func() { MaxLoopIterations = old }()
	SetMaxLoopIterations(42)
	if MaxLoopIterations != 42 {
		t.Fatalf("expected 42, got %d", MaxLoopIterations)
	}
	SetMaxLoopIterations(0)
	if MaxLoopIterations != 42 {
		t.Fatalf("zero should not change limit")
	}
}
