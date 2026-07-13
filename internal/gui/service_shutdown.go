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

var shutdownPlaywrightRuntime = playwrightrt.Shutdown

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

func runBoundedCleanup(ctx context.Context, name string, fn func()) bool {
	if fn == nil {
		return true
	}
	if ctx == nil {
		fn()
		return true
	}
	if err := ctx.Err(); err != nil {
		logx.Warn("gui shutdown cleanup skipped after deadline", "step", name, "error", err)
		return false
	}
	done := make(chan struct{})
	go func() {
		defer close(done)
		fn()
	}()
	select {
	case <-done:
		return true
	case <-ctx.Done():
		logx.Warn("gui shutdown cleanup timed out", "step", name, "error", ctx.Err())
		return false
	}
}

// Shutdown cancels in-flight work and waits for Playwright drivers to exit (avoids Node EPIPE on Ctrl+C).
func (s *Service) Shutdown(ctx context.Context) {
	if s == nil {
		return
	}
	s.CloseReportBridge()
	s.CancelRun()
	s.cancelActivePlugins()
	logx.Debug("gui shutdown: stopping tracked browser and allure processes only")
	s.closeBrowserForced()

	waitWaitGroup(ctx, &s.activePlaywright)
	waitWaitGroup(ctx, &s.activeBackground)

	if ctx != nil && ctx.Err() != nil {
		runSession := s.CurrentRunSession()
		allureState := s.AllureProcessState()
		recording := s.recorderOps().Session().LiveSession()
		recordingID, browserID := "", ""
		if recording != nil {
			recordingID = s.CurrentRecordSessionID()
			browserID = s.CurrentBrowserSessionID()
		}
		logx.Warn("gui shutdown timeout",
			"error", ctx.Err(),
			"active_runs", s.HasActiveRun(),
			"run_id", func() string {
				if runSession == nil {
					return ""
				}
				return runSession.RunID
			}(),
			"browser_pool_workers", func() int {
				if runSession == nil {
					return 0
				}
				return runSession.RequestSnapshot.Workers
			}(),
			"active_recorder", recording != nil && recording.CaptureEnabled(),
			"record_session_id", recordingID,
			"browser_session_id", browserID,
			"active_playwright_sessions", s.ActivePlaywrightSessions(),
			"allure_running", allureState.Running,
			"allure_pid", allureState.PID,
			"allure_dir", allureState.Dir,
			"dirty_tabs", s.HasDirtyTabs(),
		)
	}

	runBoundedCleanup(ctx, "temp_feature_cleanup", s.cleanupTempFeatureDirs)
	runBoundedCleanup(ctx, "allure_stop", s.stopAllureServe)
	runBoundedCleanup(ctx, "playwright_runtime_shutdown", shutdownPlaywrightRuntime)
	logx.Debug("gui shutdown complete")
}

func (s *Service) cancelActivePlugins() {
	if s == nil {
		return
	}
	s.mu.RLock()
	plugins := s.pluginService
	s.mu.RUnlock()
	if plugins != nil {
		plugins.CancelActive()
	}
}
