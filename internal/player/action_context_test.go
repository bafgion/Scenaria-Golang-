package player

import (
	"context"
	"testing"
	"time"
)

func TestTimeoutMsRespectsContextDeadline(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	ms := *timeoutMs(ctx, LocatorWaitTimeoutMs)
	if ms > 50 {
		t.Fatalf("expected timeout <= 50ms, got %v", ms)
	}
}

func TestCapWaitDuration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	got := capWaitDuration(ctx, time.Second)
	if got > 10*time.Millisecond {
		t.Fatalf("expected capped wait, got %v", got)
	}
}

func TestSleepRetryBackoffHonorsCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := sleepRetryBackoff(ctx, 1, 5*time.Second); err == nil {
		t.Fatal("expected context error")
	}
}
