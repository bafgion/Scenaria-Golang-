package player

import (
	"testing"
	"time"
)

func TestPendingAsyncCountTracksDrains(t *testing.T) {
	ResetAsyncDrainForTests(2 * time.Second)
	before := PendingAsyncCount()
	ch := make(chan error, 1)
	drainAsync(ch)
	if got := PendingAsyncCount(); got != before+1 {
		t.Fatalf("expected pending %d, got %d", before+1, got)
	}
	ch <- nil
	ResetAsyncDrainForTests(2 * time.Second)
	if got := PendingAsyncCount(); got != before {
		t.Fatalf("expected pending back to %d, got %d", before, got)
	}
}

func TestDrainPendingAsyncWaitsForAll(t *testing.T) {
	before := PendingAsyncCount()
	ch := make(chan int, 1)
	drainChan(ch)
	ch <- 1
	ResetAsyncDrainForTests(2 * time.Second)
	if got := PendingAsyncCount(); got != before {
		t.Fatalf("expected pending %d, got %d", before, got)
	}
}
