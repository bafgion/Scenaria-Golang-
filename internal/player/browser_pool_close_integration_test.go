//go:build integration

package player

import (
	"context"
	"testing"
	"time"
)

func TestParallelWorkersCloseBrowserDoesNotPoisonFollowingScenario(t *testing.T) {
	plan := ExecutionPlan{
		Cases: []RunCase{
			closeBrowserCase("First closes browser"),
			openPageCase("Second on fresh browser"),
			openPageCase("Third on fresh browser"),
		},
	}
	runner := BrowserRunner{
		Executor: NewPlaywrightExecutor(PlaywrightExecutorOptions{
			BrowserName:   "chromium",
			Headless:      true,
			CloseAfterRun: true,
		}),
		ParallelWorkers: 2,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	result, err := runner.Execute(ctx, plan)
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	if len(result.ScenarioResults) != 3 {
		t.Fatalf("expected 3 scenario results, got %d", len(result.ScenarioResults))
	}
	for _, sr := range result.ScenarioResults {
		if sr.Status != "passed" {
			t.Fatalf("scenario %q status = %q message = %q", sr.Scenario, sr.Status, sr.Message)
		}
		if sr.Message == "browser page is not available" || sr.Message == "browser session is closed" {
			t.Fatalf("scenario %q exposed low-level browser lifecycle error: %q", sr.Scenario, sr.Message)
		}
	}
}

func TestBrowserPoolReplacesSlotAfterCloseBrowser(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	options := PlaywrightExecutorOptions{BrowserName: "chromium", Headless: true}
	pool, err := newBrowserPool(ctx, options, 1)
	if err != nil {
		t.Fatalf("newBrowserPool: %v", err)
	}
	defer pool.Close()

	slot, err := pool.acquire(ctx)
	if err != nil {
		t.Fatal(err)
	}

	exec := NewPlaywrightExecutor(options)
	closeCase := closeBrowserCase("closes browser")
	_, err = exec.ExecuteScenarioOnSession(ctx, slot.session, scenarioInputFromCase(closeCase))
	if err != nil {
		t.Fatalf("close-browser scenario: %v", err)
	}
	pool.release(ctx, slot, options, ScenarioResult{Status: "passed"}, closeCase, "test-run")

	slot, err = pool.acquire(ctx)
	if err != nil {
		t.Fatal(err)
	}
	followUp := openPageCase("runs after slot replacement")
	result, err := exec.ExecuteScenarioOnSession(ctx, slot.session, scenarioInputFromCase(followUp))
	if err != nil {
		t.Fatalf("follow-up scenario: %v", err)
	}
	if result.Status != "passed" {
		t.Fatalf("follow-up status = %q message = %q", result.Status, result.Message)
	}
	pool.release(ctx, slot, options, result, followUp, "test-run")
}
