package paths

import (
	"os"
	"path/filepath"
)

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

// ConfineToProjectRoot resolves userPath relative to projectRoot and rejects paths outside the project.
func ConfineToProjectRoot(projectRoot, userPath string) (string, error) {
	if userPath == "" {
		return "", nil
	}
	return PathGuard{Root: projectRoot}.ResolveExistingOrNew(userPath)
}
