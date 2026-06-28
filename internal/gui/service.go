package gui

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/bafgion/scenaria-golang/internal/cli"
	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/recorder"
	"github.com/bafgion/scenaria-golang/internal/scenario"
	"github.com/bafgion/scenaria-golang/internal/settings"
	"github.com/bafgion/scenaria-golang/internal/stepcatalog"
	"github.com/bafgion/scenaria-golang/internal/version"
)

// Service exposes project and runner operations without UI framework dependencies.
type Service struct {
	mu           sync.RWMutex
	projectPath  string
	liveSession  *recorder.LiveSession
	recordCancel context.CancelFunc
}

func NewService() *Service {
	return &Service{}
}

type ProjectInfo struct {
	Path     string   `json:"path"`
	Features []string `json:"features"`
	Tags     []string `json:"tags"`
}

type RunRequest struct {
	Tag        string            `json:"tag"`
	TestClient string            `json:"testClient"`
	Vars       map[string]string `json:"vars"`
	DryRun     bool              `json:"dryRun"`
	Headed     bool              `json:"headed"`
	Engine     string            `json:"engine"`
	InstallPW  bool              `json:"installPlaywright"`
}

type RunResult struct {
	Output string `json:"output"`
	Error  string `json:"error"`
}

type StepCatalogEntry struct {
	Category string `json:"category"`
	Template string `json:"template"`
	Help     string `json:"help"`
}

type AppSettingsDTO struct {
	Browser           string `json:"browser"`
	Headless          bool   `json:"headless"`
	ParallelWorkers   int    `json:"parallelWorkers"`
	MaxLoopIterations int    `json:"maxLoopIterations"`
}

func (s *Service) Version() string {
	return version.String()
}

func (s *Service) ProjectPath() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.projectPath
}

func (s *Service) OpenProject(path string) (ProjectInfo, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return ProjectInfo{}, fmt.Errorf("project path is required")
	}
	info, err := os.Stat(path)
	if err != nil {
		return ProjectInfo{}, fmt.Errorf("open project: %w", err)
	}
	if !info.IsDir() {
		return ProjectInfo{}, fmt.Errorf("project path must be a directory")
	}
	s.mu.Lock()
	s.projectPath = path
	s.mu.Unlock()
	return s.projectInfo()
}

func (s *Service) projectInfo() (ProjectInfo, error) {
	s.mu.RLock()
	path := s.projectPath
	s.mu.RUnlock()
	if path == "" {
		return ProjectInfo{}, fmt.Errorf("no project opened")
	}
	store := scenario.NewFeatureStore()
	files, err := store.Discover(path)
	if err != nil {
		return ProjectInfo{}, err
	}
	tags := collectProjectTags(store, files)
	return ProjectInfo{Path: path, Features: files, Tags: tags}, nil
}

func collectProjectTags(store *scenario.FeatureStore, files []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, 16)
	for _, file := range files {
		feature, err := store.Load(file)
		if err != nil {
			continue
		}
		for _, tag := range gherkin.CollectFeatureTags(feature) {
			if _, ok := seen[tag]; ok {
				continue
			}
			seen[tag] = struct{}{}
			out = append(out, tag)
		}
	}
	return out
}

func (s *Service) ReadFeature(path string) (string, error) {
	payload, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read feature: %w", err)
	}
	return string(payload), nil
}

func (s *Service) SaveFeature(path, content string) error {
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return fmt.Errorf("save feature: %w", err)
	}
	return nil
}

func (s *Service) InitProject() (string, error) {
	path := s.ProjectPath()
	if path == "" {
		return "", fmt.Errorf("open a project folder first")
	}
	return captureCLI(func() error { return cli.RunInit([]string{path}) })
}

func (s *Service) Run(req RunRequest) RunResult {
	path := s.ProjectPath()
	if path == "" {
		return RunResult{Error: "open a project folder first"}
	}
	args := []string{path}
	if req.DryRun {
		args = append(args, "--dry-run")
	}
	if req.Tag != "" {
		args = append(args, "--tag", req.Tag)
	}
	if req.TestClient != "" {
		args = append(args, "--test-client", req.TestClient)
	}
	for key, value := range req.Vars {
		args = append(args, "--var", key+"="+value)
	}
	if req.Engine != "" {
		args = append(args, "--engine", req.Engine)
	}
	if req.Headed {
		args = append(args, "--headed")
	}
	if req.InstallPW {
		args = append(args, "--install-playwright")
	}
	out, err := captureCLI(func() error { return cli.RunRun(args) })
	if err != nil {
		return RunResult{Output: out, Error: err.Error()}
	}
	return RunResult{Output: out}
}

func (s *Service) Validate(browser string, skipBrowser bool) RunResult {
	path := s.ProjectPath()
	if path == "" {
		return RunResult{Error: "open a project folder first"}
	}
	args := []string{path}
	if skipBrowser {
		args = append(args, "--no-browser")
	} else if browser != "" {
		args = append(args, "--browser", browser)
	}
	out, err := captureCLI(func() error { return cli.RunValidate(args) })
	if err != nil {
		return RunResult{Output: out, Error: err.Error()}
	}
	return RunResult{Output: out}
}

func (s *Service) ListTestClients() ([]string, error) {
	path := s.ProjectPath()
	if path == "" {
		return nil, fmt.Errorf("open a project folder first")
	}
	return settings.ListTestClientNames(path)
}

func (s *Service) TestClientDetails(name string) (string, error) {
	path := s.ProjectPath()
	if path == "" {
		return "", fmt.Errorf("open a project folder first")
	}
	client, err := settings.LoadTestClientByName(path, name)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("name=%s base_url=%s cookies=%d local_storage=%d",
		client.Name, client.BaseURL, len(client.Cookies), len(client.LocalStorage)), nil
}

func (s *Service) SearchSteps(query string) []StepCatalogEntry {
	entries := stepcatalog.Search(query)
	out := make([]StepCatalogEntry, 0, len(entries))
	for _, entry := range entries {
		out = append(out, StepCatalogEntry{
			Category: entry.Category,
			Template: entry.Template,
			Help:     entry.Help,
		})
	}
	return out
}

func (s *Service) LoadSettings() (AppSettingsDTO, error) {
	cfg, err := settings.LoadDefaultAppSettings()
	if err != nil || cfg == nil {
		return AppSettingsDTO{Browser: "chromium", ParallelWorkers: 1, MaxLoopIterations: 100}, nil
	}
	return AppSettingsDTO{
		Browser:           cfg.Browser,
		Headless:          cfg.Headless,
		ParallelWorkers:   maxInt(1, cfg.ParallelWorkers),
		MaxLoopIterations: maxInt(1, cfg.MaxLoopIterations),
	}, nil
}

func (s *Service) SaveSettings(dto AppSettingsDTO) error {
	path := settings.DefaultAppSettingsPath()
	if path == "" {
		return fmt.Errorf("settings path unavailable")
	}
	cfg := &settings.AppSettings{
		Browser:           dto.Browser,
		Headless:          dto.Headless,
		ParallelWorkers:   maxInt(1, dto.ParallelWorkers),
		MaxLoopIterations: maxInt(1, dto.MaxLoopIterations),
	}
	if err := settings.SaveAppSettings(path, cfg); err != nil {
		return fmt.Errorf("save settings: %w", err)
	}
	return nil
}

func captureCLI(fn func() error) (string, error) {
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		return "", err
	}
	os.Stdout = w
	runErr := fn()
	_ = w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	_ = r.Close()
	return buf.String(), runErr
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func runExport(args []string) error     { return cli.RunExport(args) }
func runImportJSON(args []string) error { return cli.RunImportJSON(args) }
func runVA(args []string) error         { return cli.RunVA(args) }
