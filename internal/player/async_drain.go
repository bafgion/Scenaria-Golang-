package player

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/bafgion/scenaria-golang/internal/logx"
)

const longDrainLogAfter = 2 * time.Second

var (
	pendingAsync sync.WaitGroup
	pendingCount atomic.Int64
)

// PendingAsyncCount returns the number of in-flight async drain goroutines.
func PendingAsyncCount() int64 {
	return pendingCount.Load()
}

// ResetAsyncDrainForTests waits for pending async drains and resets observable counters.
//
//nolint:revive // exported for integration tests across packages.
func ResetAsyncDrainForTests(maxWait time.Duration) {
	drainPendingAsync(maxWait)
	pendingCount.Store(0)
}

func trackAsyncDrain(op string, fn func()) {
	pendingAsync.Add(1)
	pendingCount.Add(1)
	go func() {
		defer pendingAsync.Done()
		defer pendingCount.Add(-1)
		started := time.Now()
		warn := time.AfterFunc(longDrainLogAfter, func() {
			logx.Warn("async drain still running",
				"op", op,
				"pending", pendingCount.Load(),
				"elapsed_ms", time.Since(started).Milliseconds(),
			)
		})
		fn()
		warn.Stop()
		if elapsed := time.Since(started); elapsed >= longDrainLogAfter {
			logx.Warn("async drain completed after long wait",
				"op", op,
				"elapsed_ms", elapsed.Milliseconds(),
			)
		}
	}()
}

func drainAsync(ch <-chan error) {
	trackAsyncDrain("error", func() { <-ch })
}

func drainChan[T any](ch <-chan T) {
	trackAsyncDrain("chan", func() { <-ch })
}

func drainPendingAsync(maxWait time.Duration) {
	if maxWait <= 0 {
		pendingAsync.Wait()
		return
	}
	started := time.Now()
	done := make(chan struct{})
	go func() {
		pendingAsync.Wait()
		close(done)
	}()
	select {
	case <-done:
		if elapsed := time.Since(started); elapsed >= longDrainLogAfter {
			logx.Debug("async drains settled", "elapsed_ms", elapsed.Milliseconds())
		}
	case <-time.After(maxWait):
		logx.Warn("async drain wait timed out",
			"max_wait_ms", maxWait.Milliseconds(),
			"pending", pendingCount.Load(),
			"elapsed_ms", time.Since(started).Milliseconds(),
		)
	}
}
