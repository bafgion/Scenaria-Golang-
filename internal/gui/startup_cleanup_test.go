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

func TestCleanupStartupTempArtifactsRemovesStalePluginInstallTemps(t *testing.T) {
	root := t.TempDir()
	scenaria := filepath.Join(root, ".scenaria")
	staging := filepath.Join(scenaria, "plugin-staging")
	backups := filepath.Join(scenaria, "plugin-backups")
	staleStaging := filepath.Join(staging, "demo-111")
	freshStaging := filepath.Join(staging, "demo-222")
	staleBackup := filepath.Join(backups, "demo-333")
	unknown := filepath.Join(staging, "notowned")
	for _, dir := range []string{staleStaging, freshStaging, staleBackup, unknown} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	old := time.Now().Add(-48 * time.Hour)
	for _, dir := range []string{staleStaging, staleBackup, unknown} {
		if err := os.Chtimes(dir, old, old); err != nil {
			t.Fatal(err)
		}
	}

	if err := cleanupStartupTempArtifactsForProject(root, time.Hour, time.Now()); err != nil {
		t.Fatalf("project cleanup: %v", err)
	}

	if _, err := os.Stat(staleStaging); !os.IsNotExist(err) {
		t.Fatalf("stale plugin staging should be removed, err=%v", err)
	}
	if _, err := os.Stat(staleBackup); !os.IsNotExist(err) {
		t.Fatalf("stale plugin backup should be removed, err=%v", err)
	}
	if _, err := os.Stat(freshStaging); err != nil {
		t.Fatalf("fresh plugin staging should remain, err=%v", err)
	}
	if _, err := os.Stat(unknown); err != nil {
		t.Fatalf("unknown plugin temp name should remain, err=%v", err)
	}
}

func TestCleanupStartupTempArtifactsSkipsSymlinkEscape(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	outsideFile := filepath.Join(outside, "keep.txt")
	if err := os.WriteFile(outsideFile, []byte("keep"), 0o644); err != nil {
		t.Fatal(err)
	}
	staging := filepath.Join(root, ".scenaria", "plugin-staging")
	if err := os.MkdirAll(staging, 0o755); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(staging, "demo-escape")
	createGUIDirSymlinkOrSkip(t, outside, link)
	old := time.Now().Add(-48 * time.Hour)
	if err := os.Chtimes(outside, old, old); err != nil {
		t.Fatal(err)
	}

	if err := cleanupStartupTempArtifactsForProject(root, time.Hour, time.Now()); err != nil {
		t.Fatalf("project cleanup: %v", err)
	}

	if _, err := os.Stat(outsideFile); err != nil {
		t.Fatalf("outside data must remain, err=%v", err)
	}
	if _, err := os.Lstat(link); err != nil {
		t.Fatalf("unsafe symlink should be skipped rather than followed/removed, err=%v", err)
	}
}
