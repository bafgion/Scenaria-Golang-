package gui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bafgion/scenaria-golang/internal/scenario"
)

type ProjectReplaceRequest struct {
	Find          string `json:"find"`
	Replace       string `json:"replace"`
	CaseSensitive bool   `json:"caseSensitive"`
}

type ProjectReplaceResult struct {
	FilesChanged int      `json:"filesChanged"`
	Replacements int      `json:"replacements"`
	Files        []string `json:"files"`
}

func (s *Service) ReplaceInProject(req ProjectReplaceRequest) (ProjectReplaceResult, error) {
	path := s.ProjectPath()
	if path == "" {
		return ProjectReplaceResult{}, fmt.Errorf("open a project folder first")
	}
	find := req.Find
	if strings.TrimSpace(find) == "" {
		return ProjectReplaceResult{}, fmt.Errorf("find text is required")
	}
	store := scenario.NewFeatureStore()
	files, err := store.Discover(path)
	if err != nil {
		return ProjectReplaceResult{}, err
	}
	result := ProjectReplaceResult{Files: []string{}}
	for _, file := range files {
		payload, err := os.ReadFile(file)
		if err != nil {
			continue
		}
		original := string(payload)
		replaced := ReplaceInText(original, find, req.Replace, req.CaseSensitive)
		if replaced.Count == 0 {
			continue
		}
		if err := os.WriteFile(file, []byte(replaced.Text), 0o644); err != nil {
			return result, fmt.Errorf("write %s: %w", file, err)
		}
		result.FilesChanged++
		result.Replacements += replaced.Count
		result.Files = append(result.Files, file)
	}
	return result, nil
}

func (s *Service) DeleteFeature(path string) error {
	path = strings.TrimSpace(path)
	if path == "" {
		return fmt.Errorf("feature path is required")
	}
	project := s.ProjectPath()
	if project == "" {
		return fmt.Errorf("open a project folder first")
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	projectAbs, err := filepath.Abs(project)
	if err != nil {
		return err
	}
	rel, err := filepath.Rel(projectAbs, abs)
	if err != nil || strings.HasPrefix(rel, "..") {
		return fmt.Errorf("feature is outside the project")
	}
	if err := os.Remove(abs); err != nil {
		return fmt.Errorf("delete feature: %w", err)
	}
	return nil
}

func (s *Service) DuplicateFeature(path, newName string) (string, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return "", fmt.Errorf("feature path is required")
	}
	payload, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read feature: %w", err)
	}
	dir := filepath.Dir(path)
	ext := filepath.Ext(path)
	target := ""
	if name := strings.TrimSpace(newName); name != "" {
		name = strings.TrimSuffix(name, ext)
		name = strings.TrimSuffix(name, ".feature")
		if name == "" {
			return "", fmt.Errorf("new feature name is required")
		}
		target = filepath.Join(dir, name+ext)
		if _, err := os.Stat(target); err == nil {
			return "", fmt.Errorf("file already exists: %s", filepath.Base(target))
		} else if !os.IsNotExist(err) {
			return "", err
		}
	} else {
		base := strings.TrimSuffix(filepath.Base(path), ext)
		target = filepath.Join(dir, base+"-copy"+ext)
		for i := 2; i < 100; i++ {
			if _, err := os.Stat(target); os.IsNotExist(err) {
				break
			}
			target = filepath.Join(dir, fmt.Sprintf("%s-copy-%d%s", base, i, ext))
		}
	}
	if err := os.WriteFile(target, payload, 0o644); err != nil {
		return "", fmt.Errorf("write duplicate: %w", err)
	}
	return target, nil
}
