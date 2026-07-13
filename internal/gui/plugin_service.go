package gui

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/bafgion/scenaria-golang/internal/plugin"
)

type pluginCLIRunner interface {
	VA(ctx context.Context, args []string) (string, error)
	Run(ctx context.Context, args []string) (string, error)
}

// PluginService owns project-scoped plugin list/install/run operations.
type PluginService struct {
	projectPath func() string
	cliRunner   func() pluginCLIRunner
	mu          sync.Mutex
	nextRunID   atomic.Uint64
	active      map[uint64]pluginInvocation
}

type pluginInvocation struct {
	name   string
	cancel context.CancelFunc
}

func NewPluginService(projectPath func() string, cliRunner func() pluginCLIRunner) *PluginService {
	return &PluginService{
		projectPath: projectPath,
		cliRunner:   cliRunner,
		active:      make(map[uint64]pluginInvocation),
	}
}

func (s *PluginService) requireProject() (string, error) {
	if s == nil || s.projectPath == nil {
		return "", fmt.Errorf("open a project folder first")
	}
	path := s.projectPath()
	if path == "" {
		return "", fmt.Errorf("open a project folder first")
	}
	return path, nil
}

func (s *PluginService) List() ([]PluginEntryDTO, error) {
	path, err := s.requireProject()
	if err != nil {
		return nil, err
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

func (s *PluginService) Install(name, source string) error {
	path, err := s.requireProject()
	if err != nil {
		return err
	}
	if name == "" {
		return fmt.Errorf("plugin name is required")
	}
	return plugin.FetchAndInstall(path, name, source)
}

func (s *PluginService) Uninstall(name string) error {
	path, err := s.requireProject()
	if err != nil {
		return err
	}
	s.CancelPlugin(name)
	removed, err := plugin.Uninstall(path, name)
	if err != nil {
		return err
	}
	if !removed {
		return fmt.Errorf("plugin %q is not installed", name)
	}
	return nil
}

func (s *PluginService) Run(req PluginRunRequest) RunResult {
	path, err := s.requireProject()
	if err != nil {
		return RunResult{Error: err.Error()}
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
	cli := s.cliRunner()
	if cli == nil {
		return RunResult{Error: "cli runner is not configured"}
	}
	switch target.Runner {
	case "va":
		ctx, finish := s.startInvocation(name)
		defer finish()
		out, runErr := cli.VA(ctx, appendVanessaArgs(append([]string(nil), target.Args...), req))
		if runErr != nil {
			return RunResult{Output: out, Error: runErr.Error()}
		}
		return RunResult{Output: out}
	case "run":
		args := appendRunPluginArgs(append([]string(nil), target.Args...), req)
		ctx, finish := s.startInvocation(name)
		defer finish()
		out, runErr := cli.Run(ctx, args)
		if runErr != nil {
			return RunResult{Output: out, Error: runErr.Error()}
		}
		return RunResult{Output: out}
	default:
		return RunResult{Error: fmt.Sprintf("unsupported plugin runner %q", target.Runner)}
	}
}

func (s *PluginService) startInvocation(name string) (context.Context, func()) {
	ctx, cancel := context.WithCancel(context.Background())
	if s == nil {
		return ctx, cancel
	}
	id := s.nextRunID.Add(1)
	s.mu.Lock()
	if s.active == nil {
		s.active = make(map[uint64]pluginInvocation)
	}
	s.active[id] = pluginInvocation{name: strings.TrimSpace(name), cancel: cancel}
	s.mu.Unlock()
	return ctx, func() {
		s.mu.Lock()
		delete(s.active, id)
		s.mu.Unlock()
		cancel()
	}
}

func (s *PluginService) CancelActive() {
	if s == nil {
		return
	}
	s.mu.Lock()
	active := make([]context.CancelFunc, 0, len(s.active))
	for id, invocation := range s.active {
		active = append(active, invocation.cancel)
		delete(s.active, id)
	}
	s.mu.Unlock()
	for _, cancel := range active {
		cancel()
	}
}

func (s *PluginService) CancelPlugin(name string) {
	if s == nil {
		return
	}
	name = strings.TrimSpace(name)
	s.mu.Lock()
	active := make([]context.CancelFunc, 0)
	for id, invocation := range s.active {
		if strings.EqualFold(invocation.name, name) {
			active = append(active, invocation.cancel)
			delete(s.active, id)
		}
	}
	s.mu.Unlock()
	for _, cancel := range active {
		cancel()
	}
}

func (s *PluginService) RunVanessa(dryRun bool) RunResult {
	return s.Run(PluginRunRequest{Name: "vanessa", DryRun: dryRun})
}

func (s *PluginService) runVanessa(req PluginRunRequest) RunResult {
	path, err := s.requireProject()
	if err != nil {
		return RunResult{Error: err.Error()}
	}
	cli := s.cliRunner()
	if cli == nil {
		return RunResult{Error: "cli runner is not configured"}
	}
	args := []string{"run", "--project", path}
	if req.DryRun {
		args = append(args, "--dry-run")
	}
	args = appendVanessaArgs(args, req)
	ctx, finish := s.startInvocation(vanessaInvocationName(req))
	defer finish()
	out, err := cli.VA(ctx, args)
	if err != nil {
		return RunResult{Output: out, Error: err.Error()}
	}
	return RunResult{Output: out}
}

func vanessaInvocationName(req PluginRunRequest) string {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return "vanessa"
	}
	return name
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
	if scenario := strings.TrimSpace(req.Scenario); scenario != "" {
		args = append(args, "--scenario", scenario)
	}
	return args
}
