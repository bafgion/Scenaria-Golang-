package player

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

type slowCancellableExecutor struct {
	delay     time.Duration
	started   atomic.Int32
	completed atomic.Int32
}

func (s *slowCancellableExecutor) ExecuteScenario(ctx context.Context, input ScenarioInput) (ScenarioResult, error) {
	s.started.Add(1)
	if input.ScenarioName == "fail-now" {
		return ScenarioResult{}, fmt.Errorf("intentional failure")
	}
	select {
	case <-ctx.Done():
		return ScenarioResult{}, ctx.Err()
	case <-time.After(s.delay):
	}
	s.completed.Add(1)
	return ScenarioResult{
		Scenario: input.ScenarioName,
		Status:   "passed",
	}, nil
}

func TestParallelRunnerCancelsOnFirstError(t *testing.T) {
	const cases = 8
	plan := ExecutionPlan{Cases: make([]RunCase, cases)}
	for i := range plan.Cases {
		name := fmt.Sprintf("slow-%d", i)
		if i == 0 {
			name = "fail-now"
		}
		plan.Cases[i] = RunCase{Name: name}
	}
	executor := &slowCancellableExecutor{delay: 400 * time.Millisecond}
	runner := BrowserRunner{Executor: executor, ParallelWorkers: 2}
	_, err := runner.Execute(context.Background(), plan)
	if err == nil {
		t.Fatal("expected error")
	}
	completed := int(executor.completed.Load())
	if completed >= cases-1 {
		t.Fatalf("expected early cancellation, %d scenarios still completed", completed)
	}
}
