package gui

import (
	"os"
	"path/filepath"
	"testing"
)

func TestProjectServiceProjectInfo(t *testing.T) {
	root := t.TempDir()
	featurePath := filepath.Join(root, "demo.feature")
	if err := os.WriteFile(featurePath, []byte("Функционал: Demo\n@smoke\nСценарий: A\n  Допустим открыт \"https://example.com\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	svc := NewProjectService(func(run func() error) error { return run() })
	info, err := svc.ProjectInfo(root, 7)
	if err != nil {
		t.Fatalf("ProjectInfo: %v", err)
	}
	if info.Path != root {
		t.Fatalf("unexpected path: %q", info.Path)
	}
	if info.Version != 7 {
		t.Fatalf("unexpected version: %d", info.Version)
	}
	if len(info.Features) != 1 || info.Features[0] != featurePath {
		t.Fatalf("unexpected features: %#v", info.Features)
	}
}

func TestProjectServiceProjectInfoRequiresPath(t *testing.T) {
	svc := NewProjectService(func(run func() error) error { return run() })
	if _, err := svc.ProjectInfo("", 1); err == nil {
		t.Fatal("expected error for empty path")
	}
}
