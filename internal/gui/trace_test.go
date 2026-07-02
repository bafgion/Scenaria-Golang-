package gui

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOpenTrace_MissingDir(t *testing.T) {
	svc := NewService()
	result := svc.OpenTrace("")
	if result.Error == "" {
		t.Fatal("expected error without project")
	}
}

func TestNewestTraceZip(t *testing.T) {
	dir := t.TempDir()
	old := filepath.Join(dir, "a.zip")
	new := filepath.Join(dir, "b.zip")
	if err := os.WriteFile(old, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(new, []byte("y"), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := newestTraceZip(dir)
	if err != nil {
		t.Fatal(err)
	}
	if got != new && got != old {
		t.Fatalf("unexpected trace: %s", got)
	}
}
