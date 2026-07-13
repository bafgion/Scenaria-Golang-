package gui

import (
	"fmt"
	"os"
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
	s.tempFeatureMu.Lock()
	defer s.tempFeatureMu.Unlock()
	for _, dir := range s.tempFeatureDirs {
		if confined, err := (paths.PathGuard{Root: dir}).ResolveExistingOrNew(path); err == nil {
			return confined, nil
		}
	}
	tmpRoot := os.TempDir()
	if confined, err := (paths.PathGuard{Root: tmpRoot}).ResolveExistingOrNew(path); err == nil {
		return confined, nil
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
