package paths

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfineToProjectRoot(t *testing.T) {
	root := t.TempDir()
	inside := filepath.Join(root, "features", "demo.feature")
	if err := os.MkdirAll(filepath.Dir(inside), 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := ConfineToProjectRoot(root, "features/demo.feature"); err != nil {
		t.Fatalf("relative inside: %v", err)
	}
	got, err := ConfineToProjectRoot(root, inside)
	if err != nil {
		t.Fatalf("absolute inside: %v", err)
	}
	if !SamePath(got, inside) {
		t.Fatalf("got %q want %q", got, inside)
	}
	if _, err := ConfineToProjectRoot(root, "../escape.feature"); err == nil {
		t.Fatal("expected error for path outside project")
	}
}

func TestConfineToProjectRootAllowsDotDotPrefixNames(t *testing.T) {
	root := t.TempDir()
	got, err := ConfineToProjectRoot(root, "..hidden.feature")
	if err != nil {
		t.Fatalf("dot-dot prefix file name should remain inside project: %v", err)
	}
	want := filepath.Join(root, "..hidden.feature")
	if !SamePath(got, want) {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestConfineToProjectRootAllowsMissingDescendantParents(t *testing.T) {
	root := t.TempDir()
	got, err := ConfineToProjectRoot(root, filepath.Join("reports", "run-1", "trace.zip"))
	if err != nil {
		t.Fatalf("missing descendant parents should be allowed for new paths: %v", err)
	}
	want := filepath.Join(root, "reports", "run-1", "trace.zip")
	if !SamePath(got, want) {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestConfineToProjectRootRejectsSymlinkEscape(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	if err := os.WriteFile(filepath.Join(outside, "secret.feature"), []byte("Feature: Secret\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(root, "linked")
	createDirSymlinkOrSkip(t, outside, link)

	if _, err := ConfineToProjectRoot(root, filepath.Join("linked", "secret.feature")); err == nil {
		t.Fatal("expected existing file through symlink escape to fail")
	}
	if _, err := ConfineToProjectRoot(root, filepath.Join("linked", "new.feature")); err == nil {
		t.Fatal("expected new file through symlink escape to fail")
	}
}

func TestConfineToProjectRootRejectsBrokenSymlinkLeaf(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	link := filepath.Join(root, "broken.feature")
	createDirSymlinkOrSkip(t, filepath.Join(outside, "missing.feature"), link)

	if _, err := ConfineToProjectRoot(root, "broken.feature"); err == nil {
		t.Fatal("expected broken symlink leaf to fail instead of being treated as a new file")
	}
}

func createDirSymlinkOrSkip(t *testing.T, target, link string) {
	t.Helper()
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("symlink not available: %v", err)
	}
}
