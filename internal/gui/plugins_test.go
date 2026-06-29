package gui

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"
)

func TestListPluginsRequiresProject(t *testing.T) {
	svc := NewService()
	if _, err := svc.ListPlugins(); err == nil {
		t.Fatal("expected error without project")
	}
}

func TestPluginLifecycleViaService(t *testing.T) {
	root := t.TempDir()
	svc := NewService()
	svc.mu.Lock()
	svc.projectPath = root
	svc.mu.Unlock()

	zipPath := filepath.Join(t.TempDir(), "demo.zip")
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

	if err := svc.InstallPlugin("demo", zipPath); err != nil {
		t.Fatalf("install: %v", err)
	}
	plugins, err := svc.ListPlugins()
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(plugins) != 1 || plugins[0].Name != "demo" {
		t.Fatalf("got %+v", plugins)
	}
	if err := svc.UninstallPlugin("demo"); err != nil {
		t.Fatalf("uninstall: %v", err)
	}
}
