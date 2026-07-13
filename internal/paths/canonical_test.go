package paths

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCanonicalPathResolvesShortAndLongWindowsTempPaths(t *testing.T) {
	root := t.TempDir()
	inside := filepath.Join(root, "features", "demo.feature")
	if err := os.MkdirAll(filepath.Dir(inside), 0o755); err != nil {
		t.Fatal(err)
	}
	got, err := ConfineToProjectRoot(root, inside)
	if err != nil {
		t.Fatal(err)
	}
	if !SamePath(got, inside) {
		t.Fatalf("got %q want %q", got, inside)
	}
}

func TestSamePathTreatsEquivalentPathsAsEqual(t *testing.T) {
	root := t.TempDir()
	a := filepath.Join(root, "a.txt")
	b := filepath.Join(root, "a.txt")
	if !SamePath(a, b) {
		t.Fatalf("expected %q and %q to match", a, b)
	}
}
