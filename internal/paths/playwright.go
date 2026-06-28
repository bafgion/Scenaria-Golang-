package paths

import (
	"os"
	"path/filepath"
)

// ConfigurePlaywrightBrowsers sets PLAYWRIGHT_BROWSERS_PATH when bundled browsers exist.
func ConfigurePlaywrightBrowsers() string {
	if existing := os.Getenv("PLAYWRIGHT_BROWSERS_PATH"); existing != "" {
		return existing
	}
	if exe, err := os.Executable(); err == nil {
		bundled := filepath.Join(filepath.Dir(exe), "browsers")
		if info, err := os.Stat(bundled); err == nil && info.IsDir() {
			_ = os.Setenv("PLAYWRIGHT_BROWSERS_PATH", bundled)
			return bundled
		}
	}
	return ""
}

// BundledBrowsersDir returns the browsers directory next to the executable.
func BundledBrowsersDir() string {
	exe, err := os.Executable()
	if err != nil {
		return ""
	}
	return filepath.Join(filepath.Dir(exe), "browsers")
}
