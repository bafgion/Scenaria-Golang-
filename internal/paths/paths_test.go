package paths

import (
	"path/filepath"
	"testing"
)

func TestConfineToProjectRoot(t *testing.T) {
	root := t.TempDir()
	inside := filepath.Join(root, "features", "demo.feature")
	if _, err := ConfineToProjectRoot(root, "features/demo.feature"); err != nil {
		t.Fatalf("relative inside: %v", err)
	}
	got, err := ConfineToProjectRoot(root, inside)
	if err != nil {
		t.Fatalf("absolute inside: %v", err)
	}
	if got != filepath.Clean(inside) {
		t.Fatalf("got %q want %q", got, inside)
	}
	if _, err := ConfineToProjectRoot(root, "../escape.feature"); err == nil {
		t.Fatal("expected error for path outside project")
	}
}
