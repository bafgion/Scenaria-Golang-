package player

import (
	"context"
	"testing"
	"time"
)

func TestBrowserPoolReleaseAfterAbortDoesNotDeadlockClose(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	pool, err := newBrowserPool(ctx, PlaywrightExecutorOptions{BrowserName: "chromium", Headless: true}, 1)
	if err != nil {
		t.Skip("playwright not available:", err)
	}

	options := PlaywrightExecutorOptions{BrowserName: "chromium", Headless: true}
	slot, err := pool.acquire(ctx)
	if err != nil {
		pool.Close()
		t.Fatal(err)
	}
	slot.session.abortRun()

	done := make(chan struct{})
	go func() {
		pool.release(ctx, slot, options, ScenarioResult{}, RunCase{})
		pool.Close()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		cancel()
		t.Fatal("pool Close deadlocked after aborted release")
	}
	cancel()
}
