package gui

import (
	"bytes"
	"errors"
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

var (
	projectReplaceReadFile  = os.ReadFile
	projectReplaceStat      = os.Stat
	projectReplaceWriteFile = writeFileAtomic
)

type plannedProjectReplacement struct {
	Path  string
	Old   []byte
	New   []byte
	Count int
	Perm  os.FileMode
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
		plan, plannedResult, err := buildProjectReplacementPlan(files, find, req.Replace, req.CaseSensitive)
		if err != nil {
			return err
		}
		if err := commitProjectReplacementPlan(plan); err != nil {
			return err
		}
		result = plannedResult
		return nil
	}); err != nil {
		return ProjectReplaceResult{Files: []string{}}, err
	}
	return result, nil
}

func buildProjectReplacementPlan(files []string, find, replace string, caseSensitive bool) ([]plannedProjectReplacement, ProjectReplaceResult, error) {
	plan := make([]plannedProjectReplacement, 0)
	result := ProjectReplaceResult{Files: []string{}}
	for _, file := range files {
		info, err := projectReplaceStat(file)
		if err != nil {
			return nil, result, fmt.Errorf("stat %s: %w", filepath.Base(file), err)
		}
		payload, err := projectReplaceReadFile(file)
		if err != nil {
			return nil, result, fmt.Errorf("read %s: %w", filepath.Base(file), err)
		}
		replaced := ReplaceAllInText(string(payload), find, replace, caseSensitive, true)
		if replaced.Count == 0 {
			continue
		}
		oldContent := append([]byte(nil), payload...)
		newContent := []byte(replaced.Text)
		perm := info.Mode().Perm()
		if perm == 0 {
			perm = 0o644
		}
		plan = append(plan, plannedProjectReplacement{
			Path:  file,
			Old:   oldContent,
			New:   newContent,
			Count: replaced.Count,
			Perm:  perm,
		})
		result.FilesChanged++
		result.Replacements += replaced.Count
		result.Files = append(result.Files, file)
	}
	return plan, result, nil
}

func commitProjectReplacementPlan(plan []plannedProjectReplacement) error {
	committed := make([]plannedProjectReplacement, 0, len(plan))
	for _, item := range plan {
		if err := validateProjectReplacementUnchanged(item); err != nil {
			if rollbackErr := rollbackProjectReplacements(committed); rollbackErr != nil {
				return fmt.Errorf("%w (rollback failed: %v)", err, rollbackErr)
			}
			return err
		}
		if err := projectReplaceWriteFile(item.Path, item.New, item.Perm); err != nil {
			if rollbackErr := rollbackProjectReplacements(committed); rollbackErr != nil {
				return fmt.Errorf("write %s: %w (rollback failed: %v)", filepath.Base(item.Path), err, rollbackErr)
			}
			return fmt.Errorf("write %s: %w", filepath.Base(item.Path), err)
		}
		committed = append(committed, item)
	}
	return nil
}

func validateProjectReplacementUnchanged(item plannedProjectReplacement) error {
	current, err := projectReplaceReadFile(item.Path)
	if err != nil {
		return fmt.Errorf("verify %s before replace: %w", filepath.Base(item.Path), err)
	}
	if !bytes.Equal(current, item.Old) {
		return fmt.Errorf("replace aborted: file changed externally: %s", filepath.Base(item.Path))
	}
	return nil
}

func rollbackProjectReplacements(committed []plannedProjectReplacement) error {
	var failures []string
	for i := len(committed) - 1; i >= 0; i-- {
		item := committed[i]
		if err := projectReplaceWriteFile(item.Path, item.Old, item.Perm); err != nil {
			failures = append(failures, fmt.Sprintf("%s: %v", filepath.Base(item.Path), err))
		}
	}
	if len(failures) > 0 {
		return errors.New(strings.Join(failures, "; "))
	}
	return nil
}

func (s *Service) DeleteFeature(path string) error {
	return s.fileOperator().DeleteFeature(path)
}

func (s *Service) DuplicateFeature(path, newName string) (string, error) {
	return s.fileOperator().DuplicateFeature(path, newName)
}
