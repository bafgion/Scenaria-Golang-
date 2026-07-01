package paths

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/bafgion/scenaria-golang/internal/brand"
)

// AppDataDir is the single Scenaria application data root (settings, artifact mirrors, caches).
func AppDataDir() string {
	if runtime.GOOS == "windows" {
		if base := os.Getenv("APPDATA"); base != "" {
			return filepath.Join(base, brand.AppDataDir)
		}
	}
	if runtime.GOOS == "darwin" {
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, "Library", "Application Support", brand.AppDataDir)
		}
	}
	if xdg := os.Getenv("XDG_DATA_HOME"); xdg != "" {
		return filepath.Join(xdg, brand.AppDataDir)
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return filepath.Join(home, ".local", "share", brand.AppDataDir)
}

// LegacySettingsPaths lists deprecated settings.json locations for migration reads.
func LegacySettingsPaths() []string {
	out := make([]string, 0, 2)
	if config, err := os.UserConfigDir(); err == nil {
		out = append(out, filepath.Join(config, brand.AppDataDir, "settings.json"))
	}
	if home, err := os.UserHomeDir(); err == nil {
		out = append(out, filepath.Join(home, ".scenaria", "settings.json"))
	}
	return out
}

// LegacyProjectMirrorDir is the pre-unification artifact mirror path (~/.scenaria/projects/<slug>).
func LegacyProjectMirrorDir(projectSlug string) string {
	home, err := os.UserHomeDir()
	if err != nil || projectSlug == "" {
		return ""
	}
	return filepath.Join(home, ".scenaria", "projects", projectSlug)
}

// FileExists reports whether path is present on disk.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
