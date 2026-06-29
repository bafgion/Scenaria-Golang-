package gui

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func (s *Service) ensureInsideProject(target string) (string, error) {
	target = strings.TrimSpace(target)
	if target == "" {
		return "", fmt.Errorf("path is required")
	}
	project := s.ProjectPath()
	if project == "" {
		return "", fmt.Errorf("open a project folder first")
	}
	abs, err := filepath.Abs(target)
	if err != nil {
		return "", err
	}
	projectAbs, err := filepath.Abs(project)
	if err != nil {
		return "", err
	}
	rel, err := filepath.Rel(projectAbs, abs)
	if err != nil || strings.HasPrefix(rel, "..") {
		return "", fmt.Errorf("path is outside the project")
	}
	return abs, nil
}

// MoveFeature moves a .feature file into destDir (project folder or subfolder).
func (s *Service) MoveFeature(src, destDir string) (string, error) {
	srcAbs, err := s.ensureInsideProject(src)
	if err != nil {
		return "", err
	}
	if !strings.EqualFold(filepath.Ext(srcAbs), ".feature") {
		return "", fmt.Errorf("only .feature files can be moved")
	}
	destAbs, err := s.ensureInsideProject(destDir)
	if err != nil {
		return "", err
	}
	info, err := os.Stat(destAbs)
	if err != nil {
		return "", fmt.Errorf("destination folder: %w", err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("destination must be a folder")
	}
	target := filepath.Join(destAbs, filepath.Base(srcAbs))
	if strings.EqualFold(filepath.Clean(srcAbs), filepath.Clean(target)) {
		return srcAbs, nil
	}
	if _, err := os.Stat(target); err == nil {
		return "", fmt.Errorf("file already exists: %s", filepath.Base(target))
	} else if !os.IsNotExist(err) {
		return "", err
	}
	if err := os.Rename(srcAbs, target); err != nil {
		return "", fmt.Errorf("move feature: %w", err)
	}
	return target, nil
}

// ImportFeatures copies external .feature files into destDir inside the project.
func (s *Service) ImportFeatures(destDir string, paths []string) ([]string, error) {
	destAbs, err := s.ensureInsideProject(destDir)
	if err != nil {
		return nil, err
	}
	if info, err := os.Stat(destAbs); err != nil || !info.IsDir() {
		if err != nil {
			return nil, fmt.Errorf("destination folder: %w", err)
		}
		return nil, fmt.Errorf("destination must be a folder")
	}
	imported := make([]string, 0, len(paths))
	for _, src := range paths {
		src = strings.TrimSpace(src)
		if src == "" || !strings.EqualFold(filepath.Ext(src), ".feature") {
			continue
		}
		target := uniqueFeaturePath(filepath.Join(destAbs, filepath.Base(src)))
		if err := copyFile(src, target); err != nil {
			return imported, fmt.Errorf("import %s: %w", filepath.Base(src), err)
		}
		imported = append(imported, target)
	}
	if len(imported) == 0 {
		return nil, fmt.Errorf("no .feature files to import")
	}
	return imported, nil
}

func uniqueFeaturePath(path string) string {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return path
	}
	dir := filepath.Dir(path)
	base := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	ext := filepath.Ext(path)
	for i := 2; i < 100; i++ {
		candidate := filepath.Join(dir, fmt.Sprintf("%s-%d%s", base, i, ext))
		if _, err := os.Stat(candidate); os.IsNotExist(err) {
			return candidate
		}
	}
	return filepath.Join(dir, base+"-copy"+ext)
}

func copyFile(src, dest string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Close()
}
