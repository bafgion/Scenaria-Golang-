package gui

import (
	"fmt"

	"github.com/bafgion/scenaria-golang/internal/settings"
)

type ProjectConfigDTO struct {
	BaseURL            string `json:"baseUrl"`
	HTMLReportOpenMode string `json:"htmlReportOpenMode"`
	NavWaitUntil       string `json:"navWaitUntil"`
	DefaultRunner      string `json:"defaultRunner"`
	FeaturesRoot       string `json:"featuresRoot"`
}

func (s *Service) LoadProjectConfig() (ProjectConfigDTO, error) {
	root := s.ProjectPath()
	if root == "" {
		return ProjectConfigDTO{}, fmt.Errorf("open a project folder first")
	}
	cfg, err := settings.LoadProjectConfig(root)
	if err != nil {
		return ProjectConfigDTO{}, err
	}
	return projectConfigFromSettings(cfg), nil
}

func (s *Service) SaveProjectConfig(dto ProjectConfigDTO) error {
	root := s.ProjectPath()
	if root == "" {
		return fmt.Errorf("open a project folder first")
	}
	cfg, err := settings.LoadProjectConfig(root)
	if err != nil {
		return err
	}
	cfg.BaseURL = dto.BaseURL
	cfg.HTMLReportOpenMode = settings.NormalizeHTMLReportOpenMode(dto.HTMLReportOpenMode)
	if dto.NavWaitUntil != "" {
		cfg.NavWaitUntil = dto.NavWaitUntil
	}
	if dto.DefaultRunner != "" {
		cfg.DefaultRunner = dto.DefaultRunner
	}
	if dto.FeaturesRoot != "" {
		cfg.FeaturesRoot = dto.FeaturesRoot
	}
	return settings.SaveProjectConfig(root, cfg)
}

func projectConfigFromSettings(cfg settings.ProjectConfig) ProjectConfigDTO {
	return ProjectConfigDTO{
		BaseURL:            cfg.BaseURL,
		HTMLReportOpenMode: cfg.HTMLReportOpenMode,
		NavWaitUntil:       cfg.NavWaitUntil,
		DefaultRunner:      cfg.DefaultRunner,
		FeaturesRoot:       cfg.FeaturesRoot,
	}
}
