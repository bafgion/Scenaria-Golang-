package gui

import (
	"os"
	"path/filepath"

	"github.com/bafgion/scenaria-golang/internal/paths"
)

// ProjectArtifacts lists artifact directories for the open project.
type ProjectArtifacts struct {
	AllureDir  string `json:"allureDir"`
	TracesDir  string `json:"tracesDir"`
	VideosDir  string `json:"videosDir"`
	HTMLReport  string `json:"htmlReport"`
	JUnitReport string `json:"junitReport"`
	SummaryJSON string `json:"summaryJson"`
}

func (s *Service) scenariaDirs() []string {
	root := s.ProjectPath()
	if root == "" {
		return nil
	}
	dirs := []string{filepath.Join(root, ".scenaria")}
	if writable, err := paths.WritableScenariaDir(root); err == nil && writable != dirs[0] {
		dirs = append([]string{writable}, dirs...)
	}
	return dirs
}

func (s *Service) ScenariaArtifactPath(sub string) string {
	root := s.ProjectPath()
	if root == "" {
		return ""
	}
	path, err := paths.ScenariaArtifactPath(root, sub)
	if err != nil {
		return filepath.Join(root, ".scenaria", sub)
	}
	return path
}

func (s *Service) ProjectArtifacts() ProjectArtifacts {
	dirs := s.scenariaDirs()
	if len(dirs) == 0 {
		return ProjectArtifacts{}
	}
	out := ProjectArtifacts{}
	for _, scenaria := range dirs {
		if out.AllureDir == "" && s.ArtifactExists(filepath.Join(scenaria, "allure-results")) {
			out.AllureDir = filepath.Join(scenaria, "allure-results")
		}
		if out.TracesDir == "" && s.ArtifactExists(filepath.Join(scenaria, "traces")) {
			out.TracesDir = filepath.Join(scenaria, "traces")
		}
		if out.VideosDir == "" && s.ArtifactExists(filepath.Join(scenaria, "videos")) {
			out.VideosDir = filepath.Join(scenaria, "videos")
		}
		if out.HTMLReport == "" && s.ArtifactExists(filepath.Join(scenaria, "report.html")) {
			out.HTMLReport = filepath.Join(scenaria, "report.html")
		}
		if out.JUnitReport == "" && s.ArtifactExists(filepath.Join(scenaria, "junit.xml")) {
			out.JUnitReport = filepath.Join(scenaria, "junit.xml")
		}
		if out.SummaryJSON == "" && s.ArtifactExists(filepath.Join(scenaria, "summary.json")) {
			out.SummaryJSON = filepath.Join(scenaria, "summary.json")
		}
	}
	return out
}

func (s *Service) ParseEditorSteps(text string) []EditorStepRow {
	return ParseEditorSteps(text)
}

func (s *Service) ArtifactExists(path string) bool {
	if path == "" {
		return false
	}
	st, err := os.Stat(path)
	if err != nil {
		return false
	}
	if st.IsDir() {
		entries, err := os.ReadDir(path)
		return err == nil && len(entries) > 0
	}
	return true
}
