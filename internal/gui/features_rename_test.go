package gui

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRenameFeature(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "old.feature")
	if err := os.WriteFile(src, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	svc := NewService()
	if _, err := svc.OpenProject(root); err != nil {
		t.Fatal(err)
	}
	newPath, err := svc.RenameFeature(src, "new.feature")
	if err != nil {
		t.Fatalf("RenameFeature: %v", err)
	}
	if _, err := os.Stat(newPath); err != nil {
		t.Fatalf("expected renamed file: %v", err)
	}
	if _, err := os.Stat(src); !os.IsNotExist(err) {
		t.Fatalf("old file should be gone")
	}
}

func TestRenameFeatureAddsExtension(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "demo.feature")
	if err := os.WriteFile(src, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	svc := NewService()
	if _, err := svc.OpenProject(root); err != nil {
		t.Fatal(err)
	}
	newPath, err := svc.RenameFeature(src, "renamed")
	if err != nil {
		t.Fatalf("RenameFeature: %v", err)
	}
	if filepath.Base(newPath) != "renamed.feature" {
		t.Fatalf("got %q", newPath)
	}
}
