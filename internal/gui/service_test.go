package gui

import (
	"testing"
)

func TestServiceOpenProject(t *testing.T) {
	root := t.TempDir()
	svc := NewService()
	if _, err := svc.OpenProject(root); err != nil {
		t.Fatalf("OpenProject: %v", err)
	}
	if svc.ProjectPath() != root {
		t.Fatalf("unexpected path: %q", svc.ProjectPath())
	}
}

func TestServiceRefreshProject(t *testing.T) {
	root := t.TempDir()
	svc := NewService()
	if _, err := svc.OpenProject(root); err != nil {
		t.Fatalf("OpenProject: %v", err)
	}
	info, err := svc.RefreshProject()
	if err != nil {
		t.Fatalf("RefreshProject: %v", err)
	}
	if info.Path != root {
		t.Fatalf("unexpected path: %q", info.Path)
	}
}

func TestSearchSteps(t *testing.T) {
	svc := NewService()
	entries := svc.SearchSteps("телефон")
	if len(entries) == 0 {
		t.Fatal("expected phone step in catalog")
	}
}
