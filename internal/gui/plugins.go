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

func (s *Service) RunPlugin(req PluginRunRequest) RunResult {
	path := s.ProjectPath()
	if path == "" {
		return RunResult{Error: "open a project folder first"}
	}
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return RunResult{Error: "plugin name is required"}
	}
	desc, err := plugin.LoadDescriptor(path, name)
	if err != nil {
		if strings.EqualFold(name, "vanessa") {
			return s.runVanessa(req)
		}
		return RunResult{Error: err.Error()}
	}
	if plugin.IsVanessa(desc) {
		return s.runVanessa(req)
	}
	target, err := plugin.ResolveRun(path, desc, req.DryRun)
	if err != nil {
		return RunResult{Error: err.Error()}
	}
	switch target.Runner {
	case "va":
		out, runErr := captureCLI(func() error { return cliRunVA(appendVanessaArgs(target.Args, req)) })
		if runErr != nil {
			return RunResult{Output: out, Error: runErr.Error()}
		}
		return RunResult{Output: out}
	case "run":
		args := appendRunPluginArgs(append([]string(nil), target.Args...), req)
		out, runErr := captureCLI(func() error { return cli.RunRun(args) })
		if runErr != nil {
			return RunResult{Output: out, Error: runErr.Error()}
		}
		return RunResult{Output: out}
	default:
		return RunResult{Error: fmt.Sprintf("unsupported plugin runner %q", target.Runner)}
	}
}

func (s *Service) runVanessa(req PluginRunRequest) RunResult {
	path := s.ProjectPath()
	args := []string{"run", "--project", path}
	if req.DryRun {
		args = append(args, "--dry-run")
	}
	args = appendVanessaArgs(args, req)
	out, err := captureCLI(func() error { return cliRunVA(args) })
	if err != nil {
		return RunResult{Output: out, Error: err.Error()}
	}
	return RunResult{Output: out}
}

func appendVanessaArgs(args []string, req PluginRunRequest) []string {
	if tag := strings.TrimSpace(req.Tag); tag != "" {
		args = append(args, "--tag", tag)
	}
	for _, tag := range req.ExcludeTags {
		tag = strings.TrimSpace(tag)
		if tag != "" {
			args = append(args, "--exclude-tag", tag)
		}
	}
	if scenario := strings.TrimSpace(req.Scenario); scenario != "" {
		args = append(args, "--scenario", scenario)
	}
	if dir := strings.TrimSpace(req.RerunFailedRunDir); dir != "" {
		args = append(args, "--rerun-failed", dir)
	}
	if req.InstallEPF {
		args = append(args, "--epf-install")
	}
	if url := strings.TrimSpace(req.EPFURL); url != "" {
		args = append(args, "--epf-url", url)
	}
	if dest := strings.TrimSpace(req.EPFDest); dest != "" {
		args = append(args, "--epf-dest", dest)
	}
	if exe := strings.TrimSpace(req.PlatformExe); exe != "" {
		args = append(args, "--platform-exe", exe)
	}
	if epf := strings.TrimSpace(req.EPFPath); epf != "" {
		args = append(args, "--epf", epf)
	}
	if ib := strings.TrimSpace(req.IBConnection); ib != "" {
		args = append(args, "--ib", ib)
	}
	if req.ReportAllure {
		args = append(args, "--allure")
	}
	if dir := strings.TrimSpace(req.VaDir); dir != "" {
		args = append(args, "--dir", dir)
	}
	if files := strings.TrimSpace(req.VaFiles); files != "" {
		args = append(args, "--files", files)
	}
	return args
}

func appendRunPluginArgs(args []string, req PluginRunRequest) []string {
	if tag := strings.TrimSpace(req.Tag); tag != "" {
		args = append(args, "--tag", tag)
	}
	return args
}

// RunVanessa keeps backward compatibility for simple dry/run calls.
func (s *Service) RunVanessa(dryRun bool) RunResult {
	return s.RunPlugin(PluginRunRequest{Name: "vanessa", DryRun: dryRun})
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
