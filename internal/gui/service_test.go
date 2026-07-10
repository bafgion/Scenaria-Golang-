package gui

import (
	"os"
	"path/filepath"
	"testing"
	"time"
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

func TestServiceOpenProjectCleansStartupTemps(t *testing.T) {
	root := t.TempDir()
	scenaria := filepath.Join(root, ".scenaria")
	stale := filepath.Join(scenaria, "temp", "run-111")
	fresh := filepath.Join(scenaria, "temp", "run-222")
	if err := os.MkdirAll(stale, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(fresh, 0o755); err != nil {
		t.Fatal(err)
	}
	old := time.Now().Add(-48 * time.Hour)
	if err := os.Chtimes(stale, old, old); err != nil {
		t.Fatal(err)
	}

	svc := NewService()
	if _, err := svc.OpenProject(root); err != nil {
		t.Fatalf("OpenProject: %v", err)
	}
	if _, err := os.Stat(stale); !os.IsNotExist(err) {
		t.Fatalf("stale temp should be removed, err=%v", err)
	}
	if _, err := os.Stat(fresh); err != nil {
		t.Fatalf("fresh temp should remain, err=%v", err)
	}
}

func TestSearchSteps(t *testing.T) {
	svc := NewService()
	entries := svc.SearchSteps("телефон")
	if len(entries) == 0 {
		t.Fatal("expected phone step in catalog")
	}
}
