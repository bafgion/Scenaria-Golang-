package player

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
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

func TestBrowserPoolCloseDuringResetStopsLateRelease(t *testing.T) {
	resetStarted := make(chan struct{})
	resetContinue := make(chan struct{})
	var stopped atomic.Int32
	pool := &browserPool{
		slots:     make(chan *browserPoolSlot, 1),
		size:      1,
		closeWait: 10 * time.Millisecond,
		resetSlot: func(*browserPoolSlot) error {
			close(resetStarted)
			<-resetContinue
			return nil
		},
	}
	slot := &browserPoolSlot{index: 7, stop: func() { stopped.Add(1) }}

	done := make(chan struct{})
	go func() {
		pool.release(context.Background(), slot, PlaywrightExecutorOptions{}, ScenarioResult{}, RunCase{Name: "reset"}, "run-reset")
		close(done)
	}()
	<-resetStarted
	pool.Close()
	close(resetContinue)

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("release did not finish after close during reset")
	}
	if got := stopped.Load(); got != 1 {
		t.Fatalf("late reset release should stop worker once, got %d", got)
	}
	if len(pool.slots) != 0 {
		t.Fatalf("late reset release returned slot to closed pool")
	}
}

func TestBrowserPoolCloseDuringReplacementStopsLateRelease(t *testing.T) {
	replaceStarted := make(chan struct{})
	replaceContinue := make(chan struct{})
	var stopped atomic.Int32
	pool := &browserPool{
		slots:     make(chan *browserPoolSlot, 1),
		size:      1,
		closeWait: 10 * time.Millisecond,
		resetSlot: func(*browserPoolSlot) error {
			return errors.New("force replacement")
		},
		replaceSlot: func(_ context.Context, slot *browserPoolSlot, _ PlaywrightExecutorOptions) error {
			close(replaceStarted)
			<-replaceContinue
			slot.stop = func() { stopped.Add(1) }
			return nil
		},
	}
	slot := &browserPoolSlot{index: 8, stop: func() { stopped.Add(1) }}

	done := make(chan struct{})
	go func() {
		pool.release(context.Background(), slot, PlaywrightExecutorOptions{}, ScenarioResult{}, RunCase{Name: "replace"}, "run-replace")
		close(done)
	}()
	<-replaceStarted
	pool.Close()
	close(replaceContinue)

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("release did not finish after close during replacement")
	}
	if got := stopped.Load(); got != 1 {
		t.Fatalf("late replacement release should stop replacement worker once, got %d", got)
	}
	if len(pool.slots) != 0 {
		t.Fatalf("late replacement release returned slot to closed pool")
	}
}

func TestBrowserPoolCloseDuringUnhealthyReleaseRetiresCleanly(t *testing.T) {
	replaceStarted := make(chan struct{})
	replaceContinue := make(chan struct{})
	var stopped atomic.Int32
	pool := &browserPool{
		slots:     make(chan *browserPoolSlot, 1),
		size:      1,
		closeWait: 10 * time.Millisecond,
		resetSlot: func(*browserPoolSlot) error {
			return errors.New("force replacement")
		},
		replaceSlot: func(context.Context, *browserPoolSlot, PlaywrightExecutorOptions) error {
			close(replaceStarted)
			<-replaceContinue
			return errors.New("replacement failed")
		},
	}
	slot := &browserPoolSlot{index: 9, stop: func() { stopped.Add(1) }}

	done := make(chan struct{})
	go func() {
		pool.release(context.Background(), slot, PlaywrightExecutorOptions{}, ScenarioResult{}, RunCase{Name: "unhealthy"}, "run-unhealthy")
		close(done)
	}()
	<-replaceStarted
	pool.Close()
	close(replaceContinue)

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("unhealthy release did not finish after close")
	}
	if got := stopped.Load(); got != 1 {
		t.Fatalf("retired unhealthy slot should be stopped once, got %d", got)
	}
	if len(pool.slots) != 0 {
		t.Fatalf("unhealthy release returned slot to closed pool")
	}
}

func TestBrowserPoolConcurrentLateReleasesAfterClose(t *testing.T) {
	const workers = 4
	resetStarted := make(chan struct{}, workers)
	resetContinue := make(chan struct{})
	var stopped atomic.Int32
	pool := &browserPool{
		slots:     make(chan *browserPoolSlot, workers),
		size:      workers,
		closeWait: 10 * time.Millisecond,
		resetSlot: func(*browserPoolSlot) error {
			resetStarted <- struct{}{}
			<-resetContinue
			return nil
		},
	}

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		slot := &browserPoolSlot{index: i, stop: func() { stopped.Add(1) }}
		wg.Add(1)
		go func() {
			defer wg.Done()
			pool.release(context.Background(), slot, PlaywrightExecutorOptions{}, ScenarioResult{}, RunCase{Name: "parallel"}, "run-parallel")
		}()
	}
	for i := 0; i < workers; i++ {
		<-resetStarted
	}
	pool.Close()
	close(resetContinue)

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("concurrent late releases did not finish")
	}
	if got := stopped.Load(); got != workers {
		t.Fatalf("each late release should stop its worker, got %d", got)
	}
	if len(pool.slots) != 0 {
		t.Fatalf("concurrent late releases returned slots to closed pool")
	}
}

func TestBrowserPoolDoubleCloseStopsIdleSlotOnce(t *testing.T) {
	var stopped atomic.Int32
	pool := &browserPool{
		slots: make(chan *browserPoolSlot, 1),
		size:  1,
	}
	pool.slots <- &browserPoolSlot{index: 1, stop: func() { stopped.Add(1) }}

	pool.Close()
	pool.Close()

	if got := stopped.Load(); got != 1 {
		t.Fatalf("double close should stop idle slot once, got %d", got)
	}
}

func TestBrowserPoolReleaseAfterCloseStopsSlotWithoutSend(t *testing.T) {
	var stopped atomic.Int32
	pool := &browserPool{
		slots:  make(chan *browserPoolSlot, 1),
		size:   1,
		closed: true,
	}
	slot := &browserPoolSlot{index: 2, stop: func() { stopped.Add(1) }}

	pool.release(context.Background(), slot, PlaywrightExecutorOptions{}, ScenarioResult{}, RunCase{Name: "closed"}, "run-closed")

	if got := stopped.Load(); got != 1 {
		t.Fatalf("release after close should stop slot once, got %d", got)
	}
	if len(pool.slots) != 0 {
		t.Fatal("release after close sent slot to pool")
	}
}

func TestBrowserPoolRejectedReturnRetiresAndStopsSlot(t *testing.T) {
	var rejectedStopped atomic.Int32
	var idleStopped atomic.Int32
	pool := &browserPool{
		slots: make(chan *browserPoolSlot, 1),
		size:  1,
		resetSlot: func(*browserPoolSlot) error {
			return nil
		},
	}
	pool.slots <- &browserPoolSlot{index: 1, stop: func() { idleStopped.Add(1) }}
	rejected := &browserPoolSlot{index: 2, stop: func() { rejectedStopped.Add(1) }}

	pool.release(context.Background(), rejected, PlaywrightExecutorOptions{}, ScenarioResult{}, RunCase{Name: "full"}, "run-full")

	if got := rejectedStopped.Load(); got != 1 {
		t.Fatalf("rejected slot should be stopped once, got %d", got)
	}
	if got := idleStopped.Load(); got != 1 {
		t.Fatalf("closing exhausted pool should stop buffered idle slot once, got %d", got)
	}
	if !rejected.retired {
		t.Fatal("rejected slot should be marked retired")
	}
	if pool.retired != 1 {
		t.Fatalf("rejected slot should count as retired, got %d", pool.retired)
	}
	if len(pool.slots) != 0 {
		t.Fatal("retire-all should drain buffered idle slots")
	}
}

func TestBrowserPoolDoubleReleaseDoesNotReuseStoppedSlot(t *testing.T) {
	var stopped atomic.Int32
	pool := &browserPool{
		slots: make(chan *browserPoolSlot, 1),
		size:  1,
		resetSlot: func(*browserPoolSlot) error {
			return nil
		},
	}
	slot := &browserPoolSlot{index: 3, stop: func() { stopped.Add(1) }}

	pool.release(context.Background(), slot, PlaywrightExecutorOptions{}, ScenarioResult{}, RunCase{Name: "first"}, "run-double")
	pool.release(context.Background(), slot, PlaywrightExecutorOptions{}, ScenarioResult{}, RunCase{Name: "second"}, "run-double")

	if got := stopped.Load(); got != 1 {
		t.Fatalf("double release should stop duplicate slot once, got %d", got)
	}
	if _, err := pool.acquire(context.Background()); err == nil {
		t.Fatal("retired duplicate slot should not be reusable")
	}
}

func TestBrowserPoolRejectedReplacementSlotStopsReplacementOnce(t *testing.T) {
	var oldStopped atomic.Int32
	var replacementStopped atomic.Int32
	var idleStopped atomic.Int32
	pool := &browserPool{
		slots: make(chan *browserPoolSlot, 1),
		size:  1,
		resetSlot: func(*browserPoolSlot) error {
			return errors.New("force replacement")
		},
		replaceSlot: func(_ context.Context, slot *browserPoolSlot, _ PlaywrightExecutorOptions) error {
			slot.stopWorker()
			slot.stop = func() { replacementStopped.Add(1) }
			slot.stopOnce = sync.Once{}
			return nil
		},
	}
	pool.slots <- &browserPoolSlot{index: 4, stop: func() { idleStopped.Add(1) }}
	slot := &browserPoolSlot{index: 5, stop: func() { oldStopped.Add(1) }}

	pool.release(context.Background(), slot, PlaywrightExecutorOptions{}, ScenarioResult{}, RunCase{Name: "replacement"}, "run-replacement")

	if got := oldStopped.Load(); got != 1 {
		t.Fatalf("old worker should be stopped during replacement, got %d", got)
	}
	if got := replacementStopped.Load(); got != 1 {
		t.Fatalf("rejected replacement worker should be stopped once, got %d", got)
	}
	if got := idleStopped.Load(); got != 1 {
		t.Fatalf("retire-all should stop buffered idle worker once, got %d", got)
	}
	if !slot.retired {
		t.Fatal("rejected replacement slot should be retired")
	}
}

func TestBrowserPoolConcurrentRejectedReturnsStopEachSlot(t *testing.T) {
	const releases = 4
	var stopped atomic.Int32
	var idleStopped atomic.Int32
	pool := &browserPool{
		slots: make(chan *browserPoolSlot, 1),
		size:  1,
		resetSlot: func(*browserPoolSlot) error {
			return nil
		},
	}
	pool.slots <- &browserPoolSlot{index: 6, stop: func() { idleStopped.Add(1) }}

	var wg sync.WaitGroup
	for i := 0; i < releases; i++ {
		slot := &browserPoolSlot{index: 10 + i, stop: func() { stopped.Add(1) }}
		wg.Add(1)
		go func() {
			defer wg.Done()
			pool.release(context.Background(), slot, PlaywrightExecutorOptions{}, ScenarioResult{}, RunCase{Name: "parallel-full"}, "run-parallel-full")
		}()
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("concurrent rejected releases did not finish")
	}
	if got := stopped.Load(); got != releases {
		t.Fatalf("each rejected release should stop its slot once, got %d", got)
	}
	if got := idleStopped.Load(); got != 1 {
		t.Fatalf("buffered idle slot should be stopped once, got %d", got)
	}
}
