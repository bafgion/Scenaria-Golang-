package paths

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOpenWithDefaultAppRequiresPath(t *testing.T) {
	if err := OpenWithDefaultApp(""); err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestOpenWithDefaultAppMissingFile(t *testing.T) {
	if err := OpenWithDefaultApp(filepath.Join(t.TempDir(), "missing.html")); err == nil {
		t.Fatal("expected error for missing path")
	}
}

func TestOpenWithDefaultAppFile(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("skip shell open in CI")
	}
	tmp := t.TempDir()
	file := filepath.Join(tmp, "report.html")
	if err := os.WriteFile(file, []byte("<html></html>"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := OpenWithDefaultApp(file); err != nil {
		t.Fatalf("OpenWithDefaultApp: %v", err)
	}
}
