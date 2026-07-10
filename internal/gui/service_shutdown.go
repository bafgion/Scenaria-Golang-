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

	s.cleanupTempFeatureDirs()
	s.stopAllureServe()
	playwrightrt.Shutdown()
	logx.Debug("gui shutdown complete")
}
