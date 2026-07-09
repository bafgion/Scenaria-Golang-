package gui

import (
	"fmt"
	"os"
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
	result := ProjectReplaceResult{Files: []string{}}
	if err := s.withProjectFSWriteLock(func() error {
		store := scenario.NewFeatureStore()
		files, err := store.Discover(path)
		if err != nil {
			return err
		}
		for _, file := range files {
			payload, err := os.ReadFile(file)
			if err != nil {
				continue
			}
			original := string(payload)
			replaced := ReplaceAllInText(original, find, req.Replace, req.CaseSensitive, true)
			if replaced.Count == 0 {
				continue
			}
			if err := os.WriteFile(file, []byte(replaced.Text), 0o644); err != nil {
				return fmt.Errorf("write %s: %w", file, err)
			}
			result.FilesChanged++
			result.Replacements += replaced.Count
			result.Files = append(result.Files, file)
		}
		return nil
	}); err != nil {
		return result, err
	}
	return result, nil
}

func (s *Service) DeleteFeature(path string) error {
	return s.fileOperator().DeleteFeature(path)
}

func (s *Service) DuplicateFeature(path, newName string) (string, error) {
	return s.fileOperator().DuplicateFeature(path, newName)
}
