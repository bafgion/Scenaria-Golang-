package player

import (
	"context"
	"errors"
	"testing"
)

func TestIsRetryableAction(t *testing.T) {
	if !isRetryableAction("assert-visible") {
		t.Fatal("assert-visible should be retryable")
	}
	if !isRetryableAction("select-option") {
		t.Fatal("select-option should be retryable")
	}
	if isRetryableAction("goto") {
		t.Fatal("goto should not be retryable")
	}
}

func TestIsRetryableStepError(t *testing.T) {
	if !isRetryableStepError(errors.New("timeout 30000ms exceeded")) {
		t.Fatal("timeout should be retryable")
	}
	if !isRetryableStepError(errors.New("element is not attached to the DOM")) {
		t.Fatal("detached element should be retryable")
	}
	if isRetryableStepError(errors.New("text mismatch: expected foo")) {
		t.Fatal("assertion mismatch should not be retryable")
	}
	if isRetryableStepError(context.Canceled) {
		t.Fatal("context cancel should not be retryable")
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
