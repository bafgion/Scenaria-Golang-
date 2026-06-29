package gui

import (
	"os"
	"path/filepath"
)

// ProjectArtifacts lists artifact directories for the open project.
type ProjectArtifacts struct {
	AllureDir  string `json:"allureDir"`
	TracesDir  string `json:"tracesDir"`
	VideosDir  string `json:"videosDir"`
	HTMLReport string `json:"htmlReport"`
	JUnitReport string `json:"junitReport"`
}

func (s *Service) ProjectArtifacts() ProjectArtifacts {
	root := s.ProjectPath()
	if root == "" {
		return ProjectArtifacts{}
	}
	scenaria := filepath.Join(root, ".scenaria")
	out := ProjectArtifacts{}
	if s.ArtifactExists(filepath.Join(scenaria, "allure-results")) {
		out.AllureDir = filepath.Join(scenaria, "allure-results")
	}
	if s.ArtifactExists(filepath.Join(scenaria, "traces")) {
		out.TracesDir = filepath.Join(scenaria, "traces")
	}
	if s.ArtifactExists(filepath.Join(scenaria, "videos")) {
		out.VideosDir = filepath.Join(scenaria, "videos")
	}
	if s.ArtifactExists(filepath.Join(scenaria, "report.html")) {
		out.HTMLReport = filepath.Join(scenaria, "report.html")
	}
	if s.ArtifactExists(filepath.Join(scenaria, "junit.xml")) {
		out.JUnitReport = filepath.Join(scenaria, "junit.xml")
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
