package plugin

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInstallListUninstall(t *testing.T) {
	tmp := t.TempDir()

	if err := Install(tmp, "vanessa", "https://example.com/vanessa.zip"); err != nil {
		t.Fatalf("Install returned error: %v", err)
	}
	plugins, err := List(tmp)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(plugins) != 1 || plugins[0].Name != "vanessa" {
		t.Fatalf("unexpected plugins list: %+v", plugins)
	}

	removed, err := Uninstall(tmp, "vanessa")
	if err != nil {
		t.Fatalf("Uninstall returned error: %v", err)
	}
	if !removed {
		t.Fatal("expected plugin to be removed")
	}

	plugins, err = List(tmp)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(plugins) != 0 {
		t.Fatalf("expected empty plugin list, got: %+v", plugins)
	}
}

func TestSaveManifestAtomicTempWriteFailurePreservesOldRegistry(t *testing.T) {
	tmp := t.TempDir()
	if err := SaveManifest(tmp, Manifest{Plugins: []Entry{{Name: "demo", Source: "old"}}}); err != nil {
		t.Fatal(err)
	}
	origWrite := registryWriteTemp
	defer func() { registryWriteTemp = origWrite }()
	registryWriteTemp = func(*os.File, []byte) (int, error) {
		return 0, errors.New("disk full")
	}

	err := SaveManifest(tmp, Manifest{Plugins: []Entry{{Name: "demo", Source: "new"}}})
	if err == nil || !strings.Contains(err.Error(), "disk full") {
		t.Fatalf("expected disk full error, got %v", err)
	}
	plugins, err := List(tmp)
	if err != nil {
		t.Fatal(err)
	}
	if len(plugins) != 1 || plugins[0].Source != "old" {
		t.Fatalf("old registry should be preserved, got %#v", plugins)
	}
}

func TestSaveManifestAtomicReplaceFailurePreservesOldRegistry(t *testing.T) {
	tmp := t.TempDir()
	if err := SaveManifest(tmp, Manifest{Plugins: []Entry{{Name: "demo", Source: "old"}}}); err != nil {
		t.Fatal(err)
	}
	origRename := registryRename
	defer func() { registryRename = origRename }()
	registryRename = func(oldpath, newpath string) error {
		if strings.Contains(filepath.Base(oldpath), ".tmp-") && filepath.Base(newpath) == "plugins.json" {
			return errors.New("replace locked")
		}
		return os.Rename(oldpath, newpath)
	}

	err := SaveManifest(tmp, Manifest{Plugins: []Entry{{Name: "demo", Source: "new"}}})
	if err == nil || !strings.Contains(err.Error(), "replace locked") {
		t.Fatalf("expected replace locked error, got %v", err)
	}
	plugins, err := List(tmp)
	if err != nil {
		t.Fatal(err)
	}
	if len(plugins) != 1 || plugins[0].Source != "old" {
		t.Fatalf("old registry should be restored, got %#v", plugins)
	}
}

func TestLoadManifestReportsMalformedExistingRegistry(t *testing.T) {
	tmp := t.TempDir()
	path := ManifestPath(tmp)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(`{broken json`), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadManifest(tmp); err == nil || !strings.Contains(err.Error(), "decode plugin manifest") {
		t.Fatalf("expected decode plugin manifest error, got %v", err)
	}
}

func TestUninstallRemovesRegistryAndPluginFiles(t *testing.T) {
	tmp := t.TempDir()
	writeInstalledPluginForRegistryTest(t, tmp, "demo")
	if err := SaveManifest(tmp, Manifest{Plugins: []Entry{{Name: "demo", Source: "local"}}}); err != nil {
		t.Fatal(err)
	}

	removed, err := Uninstall(tmp, "demo")
	if err != nil {
		t.Fatalf("Uninstall: %v", err)
	}
	if !removed {
		t.Fatal("expected plugin removed")
	}
	if _, err := os.Stat(filepath.Join(tmp, "addons", "demo")); !os.IsNotExist(err) {
		t.Fatalf("plugin files should be removed, stat err=%v", err)
	}
	plugins, err := List(tmp)
	if err != nil {
		t.Fatal(err)
	}
	if len(plugins) != 0 {
		t.Fatalf("registry entry should be removed, got %#v", plugins)
	}
}

func TestUninstallMissingDirectoryRemovesRegistryEntry(t *testing.T) {
	tmp := t.TempDir()
	if err := SaveManifest(tmp, Manifest{Plugins: []Entry{{Name: "demo", Source: "local"}}}); err != nil {
		t.Fatal(err)
	}

	removed, err := Uninstall(tmp, "demo")
	if err != nil {
		t.Fatalf("Uninstall: %v", err)
	}
	if !removed {
		t.Fatal("expected registry entry removed")
	}
	plugins, err := List(tmp)
	if err != nil {
		t.Fatal(err)
	}
	if len(plugins) != 0 {
		t.Fatalf("registry entry should be removed, got %#v", plugins)
	}
}

func TestUninstallOrphanDirectoryWithoutRegistryEntryRemovesFiles(t *testing.T) {
	tmp := t.TempDir()
	writeInstalledPluginForRegistryTest(t, tmp, "demo")

	removed, err := Uninstall(tmp, "demo")
	if err != nil {
		t.Fatalf("Uninstall: %v", err)
	}
	if !removed {
		t.Fatal("expected orphan plugin files to be removed")
	}
	if _, err := os.Stat(filepath.Join(tmp, "addons", "demo")); !os.IsNotExist(err) {
		t.Fatalf("orphan plugin files should be removed, stat err=%v", err)
	}
}

func TestUninstallMissingRegistryEntryAndDirectoryReturnsFalse(t *testing.T) {
	tmp := t.TempDir()
	removed, err := Uninstall(tmp, "demo")
	if err != nil {
		t.Fatalf("Uninstall: %v", err)
	}
	if removed {
		t.Fatal("expected no-op uninstall to return false")
	}
}

func TestUninstallRegistryFailureRestoresPluginFiles(t *testing.T) {
	tmp := t.TempDir()
	writeInstalledPluginForRegistryTest(t, tmp, "demo")
	if err := SaveManifest(tmp, Manifest{Plugins: []Entry{{Name: "demo", Source: "local"}}}); err != nil {
		t.Fatal(err)
	}
	origRename := registryRename
	defer func() { registryRename = origRename }()
	registryRename = func(oldpath, newpath string) error {
		if strings.Contains(filepath.Base(oldpath), ".tmp-") && filepath.Base(newpath) == "plugins.json" {
			return errors.New("registry locked")
		}
		return os.Rename(oldpath, newpath)
	}

	removed, err := Uninstall(tmp, "demo")
	if err == nil || !strings.Contains(err.Error(), "registry locked") {
		t.Fatalf("expected registry locked error, removed=%v err=%v", removed, err)
	}
	if removed {
		t.Fatal("failed uninstall must not report removed")
	}
	if _, err := os.Stat(filepath.Join(tmp, "addons", "demo", "plugin.json")); err != nil {
		t.Fatalf("plugin files should be restored: %v", err)
	}
	plugins, err := List(tmp)
	if err != nil {
		t.Fatal(err)
	}
	if len(plugins) != 1 || plugins[0].Name != "demo" {
		t.Fatalf("registry should retain plugin, got %#v", plugins)
	}
}

func TestUninstallCleanupFailureIsExplicitAndLeavesNoRunnablePlugin(t *testing.T) {
	tmp := t.TempDir()
	writeInstalledPluginForRegistryTest(t, tmp, "demo")
	if err := SaveManifest(tmp, Manifest{Plugins: []Entry{{Name: "demo", Source: "local"}}}); err != nil {
		t.Fatal(err)
	}
	origRemoveAll := registryRemoveAll
	defer func() { registryRemoveAll = origRemoveAll }()
	registryRemoveAll = func(path string) error {
		if strings.Contains(filepath.Clean(path), filepath.Clean(filepath.Join(".scenaria", "plugin-backups"))) {
			return errors.New("cleanup denied")
		}
		return os.RemoveAll(path)
	}

	removed, err := Uninstall(tmp, "demo")
	if err == nil || !strings.Contains(err.Error(), "cleanup failed") {
		t.Fatalf("expected cleanup failure, removed=%v err=%v", removed, err)
	}
	if !removed {
		t.Fatal("registry/files should be removed from runnable location before cleanup error")
	}
	if _, err := os.Stat(filepath.Join(tmp, "addons", "demo")); !os.IsNotExist(err) {
		t.Fatalf("runnable plugin dir should be gone, stat err=%v", err)
	}
}

func TestUninstallRejectsMaliciousPluginID(t *testing.T) {
	tmp := t.TempDir()
	if removed, err := Uninstall(tmp, "../outside"); err == nil || removed {
		t.Fatalf("expected invalid plugin id error without removal, removed=%v err=%v", removed, err)
	}
}

func writeInstalledPluginForRegistryTest(t *testing.T, root, name string) {
	t.Helper()
	dir := filepath.Join(root, "addons", name)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "plugin.json"), []byte(`{"id":"`+name+`"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "runner.exe"), []byte("binary"), 0o755); err != nil {
		t.Fatal(err)
	}
}
