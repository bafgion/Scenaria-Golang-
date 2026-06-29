package gui

import (
	"fmt"
	"strings"

	"github.com/bafgion/scenaria-golang/internal/cli"
	"github.com/bafgion/scenaria-golang/internal/plugin"
)

type PluginEntryDTO struct {
	Name        string   `json:"name"`
	Source      string   `json:"source"`
	InstalledAt string   `json:"installedAt"`
	ID          string   `json:"id"`
	Description string   `json:"description"`
	Commands    []string `json:"commands"`
	Runnable    bool     `json:"runnable"`
	Vanessa     bool     `json:"vanessa"`
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
		out = append(out, pluginEntryDTO(path, entry))
	}
	return out, nil
}

func pluginEntryDTO(projectRoot string, entry plugin.Entry) PluginEntryDTO {
	dto := PluginEntryDTO{
		Name:        entry.Name,
		Source:      entry.Source,
		InstalledAt: entry.InstalledAt,
		ID:          entry.Name,
	}
	desc, err := plugin.LoadDescriptor(projectRoot, entry.Name)
	if err != nil {
		if strings.EqualFold(entry.Name, "vanessa") {
			dto.Vanessa = true
			dto.Runnable = true
			dto.Commands = []string{"va run"}
		}
		return dto
	}
	dto.ID = desc.ID
	dto.Description = desc.Description
	dto.Commands = append([]string(nil), desc.Commands...)
	dto.Vanessa = plugin.IsVanessa(desc)
	dto.Runnable = plugin.IsRunnable(desc)
	return dto
}

func (s *Service) RunPlugin(name string, dryRun bool) RunResult {
	path := s.ProjectPath()
	if path == "" {
		return RunResult{Error: "open a project folder first"}
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return RunResult{Error: "plugin name is required"}
	}
	if strings.EqualFold(name, "vanessa") {
		return s.RunVanessa(dryRun)
	}
	desc, err := plugin.LoadDescriptor(path, name)
	if err != nil {
		if strings.EqualFold(name, "vanessa") {
			return s.RunVanessa(dryRun)
		}
		return RunResult{Error: err.Error()}
	}
	target, err := plugin.ResolveRun(path, desc, dryRun)
	if err != nil {
		return RunResult{Error: err.Error()}
	}
	switch target.Runner {
	case "va":
		out, runErr := captureCLI(func() error { return cliRunVA(target.Args) })
		if runErr != nil {
			return RunResult{Output: out, Error: runErr.Error()}
		}
		return RunResult{Output: out}
	case "run":
		out, runErr := captureCLI(func() error { return cli.RunRun(target.Args) })
		if runErr != nil {
			return RunResult{Output: out, Error: runErr.Error()}
		}
		return RunResult{Output: out}
	default:
		return RunResult{Error: fmt.Sprintf("unsupported plugin runner %q", target.Runner)}
	}
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
