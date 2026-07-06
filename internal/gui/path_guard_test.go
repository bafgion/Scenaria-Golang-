package gui

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfineFeaturePathInsideProject(t *testing.T) {
	root := t.TempDir()
	svc := NewService()
	if _, err := svc.OpenProject(root); err != nil {
		t.Fatal(err)
	}
	inside := filepath.Join(root, "features", "a.feature")
	if err := os.MkdirAll(filepath.Dir(inside), 0o755); err != nil {
		t.Fatal(err)
	}
	got, err := svc.confineFeaturePath(inside)
	if err != nil {
		t.Fatal(err)
	}
	if got != inside {
		t.Fatalf("got %q", got)
	}
}

func TestConfineFeaturePathRejectsOutsideProject(t *testing.T) {
	root := t.TempDir()
	svc := NewService()
	if _, err := svc.OpenProject(root); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.confineFeaturePath(filepath.Join(os.TempDir(), "escape.feature")); err == nil {
		t.Fatal("expected error for path outside project")
	}
}

func TestReadFeatureConfined(t *testing.T) {
	root := t.TempDir()
	svc := NewService()
	if _, err := svc.OpenProject(root); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, "demo.feature")
	if err := os.WriteFile(path, []byte("Функционал: x"), 0o644); err != nil {
		t.Fatal(err)
	}
	text, err := svc.ReadFeature(path)
	if err != nil {
		t.Fatal(err)
	}
	if text == "" {
		t.Fatal("empty read")
	}
}
