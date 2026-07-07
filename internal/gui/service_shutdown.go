package gui

import (
	"context"
	"sync"
	"time"

	"github.com/bafgion/scenaria-golang/internal/logx"
	"github.com/bafgion/scenaria-golang/internal/playwrightrt"
)

// DefaultShutdownTimeout is how long OnShutdown waits for Playwright runs/recording to stop cleanly.
const DefaultShutdownTimeout = 8 * time.Second

func waitWaitGroup(ctx context.Context, wg *sync.WaitGroup) {
	if wg == nil {
		return
	}
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-ctx.Done():
	}
}

// Shutdown cancels in-flight work and waits for Playwright drivers to exit (avoids Node EPIPE on Ctrl+C).
func (s *Service) Shutdown(ctx context.Context) {
	if s == nil {
		return
	}
	s.CloseReportBridge()
	s.CancelRun()
	s.closeBrowserForced()

	waitWaitGroup(ctx, &s.activePlaywright)
	waitWaitGroup(ctx, &s.activeBackground)

	s.cleanupTempFeatureDirs()
	StopAllureServe()
	playwrightrt.Shutdown()
	logx.Debug("gui shutdown complete")
}
