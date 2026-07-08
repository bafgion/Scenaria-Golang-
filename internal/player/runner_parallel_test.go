package player

import (
	"context"
	"errors"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type blockingScenarioExecutor struct {
	block time.Duration
	done  atomic.Int32
}

func (e *blockingScenarioExecutor) ExecuteScenario(ctx context.Context, input ScenarioInput) (ScenarioResult, error) {
	if input.ScenarioName == "fail-fast" {
		return ScenarioResult{
			FeaturePath: input.FeaturePath,
			Scenario:    input.ScenarioName,
			Status:      "failed",
			Message:     "intentional failure",
		}, nil
	}
	select {
	case <-ctx.Done():
		e.done.Add(1)
		return ScenarioResult{
			FeaturePath: input.FeaturePath,
			Scenario:    input.ScenarioName,
			Status:      "failed",
			Message:     ctx.Err().Error(),
		}, ctx.Err()
	case <-time.After(e.block):
		return ScenarioResult{
			FeaturePath: input.FeaturePath,
			Scenario:    input.ScenarioName,
			Status:      "passed",
		}, nil
	}
}

func TestParallelRunnerCancelsSiblingContexts(t *testing.T) {
	executor := &blockingScenarioExecutor{block: 3 * time.Second}
	runner := BrowserRunner{Executor: executor, ParallelWorkers: 3}
	plan := ExecutionPlan{
		Cases: []RunCase{
			{Name: "fail-fast"},
			{Name: "slow-a"},
			{Name: "slow-b"},
		},
	}
	start := time.Now()
	result, err := runner.Execute(context.Background(), plan)
	elapsed := time.Since(start)
	if err == nil {
		t.Fatal("expected parallel run error")
	}
	canceled := 0
	for _, sr := range result.ScenarioResults {
		if sr.Scenario == "fail-fast" {
			continue
		}
		if sr.Status == "canceled" || (sr.Status == "failed" && strings.Contains(sr.Message, "canceled")) {
			canceled++
		}
	}
	if canceled < 1 && int(executor.done.Load()) < 1 {
		t.Fatalf("expected sibling cancel, results=%+v done=%d", result.ScenarioResults, executor.done.Load())
	}
	if elapsed > 2*time.Second {
		t.Fatalf("parallel cancel took too long: %v", elapsed)
	}
}

func TestFailFastParallelCancelHonorsContinueOnFail(t *testing.T) {
	ctx := WithContinueOnFail(context.Background(), true)
	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	failFastParallelCancel(runCtx, cancel, nil)
	if err := runCtx.Err(); err != nil {
		t.Fatalf("expected continue-on-fail to skip cancel, got %v", err)
	}
}

func TestFailFastParallelCancelStopsRunContext(t *testing.T) {
	runCtx, cancel := context.WithCancel(context.Background())
	failFastParallelCancel(runCtx, cancel, nil)
	if !errors.Is(runCtx.Err(), context.Canceled) {
		t.Fatalf("expected canceled run context, got %v", runCtx.Err())
	}
}

func TestBrowserPoolAbortActiveSessions(t *testing.T) {
	ctx := context.Background()
	pool, err := newBrowserPool(ctx, PlaywrightExecutorOptions{BrowserName: "chromium", Headless: true}, 1)
	if err != nil {
		t.Skip("playwright not available:", err)
	}
	defer pool.Close()

	pool.abortActiveSessions()
	slot, err := pool.acquire(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !slot.session.isClosed() {
		t.Fatal("expected idle pool session to be aborted")
	}
	if err := slot.session.resetForScenario(); err != nil {
		t.Fatalf("reset after abort: %v", err)
	}
	pool.release(slot)
}

var _ BrowserExecutor = (*blockingScenarioExecutor)(nil)

// Ensure blockingScenarioExecutor is safe under parallel calls.
func TestBlockingScenarioExecutorConcurrent(t *testing.T) {
	executor := &blockingScenarioExecutor{block: 10 * time.Millisecond}
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = executor.ExecuteScenario(context.Background(), ScenarioInput{ScenarioName: "ok"})
		}()
	}
	wg.Wait()
}
