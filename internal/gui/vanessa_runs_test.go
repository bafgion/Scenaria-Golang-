package gui

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestListVanessaRunDirs(t *testing.T) {
	root := t.TempDir()
	runsRoot := filepath.Join(root, "runs")
	old := os.Getenv("XDG_CONFIG_HOME")
	t.Setenv("XDG_CONFIG_HOME", root)
	defer os.Setenv("XDG_CONFIG_HOME", old)

	cfgDir := filepath.Join(root, "Scenaria")
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		t.Fatal(err)
	}
	vanessaJSON := filepath.Join(root, ".scenaria", "vanessa.json")
	if err := os.MkdirAll(filepath.Dir(vanessaJSON), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(vanessaJSON, []byte(`{"runs_dir":"`+filepath.ToSlash(runsRoot)+`"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	run1 := filepath.Join(runsRoot, "run-1")
	run2 := filepath.Join(runsRoot, "run-2")
	for _, dir := range []string{run1, run2} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	time.Sleep(10 * time.Millisecond)
	if err := os.Chtimes(run1, time.Now().Add(-time.Hour), time.Now().Add(-time.Hour)); err != nil {
		t.Fatal(err)
	}

	svc := NewService()
	if _, err := svc.OpenProject(root); err != nil {
		t.Fatal(err)
	}
	dirs, err := svc.ListVanessaRunDirs(10)
	if err != nil {
		t.Fatalf("ListVanessaRunDirs: %v", err)
	}
	if len(dirs) != 2 {
		t.Fatalf("dirs = %v", dirs)
	}
	if dirs[0] != run2 {
		t.Fatalf("newest first: got %q want %q", dirs[0], run2)
	}
}
