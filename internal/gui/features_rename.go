package gui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func (s *Service) RenameFeature(path, newName string) (string, error) {
	srcAbs, err := s.ensureInsideProject(path)
	if err != nil {
		return "", err
	}
	newName = strings.TrimSpace(newName)
	if newName == "" {
		return "", fmt.Errorf("new file name is required")
	}
	if strings.ContainsAny(newName, `/\`) {
		return "", fmt.Errorf("invalid file name %q", newName)
	}
	if !strings.EqualFold(filepath.Ext(newName), ".feature") {
		newName += ".feature"
	}
	dir := filepath.Dir(srcAbs)
	target := filepath.Join(dir, newName)
	if strings.EqualFold(filepath.Clean(srcAbs), filepath.Clean(target)) {
		return srcAbs, nil
	}
	if _, err := os.Stat(target); err == nil {
		return "", fmt.Errorf("file already exists: %s", newName)
	} else if !os.IsNotExist(err) {
		return "", err
	}
	if err := os.Rename(srcAbs, target); err != nil {
		return "", fmt.Errorf("rename feature: %w", err)
	}
	return target, nil
}
