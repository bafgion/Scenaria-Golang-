package gui

import (
	"fmt"

	"github.com/bafgion/scenaria-golang/internal/plugin"
)

type PluginEntryDTO struct {
	Name        string `json:"name"`
	Source      string `json:"source"`
	InstalledAt string `json:"installedAt"`
}

func (s *Service) ListPlugins() ([]PluginEntryDTO, error) {
	path := s.ProjectPath()
	if path == "" {
		return nil, fmt.Errorf("open a project folder first")
	}
	entries, err := plugin.List(path)
	if err != nil {
		return nil, err
	}
	out := make([]PluginEntryDTO, 0, len(entries))
	for _, entry := range entries {
		out = append(out, PluginEntryDTO{
			Name:        entry.Name,
			Source:      entry.Source,
			InstalledAt: entry.InstalledAt,
		})
	}
	return out, nil
}

func (s *Service) InstallPlugin(name, source string) error {
	path := s.ProjectPath()
	if path == "" {
		return fmt.Errorf("open a project folder first")
	}
	if name == "" {
		return fmt.Errorf("plugin name is required")
	}
	return plugin.FetchAndInstall(path, name, source)
}

func (s *Service) UninstallPlugin(name string) error {
	path := s.ProjectPath()
	if path == "" {
		return fmt.Errorf("open a project folder first")
	}
	removed, err := plugin.Uninstall(path, name)
	if err != nil {
		return err
	}
	if !removed {
		return fmt.Errorf("plugin %q is not installed", name)
	}
	return nil
}
