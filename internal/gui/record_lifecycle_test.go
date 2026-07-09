package gui

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/bafgion/scenaria-golang/internal/recorder"
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
	svc := NewService()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	begin := svc.recorderOps().Session().BeginSession(BeginRecorderSessionInput{
		Session:      recorder.NewLiveSession(),
		RecordCtx:    ctx,
		RecordCancel: cancel,
		TargetPath:   "features/a.feature",
	})
	emitted := 0
	emit := svc.guardedRecordEmit(begin.Gen, "features/a.feature", func(string, any) {
		emitted++
	})
	emit("record-step", map[string]any{"index": 0, "line": "step"})
	if emitted != 1 {
		t.Fatalf("expected 1 emit, got %d", emitted)
	}
	ctx2, cancel2 := context.WithCancel(context.Background())
	defer cancel2()
	svc.recorderOps().Session().BeginSession(BeginRecorderSessionInput{
		Session:      recorder.NewLiveSession(),
		RecordCtx:    ctx2,
		RecordCancel: cancel2,
	})
	emit("record-step", map[string]any{"index": 1, "line": "late"})
	if emitted != 1 {
		t.Fatalf("expected stale emit to be skipped, got %d", emitted)
	}
}

func TestCloseBrowserRetiresRecordGeneration(t *testing.T) {
	svc, session := attachLiveSession(t, true)
	snap := svc.recorderOps().Session().Snapshot()

	emitted := 0
	emit := svc.guardedRecordEmit(snap.Gen, snap.TargetPath, func(string, any) { emitted++ })

	svc.CloseBrowser()
	emit("record-step", map[string]any{"index": 0, "line": "late"})
	if emitted != 0 {
		t.Fatalf("expected retired generation to block emit, got %d", emitted)
	}
	if session.BrowserAlive() {
		t.Fatal("expected session cleared")
	}
	if got := svc.LastClosedBrowserSessionID(); got != snap.BrowserSessionID {
		t.Fatalf("LastClosedBrowserSessionID = %q", got)
	}
	if gen := svc.recorderOps().Session().Generation(); gen != snap.Gen+1 {
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
