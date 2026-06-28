package cli

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"

	"github.com/bafgion/scenaria-golang/internal/plugin"
)

func TestRunPluginsLifecycle(t *testing.T) {
	tmp := t.TempDir()
	zipPath := filepath.Join(tmp, "vanessa.zip")
	file, err := os.Create(zipPath)
	if err != nil {
		t.Fatal(err)
	}
	w := zip.NewWriter(file)
	entry, err := w.Create("plugin.json")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = entry.Write([]byte(`{"name":"vanessa"}`))
	_ = w.Close()
	_ = file.Close()

	if err := RunPlugins([]string{"install", "--project", tmp, "--name", "vanessa", "--source", zipPath}); err != nil {
		t.Fatalf("install command failed: %v", err)
	}
	plugins, err := plugin.List(tmp)
	if err != nil {
		t.Fatalf("plugin.List failed: %v", err)
	}
	if len(plugins) != 1 || plugins[0].Name != "vanessa" {
		t.Fatalf("unexpected plugins state: %+v", plugins)
	}

	if err := RunPlugins([]string{"uninstall", "--project", tmp, "--name", "vanessa"}); err != nil {
		t.Fatalf("uninstall command failed: %v", err)
	}
	plugins, err = plugin.List(tmp)
	if err != nil {
		t.Fatalf("plugin.List failed: %v", err)
	}
	if len(plugins) != 0 {
		t.Fatalf("expected no plugins after uninstall: %+v", plugins)
	}
}
