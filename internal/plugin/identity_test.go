package plugin

import (
	"archive/zip"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidatePluginID(t *testing.T) {
	valid := []string{
		"demo",
		"Demo_1",
		"com.example-plugin",
		"plugin.2026_07",
	}
	for _, id := range valid {
		if err := ValidatePluginID(id); err != nil {
			t.Fatalf("ValidatePluginID(%q) unexpected error: %v", id, err)
		}
	}

	invalid := []string{
		"",
		" demo",
		"demo ",
		".",
		"..",
		"../outside",
		`..\outside`,
		"plugin/name",
		`plugin\name`,
		`C:\outside`,
		`C:outside`,
		`\\server\share`,
		"demo.",
		"demo name",
		"demo\nname",
		"demo:name",
	}
	for _, id := range invalid {
		if err := ValidatePluginID(id); err == nil {
			t.Fatalf("ValidatePluginID(%q) expected error", id)
		}
	}
}

func TestAddonPathIsConfinedToAddonsRoot(t *testing.T) {
	project := t.TempDir()
	got, err := addonPath(project, "demo.plugin")
	if err != nil {
		t.Fatalf("addonPath valid id: %v", err)
	}
	want := filepath.Join(project, "addons", "demo.plugin")
	if got != want {
		t.Fatalf("addonPath = %q, want %q", got, want)
	}
	if _, err := addonPath(project, "../outside"); err == nil {
		t.Fatal("expected traversal plugin id to fail")
	}
}

func TestAddonPathRejectsAddonsSymlinkEscape(t *testing.T) {
	project := t.TempDir()
	outside := t.TempDir()
	createPluginDirSymlinkOrSkip(t, outside, filepath.Join(project, "addons"))

	if _, err := addonPath(project, "demo.plugin"); err == nil {
		t.Fatal("expected addons symlink outside project to fail")
	}
}

func TestAddonPathRejectsPluginSymlinkEscape(t *testing.T) {
	project := t.TempDir()
	outside := t.TempDir()
	addons := filepath.Join(project, "addons")
	if err := os.MkdirAll(addons, 0o755); err != nil {
		t.Fatal(err)
	}
	createPluginDirSymlinkOrSkip(t, outside, filepath.Join(addons, "demo.plugin"))

	if _, err := addonPath(project, "demo.plugin"); err == nil {
		t.Fatal("expected plugin symlink outside addons root to fail")
	}
}

func TestFetchAndInstallRejectsMalformedIDBeforeFilesystemMutation(t *testing.T) {
	tmp := t.TempDir()
	project := filepath.Join(tmp, "project")
	if err := os.MkdirAll(project, 0o755); err != nil {
		t.Fatal(err)
	}
	outside := filepath.Join(tmp, "outside")
	if err := os.MkdirAll(outside, 0o755); err != nil {
		t.Fatal(err)
	}
	sentinel := filepath.Join(outside, "sentinel.txt")
	if err := os.WriteFile(sentinel, []byte("keep"), 0o644); err != nil {
		t.Fatal(err)
	}

	zipPath := filepath.Join(tmp, "addon.zip")
	writePluginZip(t, zipPath)

	err := FetchAndInstall(project, "../outside", zipPath)
	if err == nil {
		t.Fatal("expected malformed plugin id error")
	}
	raw, readErr := os.ReadFile(sentinel)
	if readErr != nil {
		t.Fatalf("sentinel should remain readable: %v", readErr)
	}
	if string(raw) != "keep" {
		t.Fatalf("sentinel changed: %q", string(raw))
	}
	if _, statErr := os.Stat(filepath.Join(project, "addons")); !os.IsNotExist(statErr) {
		t.Fatalf("addons directory should not be created for malformed id, stat err=%v", statErr)
	}
}

func TestRegistryRejectsMalformedPluginIDs(t *testing.T) {
	project := t.TempDir()
	if err := Install(project, "demo", "source.zip"); err != nil {
		t.Fatalf("valid install: %v", err)
	}

	if err := Install(project, `..\outside`, "source.zip"); err == nil {
		t.Fatal("expected invalid install id to fail")
	}
	if removed, err := Uninstall(project, `..\outside`); err == nil || removed {
		t.Fatalf("expected invalid uninstall id to fail without removal, removed=%v err=%v", removed, err)
	}
	if err := SaveManifest(project, Manifest{Plugins: []Entry{{Name: "../outside"}}}); err == nil {
		t.Fatal("expected invalid manifest entry to fail")
	}

	plugins, err := List(project)
	if err != nil {
		t.Fatal(err)
	}
	if len(plugins) != 1 || plugins[0].Name != "demo" {
		t.Fatalf("valid manifest entry should remain, got %#v", plugins)
	}
}

func TestLoadDescriptorRejectsMalformedID(t *testing.T) {
	project := t.TempDir()
	outside := filepath.Join(project, "..", "outside")
	if err := os.MkdirAll(outside, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(outside, "plugin.json"), []byte(`{"id":"outside"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := LoadDescriptor(project, "../outside")
	if err == nil {
		t.Fatal("expected malformed descriptor id to fail")
	}
	if !strings.Contains(err.Error(), "invalid plugin id") {
		t.Fatalf("expected validation error, got %v", err)
	}
}

func writePluginZip(t *testing.T, zipPath string) {
	t.Helper()
	file, err := os.Create(zipPath)
	if err != nil {
		t.Fatal(err)
	}
	w := zip.NewWriter(file)
	entry, err := w.Create("plugin.json")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := entry.Write([]byte(`{"name":"demo"}`)); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
}

func createPluginDirSymlinkOrSkip(t *testing.T, target, link string) {
	t.Helper()
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("symlink not available: %v", err)
	}
}
