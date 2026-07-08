package gui

import (
	"context"
	"testing"
)

func TestCurrentRunSessionReturnsClonedSnapshot(t *testing.T) {
	svc := NewService()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	svc.mu.Lock()
	svc.runSession = &RunSession{
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
	svc.mu.Unlock()

	got := svc.CurrentRunSession()
	if got == nil {
		t.Fatal("expected run session")
	}
	if got.RunID != "run-1" || got.ProjectVersion != 3 {
		t.Fatalf("unexpected run session header: %+v", got)
	}

	got.RequestSnapshot.Targets[0] = "changed.feature"
	got.RequestSnapshot.Vars["A"] = "2"
	got.TempResources[0] = "changed"

	svc.mu.RLock()
	stored := svc.runSession
	svc.mu.RUnlock()
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
