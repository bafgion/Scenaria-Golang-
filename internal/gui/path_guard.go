package gui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bafgion/scenaria-golang/internal/paths"
)

func (s *Service) confineFeaturePath(path string) (string, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return "", fmt.Errorf("feature path is required")
	}
	if root := s.ProjectPath(); root != "" {
		return paths.ConfineToProjectRoot(root, path)
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	s.tempFeatureMu.Lock()
	defer s.tempFeatureMu.Unlock()
	for _, dir := range s.tempFeatureDirs {
		dirAbs, err := filepath.Abs(dir)
		if err != nil {
			continue
		}
		rel, err := filepath.Rel(dirAbs, abs)
		if err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
			return abs, nil
		}
	}
	tmpRoot := os.TempDir()
	if rel, err := filepath.Rel(tmpRoot, abs); err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
		return abs, nil
	}
	return "", fmt.Errorf("feature path is outside the allowed temp run directories")
}

func (s *Service) confineArtifactPath(path string) (string, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return "", fmt.Errorf("artifact path is required")
	}
	root := s.ProjectPath()
	if root == "" {
		return "", fmt.Errorf("open a project folder first")
	}
	return paths.ConfineToProjectRoot(root, path)
}
