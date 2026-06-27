package cli

import (
	"testing"

	"github.com/bafgion/scenaria-golang/internal/plugin"
)

func TestRunPluginsLifecycle(t *testing.T) {
	tmp := t.TempDir()

	if err := RunPlugins([]string{"install", "--project", tmp, "--name", "vanessa", "--source", "https://example.com/vanessa.zip"}); err != nil {
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

func TestRunPluginsErrors(t *testing.T) {
	if err := RunPlugins(nil); err == nil {
		t.Fatal("expected usage error")
	}
	if err := RunPlugins([]string{"install"}); err == nil {
		t.Fatal("expected missing flags error")
	}
	if err := RunPlugins([]string{"unknown"}); err == nil {
		t.Fatal("expected unknown subcommand error")
	}
}
