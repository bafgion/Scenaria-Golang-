package paths

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/bafgion/scenaria-golang/internal/brand"
)

func AppDataDir() string {
	if runtime.GOOS == "windows" {
		base := os.Getenv("APPDATA")
		if base != "" {
			return filepath.Join(base, brand.AppDataDir)
		}
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return filepath.Join(home, ".scenaria")
}

func InferProjectRoot(paths []string) string {
	if len(paths) == 0 {
		return ""
	}
	files := make([]string, 0)
	dirs := make([]string, 0)
	for _, raw := range paths {
		info, err := os.Stat(raw)
		if err != nil {
			continue
		}
		if info.IsDir() {
			dirs = append(dirs, raw)
		} else {
			files = append(files, raw)
		}
	}
	if len(files) > 0 {
		parents := map[string]struct{}{}
		for _, file := range files {
			parents[filepath.Dir(file)] = struct{}{}
		}
		if len(parents) == 1 {
			for parent := range parents {
				return parent
			}
		}
	}
	if len(dirs) == 1 {
		return dirs[0]
	}
	if len(files) == 1 {
		return filepath.Dir(files[0])
	}
	return ""
}

func ScenariaProjectDir(projectRoot string) string {
	return filepath.Join(projectRoot, ".scenaria")
}
