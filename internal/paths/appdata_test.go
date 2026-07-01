package paths

import (
	"path/filepath"
	"runtime"
	"testing"
)

func TestAppDataDirUsesXDGOnLinux(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("linux-specific")
	}
	xdg := t.TempDir()
	t.Setenv("XDG_DATA_HOME", xdg)
	got := AppDataDir()
	want := filepath.Join(xdg, "Scenaria")
	if got != want {
		t.Fatalf("AppDataDir = %q, want %q", got, want)
	}
}

func TestLegacySettingsPathsIncludeConfigDir(t *testing.T) {
	paths := LegacySettingsPaths()
	if len(paths) < 2 {
		t.Fatalf("expected legacy settings paths, got %v", paths)
	}
}
