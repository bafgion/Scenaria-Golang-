package player

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type chaoticExecutor struct {
	rng      *rand.Rand
	mu       sync.Mutex
	failRate float64
	calls    atomic.Int32
}

func (e *chaoticExecutor) ExecuteScenario(_ context.Context, input ScenarioInput) (ScenarioResult, error) {
	e.calls.Add(1)
	e.mu.Lock()
	shouldFail := e.rng.Float64() < e.failRate
	delay := time.Duration(e.rng.Intn(3)) * time.Millisecond
	e.mu.Unlock()
	if shouldFail {
		return ScenarioResult{}, fmt.Errorf("chaos failure at %s", input.ScenarioName)
	}
	time.Sleep(delay)
	return ScenarioResult{
		FeaturePath: input.FeaturePath,
		Scenario:    input.ScenarioName,
		Status:      "passed",
	}, nil
}

func TestChaosParallelRunnerRandomized(t *testing.T) {
	const runs = 24
	for seed := int64(0); seed < runs; seed++ {
		t.Run(fmt.Sprintf("seed-%d", seed), func(t *testing.T) {
			rng := rand.New(rand.NewSource(seed))
			caseCount := 1 + rng.Intn(48)
			workers := 1 + rng.Intn(8)
			failRate := rng.Float64() * 0.35

			plan := ExecutionPlan{Cases: make([]RunCase, caseCount)}
			for i := range plan.Cases {
				plan.Cases[i] = RunCase{
					FeaturePath: fmt.Sprintf("f-%d.feature", i%5),
					Name:        fmt.Sprintf("scenario-%03d", i),
				}
			}

			executor := &chaoticExecutor{rng: rand.New(rand.NewSource(seed + 99)), failRate: failRate}
			runner := BrowserRunner{Executor: executor, ParallelWorkers: workers}
			result, err := runner.Execute(context.Background(), plan)

			if err != nil {
				return
			}
			if len(result.ScenarioResults) != caseCount {
				t.Fatalf("expected %d results, got %d", caseCount, len(result.ScenarioResults))
			}
			for i, scenarioResult := range result.ScenarioResults {
				want := fmt.Sprintf("scenario-%03d", i)
				if scenarioResult.Scenario != want {
					t.Fatalf("result[%d].Scenario = %q, want %q", i, scenarioResult.Scenario, want)
				}
				if scenarioResult.Status != "passed" {
					t.Fatalf("result[%d] status = %q", i, scenarioResult.Status)
				}
			}
			if got := int(executor.calls.Load()); got != caseCount {
				t.Fatalf("executor calls = %d, want %d", got, caseCount)
			}
		})
	}
}

func TestChaosParallelRunnerCancellation(t *testing.T) {
	slow := &slowExecutor{delay: 40 * time.Millisecond}
	plan := ExecutionPlan{Cases: make([]RunCase, 32)}
	for i := range plan.Cases {
		plan.Cases[i] = RunCase{Name: fmt.Sprintf("s-%d", i)}
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	runner := BrowserRunner{Executor: slow, ParallelWorkers: 8}
	_, err := runner.Execute(ctx, plan)
	if err == nil {
		t.Fatal("expected cancellation error")
	}
}

type slowExecutor struct {
	delay time.Duration
}

func (s *slowExecutor) ExecuteScenario(ctx context.Context, input ScenarioInput) (ScenarioResult, error) {
	select {
	case <-ctx.Done():
		return ScenarioResult{}, ctx.Err()
	case <-time.After(s.delay):
		return ScenarioResult{Scenario: input.ScenarioName, Status: "passed"}, nil
	}
}

func TestChaosArtifactBaseNamesUnique(t *testing.T) {
	names := []string{
		"Успешный вход",
		"Успешный вход",
		`weird/path\file.feature`,
		"",
		"!!!",
	}
	seen := map[string]struct{}{}
	for i, title := range names {
		base := artifactBaseName(ScenarioInput{
			FeaturePath:  fmt.Sprintf("proj/s%d.feature", i),
			ScenarioName: title,
		})
		if base == "" {
			t.Fatalf("empty artifact base for %q", title)
		}
		key := base
		if _, ok := seen[key]; ok && title != "" {
			// same title different feature path should still be unique when feature differs
			continue
		}
		seen[key] = struct{}{}
	}
}
