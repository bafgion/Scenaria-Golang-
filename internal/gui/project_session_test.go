package gui

import (
	"context"
	"testing"
	"time"

	"github.com/bafgion/scenaria-golang/internal/recorder"
)

func TestOpenProjectIncrementsProjectVersionAndCancelsPreviousSession(t *testing.T) {
	svc := NewService()
	first := t.TempDir()
	if _, err := svc.OpenProject(first); err != nil {
		t.Fatalf("OpenProject first: %v", err)
	}

	svc.mu.RLock()
	firstSession := svc.projectSession
	firstVersion := svc.projectVersion
	svc.mu.RUnlock()
	if firstSession == nil {
		t.Fatal("expected first project session")
	}
	if firstVersion == 0 {
		t.Fatal("expected non-zero project version")
	}

	second := t.TempDir()
	if _, err := svc.OpenProject(second); err != nil {
		t.Fatalf("OpenProject second: %v", err)
	}

	svc.mu.RLock()
	secondSession := svc.projectSession
	secondVersion := svc.projectVersion
	svc.mu.RUnlock()
	if secondSession == nil {
		t.Fatal("expected second project session")
	}
	if secondVersion <= firstVersion {
		t.Fatalf("expected project version increment: first=%d second=%d", firstVersion, secondVersion)
	}
	if secondSession.Root != second {
		t.Fatalf("unexpected session root: %q", secondSession.Root)
	}

	select {
	case <-firstSession.Context.Done():
	case <-time.After(200 * time.Millisecond):
		t.Fatal("expected previous project session context to be canceled")
	}
}

func TestOpenProjectCancelsRunValidateAndRecordContexts(t *testing.T) {
	svc := NewService()
	root := t.TempDir()

	begin, err := svc.runner().TryBegin(1, RunRequest{Targets: []string{"a.feature"}}, nil, DefaultRunTimeout)
	if err != nil {
		t.Fatalf("TryBegin: %v", err)
	}
	runCtx := begin.Context

	validateCtx, validateCancel := context.WithCancel(context.Background())
	recordCtx, recordCancel := context.WithCancel(context.Background())

	svc.mu.Lock()
	svc.validateCancel = validateCancel
	svc.mu.Unlock()
	svc.recorderOps().Session().BeginSession(BeginRecorderSessionInput{
		Session:      recorder.NewLiveSession(),
		RecordCtx:    recordCtx,
		RecordCancel: recordCancel,
	})

	if _, err := svc.OpenProject(root); err != nil {
		t.Fatalf("OpenProject: %v", err)
	}

	select {
	case <-runCtx.Done():
	case <-time.After(200 * time.Millisecond):
		t.Fatal("expected run context canceled")
	}
	select {
	case <-validateCtx.Done():
	case <-time.After(200 * time.Millisecond):
		t.Fatal("expected validate context canceled")
	}
	select {
	case <-recordCtx.Done():
	case <-time.After(200 * time.Millisecond):
		t.Fatal("expected record context canceled")
	}
}
