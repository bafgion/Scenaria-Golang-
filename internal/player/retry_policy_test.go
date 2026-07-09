package player

import (
	"context"
	"errors"
	"testing"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
)

func TestShouldRetryActionSafeByDefaultRiskyOptIn(t *testing.T) {
	exec := NewStepExecutor(ExecutorOptions{})
	if !exec.shouldRetryAction("assert-visible") {
		t.Fatal("assertions should retry by default")
	}
	if !exec.shouldRetryAction("wait-visible") {
		t.Fatal("waits should retry by default")
	}
	if exec.shouldRetryAction("click") {
		t.Fatal("click should not retry by default")
	}
	if exec.shouldRetryAction("fill") {
		t.Fatal("fill should not retry by default")
	}
	if exec.shouldRetryAction("download-click") {
		t.Fatal("download-click should not retry by default")
	}
}

func TestShouldRetryActionRiskyWhenEnabled(t *testing.T) {
	enabled := true
	exec := NewStepExecutor(ExecutorOptions{
		RetryPolicy: RetryPolicy{Actions: &enabled},
	})
	if !exec.shouldRetryAction("click") {
		t.Fatal("click should retry when RetryActions enabled")
	}
}

func TestRunWithRetriesRecordsAttemptCount(t *testing.T) {
	exec := NewStepExecutor(ExecutorOptions{MaxActionRetries: 2})
	calls := 0
	attempts, err := exec.runWithRetries(context.Background(), "wait-visible", func() error {
		calls++
		if calls < 3 {
			return errors.New("timeout 30000ms exceeded")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if attempts != 3 {
		t.Fatalf("expected 3 attempts, got %d", attempts)
	}
}

func TestStepRecordStoresRetryAttempts(t *testing.T) {
	ctx := NewRunContext(nil, 1, t.TempDir())
	idx := ctx.beginLeafStep(gherkinStep("жду 1 сек"))
	ctx.setStepRetryAttempts(idx, 2)
	recs := ctx.StepRecords()
	if len(recs) != 1 || recs[0].RetryAttempts != 2 {
		t.Fatalf("unexpected retry attempts: %+v", recs)
	}
}

func gherkinStep(text string) gherkin.Step {
	return gherkin.Step{Line: 1, Keyword: "Когда", Text: text}
}