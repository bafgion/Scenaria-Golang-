package player

import playwright "github.com/mxschmitt/playwright-go"

func startTraceRecording(session *browserSession) error {
	if session == nil || session.context == nil {
		return nil
	}
	return session.context.Tracing().Start(playwright.TracingStartOptions{
		Screenshots: playwright.Bool(true),
		Snapshots:   playwright.Bool(true),
	})
}

// discardTraceRecording drops in-memory trace after a passed scenario and starts a fresh recording.
// Trace archives are only persisted on failure (see captureFailureArtifacts).
func discardTraceRecording(session *browserSession) {
	if session == nil || !session.traceEnabled {
		return
	}
	session.mu.Lock()
	defer session.mu.Unlock()
	if session.isClosed() || session.context == nil || session.traceStopped {
		return
	}
	_ = session.context.Tracing().Stop()
	session.traceStopped = true
	if err := startTraceRecording(session); err == nil {
		session.traceStopped = false
	}
}

// restartTraceRecording begins a new trace segment after a failure archive was saved.
func restartTraceRecording(session *browserSession) {
	if session == nil || !session.traceEnabled {
		return
	}
	session.mu.Lock()
	defer session.mu.Unlock()
	if session.isClosed() || session.context == nil || !session.traceStopped {
		return
	}
	if err := startTraceRecording(session); err == nil {
		session.traceStopped = false
	}
}
