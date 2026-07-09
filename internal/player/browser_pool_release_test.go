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
		pool.release(ctx, slot, options, ScenarioResult{}, RunCase{}, "test-run")
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

func TestBrowserPoolCloseIgnoresRetiredSlots(t *testing.T) {
	pool := &browserPool{
		slots:   make(chan *browserPoolSlot, 1),
		size:    1,
		retired: 1,
	}
	done := make(chan struct{})
	go func() {
		pool.Close()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("pool Close waited for a retired slot")
	}
}

func TestBrowserPoolAcquireFailsWhenAllSlotsRetired(t *testing.T) {
	pool := &browserPool{
		slots:   make(chan *browserPoolSlot, 1),
		size:    1,
		retired: 1,
	}
	if _, err := pool.acquire(context.Background()); err == nil {
		t.Fatal("expected acquire to fail when all slots are retired")
	}
}
