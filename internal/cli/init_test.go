package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunInit(t *testing.T) {
	tmp := t.TempDir()
	if err := RunInit([]string{tmp}); err != nil {
		t.Fatalf("RunInit failed: %v", err)
	}
	if _, err := os.Stat(filepath.Join(tmp, ".scenaria", "project.json")); err != nil {
		t.Fatalf("project.json missing: %v", err)
	}
}
