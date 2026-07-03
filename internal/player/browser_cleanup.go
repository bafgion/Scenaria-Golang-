package player

import (
	"context"
	"sync"

	"github.com/bafgion/scenaria-golang/internal/logx"
	playwright "github.com/mxschmitt/playwright-go"
)

// watchContext closes the browser session when ctx is cancelled (parallel run abort, user stop).
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
				logx.Debug("closing browser session after context cancellation")
				s.closeLocked(false)
			})
		case <-done:
		}
	}()
	return func() { close(done) }
}

// startPlaywright runs Playwright; callers invoke the returned stop() when done (see CloseAfterRun).
func startPlaywright(ctx context.Context) (*playwright.Playwright, func(), error) {
	pw, err := playwright.Run()
	if err != nil {
		return nil, nil, err
	}
	var once sync.Once
	stop := func() {
		once.Do(func() {
			closeBrowserResource("playwright", pw.Stop)
		})
	}
	_ = ctx // lifecycle: callers stop Playwright explicitly (see CloseAfterRun).
	return pw, stop, nil
}
