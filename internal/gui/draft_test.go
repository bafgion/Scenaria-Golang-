package gui

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFeatureDraftRoundTrip(t *testing.T) {
	root := t.TempDir()
	feature := filepath.Join(root, "demo.feature")
	if err := os.WriteFile(feature, []byte("Функционал: X\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	svc := NewService()
	svc.projectPath = root
	content := "Функционал: X\n  Сценарий: Y\n"
	if err := svc.SaveFeatureDraft(feature, content); err != nil {
		t.Fatal(err)
	}
	loaded, err := svc.LoadFeatureDraft(feature)
	if err != nil {
		t.Fatal(err)
	}
	if loaded != content {
		t.Fatalf("draft mismatch: %q", loaded)
	}
	if err := svc.ClearFeatureDraft(feature); err != nil {
		t.Fatal(err)
	}
	again, err := svc.LoadFeatureDraft(feature)
	if err != nil {
		t.Fatal(err)
	}
	if again != "" {
		t.Fatalf("expected empty draft, got %q", again)
	}
}
