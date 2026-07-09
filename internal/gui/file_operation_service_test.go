package gui

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileOperationServiceReadAndSaveFeature(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "demo.feature")
	if err := os.WriteFile(path, []byte("Функционал: Demo"), 0o644); err != nil {
		t.Fatal(err)
	}
	svc := NewFileOperationService(
		func(_ string) (string, error) { return path, nil },
		func() string { return root },
		func(run func() error) error { return run() },
		func(run func() error) error { return run() },
	)

	got, err := svc.ReadFeature(path)
	if err != nil {
		t.Fatalf("ReadFeature: %v", err)
	}
	if got != "Функционал: Demo" {
		t.Fatalf("unexpected content: %q", got)
	}
	if err := svc.SaveFeature(path, "Функционал: Updated"); err != nil {
		t.Fatalf("SaveFeature: %v", err)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(raw) != "Функционал: Updated" {
		t.Fatalf("unexpected saved content: %q", string(raw))
	}
}

func TestFileOperationServiceDeleteFeature(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "demo.feature")
	if err := os.WriteFile(path, []byte("Feature:"), 0o644); err != nil {
		t.Fatal(err)
	}
	svc := NewFileOperationService(
		func(p string) (string, error) { return p, nil },
		func() string { return root },
		func(run func() error) error { return run() },
		func(run func() error) error { return run() },
	)
	if err := svc.DeleteFeature(path); err != nil {
		t.Fatalf("DeleteFeature: %v", err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("expected file removed, stat err=%v", err)
	}
}
