package plugin

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"
)

func TestFetchAndInstallLocalZip(t *testing.T) {
	tmp := t.TempDir()
	zipPath := filepath.Join(tmp, "addon.zip")
	file, err := os.Create(zipPath)
	if err != nil {
		t.Fatal(err)
	}
	w := zip.NewWriter(file)
	entry, err := w.Create("plugin.json")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = entry.Write([]byte(`{"name":"demo"}`))
	_ = w.Close()
	_ = file.Close()

	project := filepath.Join(tmp, "project")
	if err := FetchAndInstall(project, "demo", zipPath); err != nil {
		t.Fatalf("FetchAndInstall: %v", err)
	}
	if _, err := os.Stat(filepath.Join(project, "addons", "demo", "plugin.json")); err != nil {
		t.Fatalf("expected extracted plugin.json: %v", err)
	}
	plugins, err := List(project)
	if err != nil {
		t.Fatal(err)
	}
	if len(plugins) != 1 || plugins[0].Name != "demo" {
		t.Fatalf("unexpected manifest: %#v", plugins)
	}
}
