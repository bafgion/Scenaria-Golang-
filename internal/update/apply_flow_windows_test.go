//go:build windows

package update

import (
	"archive/zip"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bafgion/scenaria-golang/internal/brand"
)

func TestApplyDownloadedPortableStagesAndScript(t *testing.T) {
	installDir := t.TempDir()
	zipPath := filepath.Join(t.TempDir(), "update.zip")
	writePortableZip(t, zipPath, map[string]string{
		"Scenaria/" + brand.GUIExeName: "gui",
		"Scenaria/readme.txt":          "v2",
	})

	var launchedScript string
	oldHook := launchHiddenBatchHook
	launchHiddenBatchHook = func(scriptPath string) error {
		launchedScript = scriptPath
		return nil
	}
	defer func() { launchHiddenBatchHook = oldHook }()

	if err := ApplyDownloaded(zipPath, ApplyKindPortable, installDir, 4242, nil); err != nil {
		t.Fatalf("ApplyDownloaded: %v", err)
	}
	staging := filepath.Join(installDir, updateStagingDir)
	if _, err := os.Stat(filepath.Join(staging, brand.GUIExeName)); err != nil {
		t.Fatalf("staged gui exe: %v", err)
	}
	if _, err := os.Stat(filepath.Join(staging, "readme.txt")); err != nil {
		t.Fatalf("staged readme: %v", err)
	}
	scriptPath := filepath.Join(installDir, updateBatName)
	if launchedScript != scriptPath {
		t.Fatalf("launched=%q want %q", launchedScript, scriptPath)
	}
	body, err := os.ReadFile(scriptPath)
	if err != nil {
		t.Fatal(err)
	}
	script := string(body)
	if !strings.Contains(script, "robocopy") {
		t.Fatalf("missing robocopy: %s", script)
	}
	if !strings.Contains(script, "cd /d") {
		t.Fatalf("expected working dir change in script")
	}
	if !strings.Contains(script, "/D") {
		t.Fatalf("expected start /D in script")
	}
	if !strings.Contains(script, "4242") {
		t.Fatalf("missing parent pid in script")
	}
}

func TestApplyDownloadedRejectsEmptyPaths(t *testing.T) {
	if err := ApplyDownloaded("", ApplyKindPortable, t.TempDir(), 1, nil); err == nil {
		t.Fatal("expected error for empty asset path")
	}
	if err := ApplyDownloaded("x.zip", ApplyKindPortable, "", 1, nil); err == nil {
		t.Fatal("expected error for empty install dir")
	}
}

func TestApplyDownloadedUnsupportedKind(t *testing.T) {
	err := ApplyDownloaded("file.exe", ApplyKind("msi"), t.TempDir(), 1, nil)
	if err == nil || !strings.Contains(err.Error(), "unsupported") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func writePortableZip(t *testing.T, zipPath string, files map[string]string) {
	t.Helper()
	out, err := os.Create(zipPath)
	if err != nil {
		t.Fatal(err)
	}
	w := zip.NewWriter(out)
	for name, content := range files {
		f, err := w.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := f.Write([]byte(content)); err != nil {
			t.Fatal(err)
		}
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	if err := out.Close(); err != nil {
		t.Fatal(err)
	}
}
