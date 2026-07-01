package player

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
)

type failingParallelExecutor struct {
	calls atomic.Int32
}

func (f *failingParallelExecutor) ExecuteScenario(_ context.Context, input ScenarioInput) (ScenarioResult, error) {
	f.calls.Add(1)
	if input.ScenarioName == "fail-me" {
		return ScenarioResult{}, fmt.Errorf("boom at %s", input.ScenarioName)
	}
	return ScenarioResult{
		Scenario: input.ScenarioName,
		Status:   "passed",
	}, nil
}

func TestParallelRunnerFailureRecordsResultUnderMutex(t *testing.T) {
	plan := ExecutionPlan{
		Cases: []RunCase{
			{Name: "ok-1"},
			{Name: "fail-me"},
			{Name: "ok-2"},
		},
	}
	executor := &failingParallelExecutor{}
	runner := BrowserRunner{Executor: executor, ParallelWorkers: 3}
	_, err := runner.Execute(context.Background(), plan)
	if err == nil {
		t.Fatal("expected error")
	}
	calls := int(executor.calls.Load())
	if calls < 2 || calls > 3 {
		t.Fatalf("expected 2-3 calls before cancel on first error, got %d", calls)
	}
}
