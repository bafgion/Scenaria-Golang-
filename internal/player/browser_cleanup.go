package player

import (
	"context"
	"sync"
	"time"

	"github.com/bafgion/scenaria-golang/internal/logx"
	"github.com/bafgion/scenaria-golang/internal/playwrightrt"
	playwright "github.com/mxschmitt/playwright-go"
)

// WatchContext aborts in-flight Playwright work when ctx is cancelled (user stop, run timeout).
func (s *browserSession) WatchContext(ctx context.Context) func() {
	return s.watchContext(ctx)
}

// watchContext marks the session cancelled when ctx is cancelled (user stop, run timeout).
// Browser resources are torn down by the runner/pool defer — not here — to avoid EPIPE on the Playwright driver.
func (s *browserSession) watchContext(ctx context.Context) func() {
	if s == nil || ctx == nil {
		return func() {}
	}
	done := make(chan struct{})
	var once sync.Once
	go func() {
		select {
		case <-ctx.Done():
			once.Do(func() {
				logx.Debug("aborting browser session after context cancellation")
				s.abortRun()
			})
		case <-done:
		}
	}()
	return func() { close(done) }
}

// abortRun stops new steps without closing Playwright sockets (prevents Node EPIPE crashes).
func (s *browserSession) abortRun() {
	if s == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.external {
		s.setClosed()
		return
	}
	s.setClosed()
}

// startPlaywright runs Playwright; callers invoke the returned stop() when done (see CloseAfterRun).
func startPlaywright(ctx context.Context) (*playwright.Playwright, func(), error) {
	pw, release, err := playwrightrt.Acquire(ctx)
	if err != nil {
		return nil, nil, err
	}
	var once sync.Once
	stop := func() {
		once.Do(func() {
			drainPendingAsync(2 * time.Second)
			time.Sleep(playwrightDrainDelay)
			release()
		})
	}
	return pw, stop, nil
}

const playwrightDrainDelay = 75 * time.Millisecond

// stopBrowserWorker tears down a pool worker; brief delay before pw.Stop reduces Playwright Node EPIPE crashes.
func stopBrowserWorker(stopWatch func(), session *browserSession, stopPW func()) {
	if stopWatch != nil {
		stopWatch()
	}
	if session != nil {
		session.close()
		time.Sleep(playwrightDrainDelay)
	}
	if stopPW != nil {
		stopPW()
	}
}
