package gui

import (
	"os"
	"path/filepath"

	"github.com/bafgion/scenaria-golang/internal/paths"
)

type projectPathProvider func() string

// ReportService owns report/artifact path discovery for a project.
type ReportService struct {
	projectPath projectPathProvider
}

func NewReportService(projectPath projectPathProvider) *ReportService {
	return &ReportService{projectPath: projectPath}
}

func (s *ReportService) scenariaDirs() []string {
	if s == nil || s.projectPath == nil {
		return nil
	}
	root := s.projectPath()
	if root == "" {
		return nil
	}
	dirs := []string{filepath.Join(root, ".scenaria")}
	if writable, err := paths.WritableScenariaDir(root); err == nil && writable != dirs[0] {
		dirs = append([]string{writable}, dirs...)
	}
	return dirs
}

func (s *ReportService) ScenariaArtifactPath(sub string) string {
	if s == nil || s.projectPath == nil {
		return ""
	}
	root := s.projectPath()
	if root == "" {
		return ""
	}
	path, err := paths.ScenariaArtifactPath(root, sub)
	if err != nil {
		return filepath.Join(root, ".scenaria", sub)
	}
	return path
}

func (s *ReportService) ArtifactExists(path string) bool {
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

func (s *ReportService) ProjectArtifacts() ProjectArtifacts {
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
