package playwrightrt

import (
	"context"
	"sync"
	"time"

	"github.com/bafgion/scenaria-golang/internal/logx"
	playwright "github.com/mxschmitt/playwright-go"
)

const (
	stopDrainDelay          = 75 * time.Millisecond
	startupCancelDrainLimit = 5 * time.Second
)

// Runtime shares one Playwright Node driver per process (ref-counted).
type Runtime struct {
	mu   sync.Mutex
	pw   *playwright.Playwright
	refs int
}

var defaultRuntime Runtime

// Acquire returns the shared Playwright driver; call release when done.
func Acquire(ctx context.Context) (*playwright.Playwright, func(), error) {
	return defaultRuntime.acquire(ctx)
}

// Shutdown stops the shared driver (app exit); ignores remaining refs.
func Shutdown() {
	defaultRuntime.shutdown()
}

func (r *Runtime) acquire(ctx context.Context) (*playwright.Playwright, func(), error) {
	if ctx != nil {
		if err := ctx.Err(); err != nil {
			return nil, nil, err
		}
	}

	r.mu.Lock()
	if r.pw != nil {
		r.refs++
		pw := r.pw
		r.mu.Unlock()
		return pw, r.releaseLocked, nil
	}
	r.mu.Unlock()

	started, err := runPlaywright(ctx)
	if err != nil {
		return nil, nil, err
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if r.pw != nil {
		_ = started.Stop()
		r.refs++
		return r.pw, r.releaseLocked, nil
	}
	r.pw = started
	r.refs = 1
	return r.pw, r.releaseLocked, nil
}

func (r *Runtime) releaseLocked() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.refs > 0 {
		r.refs--
	}
	if r.refs > 0 || r.pw == nil {
		return
	}
	pw := r.pw
	r.pw = nil
	go func() {
		time.Sleep(stopDrainDelay)
		_ = pw.Stop()
	}()
}

func (r *Runtime) shutdown() {
	r.mu.Lock()
	pw := r.pw
	r.pw = nil
	r.refs = 0
	r.mu.Unlock()
	if pw == nil {
		return
	}
	time.Sleep(stopDrainDelay)
	_ = pw.Stop()
}

func runPlaywright(ctx context.Context) (*playwright.Playwright, error) {
	if ctx == nil {
		return playwright.Run()
	}
	type pwResult struct {
		pw  *playwright.Playwright
		err error
	}
	ch := make(chan pwResult, 1)
	go func() {
		pw, err := playwright.Run()
		ch <- pwResult{pw: pw, err: err}
	}()
	select {
	case <-ctx.Done():
		logx.Debug("playwright startup cancelled, draining")
		drainStart := time.Now()
		select {
		case r := <-ch:
			if r.err == nil && r.pw != nil {
				_ = r.pw.Stop()
			}
		case <-time.After(startupCancelDrainLimit):
			logx.Warn("playwright startup cancel drain timed out",
				"max_wait_ms", startupCancelDrainLimit.Milliseconds(),
				"elapsed_ms", time.Since(drainStart).Milliseconds(),
			)
		}
		return nil, ctx.Err()
	case r := <-ch:
		if r.err != nil {
			return nil, r.err
		}
		return r.pw, nil
	}
}
