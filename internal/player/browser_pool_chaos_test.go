package player

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestIsRetryableGotoError(t *testing.T) {
	if !isRetryableGotoError(errors.New("goto failed: timeout 30000ms exceeded")) {
		t.Fatal("timeout should be retryable")
	}
	if !isRetryableGotoError(errors.New("net::ERR_CONNECTION_RESET")) {
		t.Fatal("net error should be retryable")
	}
	if isRetryableGotoError(errors.New("404 not found")) {
		t.Fatal("404 should not be retryable")
	}
}

func TestBrowserPoolCancelWithTwoWorkersDoesNotDeadlock(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	pool, err := newBrowserPool(ctx, PlaywrightExecutorOptions{BrowserName: "chromium", Headless: true}, 2)
	if err != nil {
		t.Skip("playwright not available:", err)
	}

	slot1, err := pool.acquire(ctx)
	if err != nil {
		pool.Close()
		t.Fatal(err)
	}
	slot2, err := pool.acquire(ctx)
	if err != nil {
		pool.release(slot1)
		pool.Close()
		t.Fatal(err)
	}
	cancel()

	done := make(chan struct{})
	go func() {
		slot1.session.abortRun()
		pool.release(slot1)
		slot2.session.abortRun()
		pool.release(slot2)
		pool.Close()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Fatal("pool Close deadlocked after cancel with two workers")
	}
}
