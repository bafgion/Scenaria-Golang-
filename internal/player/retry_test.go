package player

import "testing"

func TestIsRetryableAction(t *testing.T) {
	if !isRetryableAction("assert-visible") {
		t.Fatal("assert-visible should be retryable")
	}
	if isRetryableAction("goto") {
		t.Fatal("goto should not be retryable")
	}
}

func TestMaxActionRetriesDefaults(t *testing.T) {
	exec := NewStepExecutor(ExecutorOptions{})
	if got := exec.maxActionRetries(); got != DefaultMaxActionRetries {
		t.Fatalf("expected default %d, got %d", DefaultMaxActionRetries, got)
	}
	exec = NewStepExecutor(ExecutorOptions{MaxActionRetries: -1})
	if got := exec.maxActionRetries(); got != 0 {
		t.Fatalf("expected 0 retries, got %d", got)
	}
}
