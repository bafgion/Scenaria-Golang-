package paths

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
	userPath = strings.TrimSpace(userPath)
	if userPath == "" {
		return "", nil
	}
	absRoot, err := filepath.Abs(projectRoot)
	if err != nil {
		return "", fmt.Errorf("resolve project root: %w", err)
	}
	var absPath string
	if filepath.IsAbs(userPath) {
		absPath = filepath.Clean(userPath)
	} else {
		absPath = filepath.Clean(filepath.Join(absRoot, userPath))
	}
	rel, err := filepath.Rel(absRoot, absPath)
	if err != nil || strings.HasPrefix(rel, "..") {
		return "", fmt.Errorf("path must be inside project: %s", userPath)
	}
	return absPath, nil
}
