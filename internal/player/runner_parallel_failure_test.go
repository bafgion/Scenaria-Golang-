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

type statusFailingParallelExecutor struct {
	calls atomic.Int32
}

func (f *statusFailingParallelExecutor) ExecuteScenario(_ context.Context, input ScenarioInput) (ScenarioResult, error) {
	f.calls.Add(1)
	if input.ScenarioName == "fail-status" {
		return ScenarioResult{
			Scenario: input.ScenarioName,
			Status:   "failed",
			Message:  "boom",
		}, nil
	}
	return ScenarioResult{
		Scenario: input.ScenarioName,
		Status:   "passed",
	}, nil
}

func TestParallelRunnerFailureReturnsPartialResults(t *testing.T) {
	plan := ExecutionPlan{
		Cases: []RunCase{
			{Name: "ok-1"},
			{Name: "fail-status"},
			{Name: "ok-2"},
		},
	}
	executor := &statusFailingParallelExecutor{}
	runner := BrowserRunner{Executor: executor, ParallelWorkers: 3}
	result, err := runner.Execute(context.Background(), plan)
	if err == nil {
		t.Fatal("expected error")
	}
	if len(result.ScenarioResults) != len(plan.Cases) {
		t.Fatalf("expected %d partial results, got %d", len(plan.Cases), len(result.ScenarioResults))
	}
	if result.ScenarioResults[1].Status != "failed" {
		t.Fatalf("expected failed status at index 1, got %+v", result.ScenarioResults[1])
	}
}

func TestSequentialRunnerFailureReturnsPartialResults(t *testing.T) {
	plan := ExecutionPlan{
		Cases: []RunCase{
			{Name: "ok-1"},
			{Name: "fail-status"},
		},
	}
	executor := &statusFailingParallelExecutor{}
	runner := BrowserRunner{Executor: executor, ParallelWorkers: 1}
	result, err := runner.Execute(context.Background(), plan)
	if err == nil {
		t.Fatal("expected error")
	}
	if len(result.ScenarioResults) != 2 {
		t.Fatalf("expected 2 partial results, got %d", len(result.ScenarioResults))
	}
	if result.ScenarioResults[0].Status != "passed" || result.ScenarioResults[1].Status != "failed" {
		t.Fatalf("unexpected statuses: %+v", result.ScenarioResults)
	}
}

func TestParallelRunnerFailureRecordsFailedStatus(t *testing.T) {
	plan := ExecutionPlan{
		Cases: []RunCase{
			{Name: "ok-1"},
			{Name: "fail-status"},
			{Name: "ok-2"},
		},
	}
	executor := &statusFailingParallelExecutor{}
	runner := BrowserRunner{Executor: executor, ParallelWorkers: 3}
	_, err := runner.Execute(context.Background(), plan)
	if err == nil {
		t.Fatal("expected error")
	}
	calls := int(executor.calls.Load())
	if calls < 2 || calls > 3 {
		t.Fatalf("expected 2-3 calls before cancel on failed status, got %d", calls)
	}
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
