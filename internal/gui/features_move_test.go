package gui

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMoveFeature(t *testing.T) {
	root := t.TempDir()
	featuresDir := filepath.Join(root, "features")
	if err := os.MkdirAll(filepath.Join(featuresDir, "nested"), 0o755); err != nil {
		t.Fatal(err)
	}
	src := filepath.Join(featuresDir, "a.feature")
	if err := os.WriteFile(src, []byte("Feature: A\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	svc := NewService()
	svc.mu.Lock()
	svc.projectPath = root
	svc.mu.Unlock()

	dest := filepath.Join(featuresDir, "nested")
	newPath, err := svc.MoveFeature(src, dest)
	if err != nil {
		t.Fatalf("MoveFeature: %v", err)
	}
	if _, err := os.Stat(newPath); err != nil {
		t.Fatalf("moved file missing: %v", err)
	}
	if _, err := os.Stat(src); !os.IsNotExist(err) {
		t.Fatal("source should be gone")
	}
}

func TestImportFeatures(t *testing.T) {
	root := t.TempDir()
	external := t.TempDir()
	extFile := filepath.Join(external, "imported.feature")
	if err := os.WriteFile(extFile, []byte("Feature: X\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	svc := NewService()
	svc.mu.Lock()
	svc.projectPath = root
	svc.mu.Unlock()

	paths, err := svc.ImportFeatures(root, []string{extFile})
	if err != nil {
		t.Fatalf("ImportFeatures: %v", err)
	}
	if len(paths) != 1 {
		t.Fatalf("expected 1 import, got %d", len(paths))
	}
	if _, err := os.Stat(paths[0]); err != nil {
		t.Fatalf("imported file missing: %v", err)
	}
}
