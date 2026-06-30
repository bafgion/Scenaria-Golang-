package brand

import (
	"os"
	"path/filepath"
	"runtime"
)

const assetsRel = "assets/branding"

// Asset file names (Python app/qt/branding.py).
const (
	IconMasterPNG = "icon-variant-b-monogram-su.png"
	IconSquarePNG = "app-icon-square.png"
	IconMarkPNG   = "app-icon-mark.png"
	IconICO       = "app.ico"
)

// RequiredAssets must exist in the branding folder for release builds.
var RequiredAssets = []string{
	IconMasterPNG,
	IconMarkPNG,
	IconSquarePNG,
	IconICO,
}

// Dir returns the branding assets directory (dev checkout or next to the executable).
func Dir() string {
	candidates := []string{
		filepath.Join(moduleRoot(), assetsRel),
	}
	if exe, err := os.Executable(); err == nil {
		candidates = append(candidates, filepath.Join(filepath.Dir(exe), assetsRel))
	}
	for _, path := range candidates {
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			return path
		}
	}
	return candidates[0]
}

// AssetPath joins Dir() with a file inside assets/branding.
func AssetPath(name string) string {
	return filepath.Join(Dir(), name)
}

func moduleRoot() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "."
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}
