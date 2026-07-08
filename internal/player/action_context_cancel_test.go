package player

import (
	"context"
	"testing"
	"time"

	playwright "github.com/mxschmitt/playwright-go"
)

func TestPageGotoRespectsContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := pageGoto(ctx, nil, "https://example.com", nil)
	if err == nil {
		t.Fatal("expected context error")
	}
}

func TestWaitForLocatorRespectsContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := waitForLocator(ctx, nil, playwright.LocatorWaitForOptions{})
	if err == nil {
		t.Fatal("expected context error")
	}
}

func TestPageGotoOnceRespectsPreCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := pageGotoOnce(ctx, nil, "https://example.com", nil)
	if err == nil {
		t.Fatal("expected context error")
	}
}

func TestCanceledDrainDoesNotLeakPendingCount(t *testing.T) {
	before := PendingAsyncCount()
	for i := 0; i < 20; i++ {
		ctx, cancel := context.WithCancel(context.Background())
		ch := make(chan error, 1)
		go func() {
			select {
			case <-ctx.Done():
				ch <- ctx.Err()
			case <-time.After(5 * time.Second):
				ch <- nil
			}
		}()
		cancel()
		drainAsync(ch)
	}
	drainPendingAsync(3 * time.Second)
	if got := PendingAsyncCount(); got != before {
		t.Fatalf("expected pending %d after 20 cancels, got %d", before, got)
	}
}
