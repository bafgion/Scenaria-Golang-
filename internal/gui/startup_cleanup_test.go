package gui

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCleanupStartupTempArtifactsRemovesStaleProjectAndGlobalTemps(t *testing.T) {
	root := t.TempDir()
	scenaria := filepath.Join(root, ".scenaria")
	projectTemp := filepath.Join(scenaria, "temp")
	runsDir := filepath.Join(scenaria, "runs", "run-1")
	staleProjectTemp := filepath.Join(projectTemp, "run-111")
	freshProjectTemp := filepath.Join(projectTemp, "run-222")
	staleRunTemp := filepath.Join(runsDir, ".scenaria-write-stale")
	freshRunTemp := filepath.Join(runsDir, ".scenaria-write-fresh")
	if err := os.MkdirAll(staleProjectTemp, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(freshProjectTemp, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(runsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(staleProjectTemp, "scenario.feature"), []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(freshProjectTemp, "scenario.feature"), []byte("new"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(staleRunTemp, []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(freshRunTemp, []byte("new"), 0o644); err != nil {
		t.Fatal(err)
	}
	old := time.Now().Add(-48 * time.Hour)
	if err := os.Chtimes(staleProjectTemp, old, old); err != nil {
		t.Fatal(err)
	}
	if err := os.Chtimes(staleRunTemp, old, old); err != nil {
		t.Fatal(err)
	}

	globalRoot := t.TempDir()
	staleGlobal := filepath.Join(globalRoot, "scenaria-run-stale")
	freshGlobal := filepath.Join(globalRoot, "scenaria-run-fresh")
	if err := os.MkdirAll(staleGlobal, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(freshGlobal, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Chtimes(staleGlobal, old, old); err != nil {
		t.Fatal(err)
	}

	if err := cleanupStartupTempArtifactsForProject(root, time.Hour, time.Now()); err != nil {
		t.Fatalf("project cleanup: %v", err)
	}
	if err := cleanupGlobalStartupTemps(globalRoot, time.Hour, time.Now()); err != nil {
		t.Fatalf("global cleanup: %v", err)
	}

	if _, err := os.Stat(staleProjectTemp); !os.IsNotExist(err) {
		t.Fatalf("stale project temp should be removed, err=%v", err)
	}
	if _, err := os.Stat(freshProjectTemp); err != nil {
		t.Fatalf("fresh project temp should remain, err=%v", err)
	}
	if _, err := os.Stat(staleRunTemp); !os.IsNotExist(err) {
		t.Fatalf("stale run temp should be removed, err=%v", err)
	}
	if _, err := os.Stat(freshRunTemp); err != nil {
		t.Fatalf("fresh run temp should remain, err=%v", err)
	}
	if _, err := os.Stat(staleGlobal); !os.IsNotExist(err) {
		t.Fatalf("stale global temp should be removed, err=%v", err)
	}
	if _, err := os.Stat(freshGlobal); err != nil {
		t.Fatalf("fresh global temp should remain, err=%v", err)
	}
}
