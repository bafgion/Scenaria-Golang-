package plugin

import (
	"archive/zip"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bafgion/scenaria-golang/internal/paths"
)

func TestFetchAndInstallUpdatesPluginTransactionally(t *testing.T) {
	project := t.TempDir()
	createInstalledPlugin(t, project, "demo", "old")
	if err := Install(project, "demo", "old-source"); err != nil {
		t.Fatal(err)
	}

	zipPath := filepath.Join(t.TempDir(), "new.zip")
	writePluginZipWithFiles(t, zipPath, map[string]string{
		"plugin.json": `{"id":"demo","name":"Demo"}`,
		"new.txt":     "new",
	})

	if err := FetchAndInstall(project, "demo", zipPath); err != nil {
		t.Fatalf("FetchAndInstall update: %v", err)
	}
	if _, err := os.Stat(filepath.Join(project, "addons", "demo", "old.txt")); !os.IsNotExist(err) {
		t.Fatalf("old file should be replaced, stat err=%v", err)
	}
	assertFileContent(t, filepath.Join(project, "addons", "demo", "new.txt"), "new")

	plugins, err := List(project)
	if err != nil {
		t.Fatal(err)
	}
	if len(plugins) != 1 || plugins[0].Name != "demo" || plugins[0].Source != zipPath {
		t.Fatalf("registry not updated: %#v", plugins)
	}
}

func TestFetchAndInstallMissingDescriptorLeavesExistingPlugin(t *testing.T) {
	project := t.TempDir()
	createInstalledPlugin(t, project, "demo", "old")
	if err := Install(project, "demo", "old-source"); err != nil {
		t.Fatal(err)
	}

	zipPath := filepath.Join(t.TempDir(), "broken.zip")
	writePluginZipWithFiles(t, zipPath, map[string]string{"readme.txt": "no descriptor"})

	err := FetchAndInstall(project, "demo", zipPath)
	if err == nil || !strings.Contains(err.Error(), "plugin.json is required") {
		t.Fatalf("expected missing descriptor error, got %v", err)
	}
	assertFileContent(t, filepath.Join(project, "addons", "demo", "old.txt"), "old")
	assertRegistrySource(t, project, "demo", "old-source")
	assertNoStagingDirs(t, project)
}

func TestFetchAndInstallDescriptorIDMismatchLeavesExistingPlugin(t *testing.T) {
	project := t.TempDir()
	createInstalledPlugin(t, project, "demo", "old")
	if err := Install(project, "demo", "old-source"); err != nil {
		t.Fatal(err)
	}

	zipPath := filepath.Join(t.TempDir(), "mismatch.zip")
	writePluginZipWithFiles(t, zipPath, map[string]string{
		"plugin.json": `{"id":"other"}`,
		"new.txt":     "new",
	})

	err := FetchAndInstall(project, "demo", zipPath)
	if err == nil || !strings.Contains(err.Error(), "does not match") {
		t.Fatalf("expected descriptor mismatch error, got %v", err)
	}
	assertFileContent(t, filepath.Join(project, "addons", "demo", "old.txt"), "old")
	assertRegistrySource(t, project, "demo", "old-source")
}

func TestFetchAndInstallCommitRenameFailureRestoresExistingPlugin(t *testing.T) {
	project := t.TempDir()
	createInstalledPlugin(t, project, "demo", "old")
	if err := Install(project, "demo", "old-source"); err != nil {
		t.Fatal(err)
	}
	dest, err := addonPath(project, "demo")
	if err != nil {
		t.Fatal(err)
	}

	zipPath := filepath.Join(t.TempDir(), "new.zip")
	writePluginZipWithFiles(t, zipPath, map[string]string{
		"plugin.json": `{"id":"demo"}`,
		"new.txt":     "new",
	})

	origRename := installRename
	defer func() { installRename = origRename }()
	installRename = func(oldpath, newpath string) error {
		stagingMarker := filepath.Join(".scenaria", "plugin-staging")
		if paths.SamePath(newpath, dest) && strings.Contains(filepath.ToSlash(filepath.Clean(oldpath)), filepath.ToSlash(stagingMarker)) {
			return errors.New("simulated final rename failure")
		}
		return os.Rename(oldpath, newpath)
	}

	err = FetchAndInstall(project, "demo", zipPath)
	if err == nil || !strings.Contains(err.Error(), "commit plugin files") {
		t.Fatalf("expected commit failure, got %v", err)
	}
	assertFileContent(t, filepath.Join(project, "addons", "demo", "old.txt"), "old")
	if _, err := os.Stat(filepath.Join(project, "addons", "demo", "new.txt")); !os.IsNotExist(err) {
		t.Fatalf("new file should not remain after rollback, stat err=%v", err)
	}
	assertRegistrySource(t, project, "demo", "old-source")
}

func TestFetchAndInstallRegistryFailureRestoresExistingPlugin(t *testing.T) {
	project := t.TempDir()
	createInstalledPlugin(t, project, "demo", "old")
	if err := Install(project, "demo", "old-source"); err != nil {
		t.Fatal(err)
	}
	manifest := ManifestPath(project)
	if err := os.WriteFile(manifest, []byte(`{broken json`), 0o644); err != nil {
		t.Fatal(err)
	}

	zipPath := filepath.Join(t.TempDir(), "new.zip")
	writePluginZipWithFiles(t, zipPath, map[string]string{
		"plugin.json": `{"id":"demo"}`,
		"new.txt":     "new",
	})

	err := FetchAndInstall(project, "demo", zipPath)
	if err == nil || !strings.Contains(err.Error(), "decode plugin manifest") {
		t.Fatalf("expected registry decode error, got %v", err)
	}
	assertFileContent(t, filepath.Join(project, "addons", "demo", "old.txt"), "old")
	if _, err := os.Stat(filepath.Join(project, "addons", "demo", "new.txt")); !os.IsNotExist(err) {
		t.Fatalf("new file should not remain after registry rollback, stat err=%v", err)
	}
}

func TestFetchAndInstallMissingSourceLeavesExistingPluginAndCleansStaging(t *testing.T) {
	project := t.TempDir()
	createInstalledPlugin(t, project, "demo", "old")

	err := FetchAndInstall(project, "demo", filepath.Join(t.TempDir(), "missing"))
	if err == nil || !strings.Contains(err.Error(), "plugin source not found") {
		t.Fatalf("expected missing source error, got %v", err)
	}
	assertFileContent(t, filepath.Join(project, "addons", "demo", "old.txt"), "old")
	assertNoStagingDirs(t, project)
}

func createInstalledPlugin(t *testing.T, project, name, marker string) {
	t.Helper()
	dir := filepath.Join(project, "addons", name)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "plugin.json"), []byte(`{"id":"`+name+`"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, marker+".txt"), []byte(marker), 0o644); err != nil {
		t.Fatal(err)
	}
}

func writePluginZipWithFiles(t *testing.T, zipPath string, files map[string]string) {
	t.Helper()
	file, err := os.Create(zipPath)
	if err != nil {
		t.Fatal(err)
	}
	w := zip.NewWriter(file)
	for name, content := range files {
		entry, err := w.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := entry.Write([]byte(content)); err != nil {
			t.Fatal(err)
		}
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
}

func assertFileContent(t *testing.T, path, want string) {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	if string(raw) != want {
		t.Fatalf("%s = %q, want %q", path, string(raw), want)
	}
}

func assertRegistrySource(t *testing.T, project, name, want string) {
	t.Helper()
	plugins, err := List(project)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range plugins {
		if entry.Name == name {
			if entry.Source != want {
				t.Fatalf("registry source = %q, want %q", entry.Source, want)
			}
			return
		}
	}
	t.Fatalf("plugin %q not found in registry: %#v", name, plugins)
}

func assertNoStagingDirs(t *testing.T, project string) {
	t.Helper()
	stagingRoot := filepath.Join(project, ".scenaria", "plugin-staging")
	entries, err := os.ReadDir(stagingRoot)
	if os.IsNotExist(err) {
		return
	}
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected no staging dirs, got %v", entries)
	}
}
