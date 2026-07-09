package gui

import (
	"context"
	"testing"
)

func TestCurrentRunSessionReturnsClonedSnapshot(t *testing.T) {
	runSvc := NewRunService(func() string { return "" })
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	runSvc.mu.Lock()
	runSvc.session = &RunSession{
		RunID:          "run-1",
		ProjectVersion: 3,
		RequestSnapshot: RunRequest{
			Scenario: "Login",
			Targets:  []string{"a.feature"},
			Vars:     map[string]string{"A": "1"},
		},
		Context:       ctx,
		TempResources: []string{"tmp/a.feature"},
	}
	runSvc.mu.Unlock()

	got := runSvc.CurrentSession()
	if got == nil {
		t.Fatal("expected run session")
	}
	if got.RunID != "run-1" || got.ProjectVersion != 3 {
		t.Fatalf("unexpected run session header: %+v", got)
	}

	got.RequestSnapshot.Targets[0] = "changed.feature"
	got.RequestSnapshot.Vars["A"] = "2"
	got.TempResources[0] = "changed"

	runSvc.mu.RLock()
	stored := runSvc.session
	runSvc.mu.RUnlock()
	if stored == nil {
		t.Fatal("stored run session disappeared")
	}
	if stored.RequestSnapshot.Targets[0] != "a.feature" {
		t.Fatalf("targets were not cloned: %#v", stored.RequestSnapshot.Targets)
	}
	if stored.RequestSnapshot.Vars["A"] != "1" {
		t.Fatalf("vars were not cloned: %#v", stored.RequestSnapshot.Vars)
	}
	if stored.TempResources[0] != "tmp/a.feature" {
		t.Fatalf("temp resources were not cloned: %#v", stored.TempResources)
	}
}

func TestRunServiceTryBeginRejectsOverlappingRuns(t *testing.T) {
	runSvc := NewRunService(func() string { return "" })
	first, err := runSvc.TryBegin(1, RunRequest{Targets: []string{"a.feature"}}, nil, DefaultRunTimeout)
	if err != nil {
		t.Fatalf("TryBegin: %v", err)
	}
	defer first.Cancel()
	defer runSvc.Finish(first.Gen)

	_, err = runSvc.TryBegin(1, RunRequest{Targets: []string{"b.feature"}}, nil, DefaultRunTimeout)
	if err == nil {
		t.Fatal("expected overlapping run rejection")
	}
}

func TestServiceCurrentRunSessionDelegatesToRunService(t *testing.T) {
	svc := NewService()

	begin, err := svc.runner().TryBegin(2, RunRequest{Scenario: "S"}, nil, DefaultRunTimeout)
	if err != nil {
		t.Fatalf("TryBegin: %v", err)
	}
	defer begin.Cancel()
	defer svc.runner().Finish(begin.Gen)

	got := svc.CurrentRunSession()
	if got == nil || got.RunID != begin.RunID {
		t.Fatalf("unexpected delegated session: %+v", got)
	}
}
