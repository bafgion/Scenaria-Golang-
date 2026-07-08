package gui

import (
	"context"
	"path/filepath"
	"testing"
)

func TestStopRecordingCaptureIdempotent(t *testing.T) {
	svc, session := attachLiveSession(t, true)
	if err := session.BeginCapture(); err != nil {
		t.Fatal(err)
	}
	stopped, err := svc.StopRecordingCapture()
	if err != nil || !stopped {
		t.Fatalf("first stop: stopped=%v err=%v", stopped, err)
	}
	stoppedAgain, err := svc.StopRecordingCapture()
	if err != nil {
		t.Fatalf("second stop: %v", err)
	}
	if stoppedAgain {
		t.Fatal("expected second stop to be a no-op")
	}
}

func TestGuardedRecordEmitSkipsRetiredGeneration(t *testing.T) {
	svc := &Service{}
	svc.mu.Lock()
	svc.recordGen = 1
	svc.mu.Unlock()
	emitted := 0
	emit := svc.guardedRecordEmit(1, "features/a.feature", func(string, any) {
		emitted++
	})
	emit("record-step", map[string]any{"index": 0, "line": "step"})
	if emitted != 1 {
		t.Fatalf("expected 1 emit, got %d", emitted)
	}
	svc.mu.Lock()
	svc.recordGen = 2
	svc.mu.Unlock()
	emit("record-step", map[string]any{"index": 1, "line": "late"})
	if emitted != 1 {
		t.Fatalf("expected stale emit to be skipped, got %d", emitted)
	}
}

func TestCloseBrowserRetiresRecordGeneration(t *testing.T) {
	svc, session := attachLiveSession(t, true)
	svc.mu.Lock()
	svc.recordGen = 3
	svc.recordSessionID = "record-3"
	svc.browserSessionID = "browser-3"
	svc.recordTargetPath = "features/live.feature"
	ctx, cancel := context.WithCancel(context.Background())
	svc.recordCtx = ctx
	svc.recordCancel = cancel
	svc.mu.Unlock()

	emitted := 0
	emit := svc.guardedRecordEmit(3, svc.recordTargetPath, func(string, any) { emitted++ })

	svc.CloseBrowser()
	emit("record-step", map[string]any{"index": 0, "line": "late"})
	if emitted != 0 {
		t.Fatalf("expected retired generation to block emit, got %d", emitted)
	}
	if session.BrowserAlive() {
		t.Fatal("expected session cleared")
	}
	if got := svc.LastClosedBrowserSessionID(); got != "browser-3" {
		t.Fatalf("LastClosedBrowserSessionID = %q", got)
	}
	svc.mu.RLock()
	gen := svc.recordGen
	svc.mu.RUnlock()
	if gen != 4 {
		t.Fatalf("expected recordGen increment on close, got %d", gen)
	}
}

func TestResolveRecordTargetPathPrefersAppendTo(t *testing.T) {
	root := t.TempDir()
	req := RecordRequest{
		Output:   "out/new.feature",
		AppendTo: "features/existing.feature",
	}
	got := resolveRecordTargetPath(root, req)
	want := filepath.Join(root, "features", "existing.feature")
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}
