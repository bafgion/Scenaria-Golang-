package player

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
)

type countingParallelExecutor struct {
	calls atomic.Int32
	names []string
	mu    sync.Mutex
}

func (e *countingParallelExecutor) ExecuteScenario(_ context.Context, input ScenarioInput) (ScenarioResult, error) {
	e.calls.Add(1)
	e.mu.Lock()
	e.names = append(e.names, input.ScenarioName)
	e.mu.Unlock()
	return ScenarioResult{
		Scenario: input.ScenarioName,
		Status:   "passed",
	}, nil
}

type variableIsolationExecutor struct {
	mu sync.Mutex
}

func (e *variableIsolationExecutor) ExecuteScenario(_ context.Context, input ScenarioInput) (ScenarioResult, error) {
	marker := "marker_" + input.ScenarioName
	input.Variables[marker] = "owned"
	e.mu.Lock()
	defer e.mu.Unlock()
	for key, value := range input.Variables {
		if strings.HasPrefix(key, "marker_") && value == "owned" && key != marker {
			return ScenarioResult{}, fmt.Errorf("variable leak: %s visible in %s", key, input.ScenarioName)
		}
	}
	return ScenarioResult{
		Scenario: input.ScenarioName,
		Status:   "passed",
	}, nil
}

func makeWorkerParityPlan(count int) ExecutionPlan {
	cases := make([]RunCase, 0, count)
	for i := 0; i < count; i++ {
		name := fmt.Sprintf("scenario-%d", i+1)
		cases = append(cases, RunCase{
			CaseID:    BuildCaseID("batch.feature", name, 0),
			Name:      name,
			Variables: map[string]string{"BASE": "shared"},
		})
	}
	return ExecutionPlan{Cases: cases}
}

func TestParallelWorkersExecuteEachCaseOnce(t *testing.T) {
	workerCounts := []int{1, 2, 4}
	for _, workers := range workerCounts {
		t.Run(fmt.Sprintf("workers-%d", workers), func(t *testing.T) {
			executor := &countingParallelExecutor{}
			runner := BrowserRunner{Executor: executor, ParallelWorkers: workers}
			plan := makeWorkerParityPlan(6)
			result, err := runner.Execute(context.Background(), plan)
			if err != nil {
				t.Fatalf("execute failed: %v", err)
			}
			if int(executor.calls.Load()) != len(plan.Cases) {
				t.Fatalf("expected %d calls, got %d", len(plan.Cases), executor.calls.Load())
			}
			if len(result.ScenarioResults) != len(plan.Cases) {
				t.Fatalf("expected %d results, got %d", len(plan.Cases), len(result.ScenarioResults))
			}
			seen := map[string]int{}
			for _, scenarioResult := range result.ScenarioResults {
				seen[scenarioResult.Scenario]++
			}
			for _, runCase := range plan.Cases {
				if seen[runCase.Name] != 1 {
					t.Fatalf("expected scenario %q exactly once, got %d", runCase.Name, seen[runCase.Name])
				}
			}
		})
	}
}

func TestParallelWorkersIsolateScenarioVariables(t *testing.T) {
	executor := &variableIsolationExecutor{}
	runner := BrowserRunner{Executor: executor, ParallelWorkers: 4}
	plan := makeWorkerParityPlan(8)
	if _, err := runner.Execute(context.Background(), plan); err != nil {
		t.Fatalf("execute failed: %v", err)
	}
}

func TestParallelWorkersStressTenScenarios(t *testing.T) {
	for _, workers := range []int{2, 4} {
		t.Run(fmt.Sprintf("workers-%d", workers), func(t *testing.T) {
			executor := &countingParallelExecutor{}
			runner := BrowserRunner{Executor: executor, ParallelWorkers: workers}
			plan := makeWorkerParityPlan(10)
			result, err := runner.Execute(context.Background(), plan)
			if err != nil {
				t.Fatalf("execute failed: %v", err)
			}
			if len(result.ScenarioResults) != 10 {
				t.Fatalf("expected 10 results, got %d", len(result.ScenarioResults))
			}
			for _, sr := range result.ScenarioResults {
				if sr.Status != "passed" {
					t.Fatalf("scenario %q status = %q", sr.Scenario, sr.Status)
				}
			}
		})
	}
}

func TestCloneVariables(t *testing.T) {
	source := map[string]string{"BASE": "https://example.com"}
	cloned := CloneVariables(source)
	cloned["EXTRA"] = "x"
	source["BASE"] = "mutated"
	if cloned["BASE"] != "https://example.com" {
		t.Fatalf("clone should not track source mutations: %#v", cloned)
	}
	if source["EXTRA"] != "" {
		t.Fatalf("clone should not mutate source: %#v", source)
	}
}
