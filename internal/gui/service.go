package gui

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/bafgion/scenaria-golang/internal/cli"
	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/recorder"
	"github.com/bafgion/scenaria-golang/internal/runstatus"
	"github.com/bafgion/scenaria-golang/internal/scenario"
	"github.com/bafgion/scenaria-golang/internal/selector"
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
	Path         string              `json:"path"`
	Features     []string            `json:"features"`
	Tags         []string            `json:"tags"`
	FeatureTags  map[string][]string `json:"featureTags"`
}

type RunRequest struct {
	Tag        string            `json:"tag"`
	Scenario   string            `json:"scenario"`
	TestClient string            `json:"testClient"`
	Vars       map[string]string `json:"vars"`
	DryRun     bool              `json:"dryRun"`
	Headed     bool              `json:"headed"`
	Engine     string            `json:"engine"`
	InstallPW  bool              `json:"installPlaywright"`
	AllureDir  string            `json:"allureDir"`
	TraceDir   string            `json:"traceDir"`
	VideoDir   string            `json:"videoDir"`
	HTMLPath      string            `json:"htmlPath"`
	JUnitPath     string            `json:"junitPath"`
	SummaryJSON   string            `json:"summaryJson"`
	Targets       []string          `json:"targets"`
	Browser       string            `json:"browser"`
	Workers       int               `json:"workers"`
	SlowMo        int               `json:"slowMo"`
	BaseURL       string            `json:"baseUrl"`
	StartStep     int               `json:"startStep"`
	EndStep       int               `json:"endStep"`
}

type ValidateRequest struct {
	Browser     string   `json:"browser"`
	SkipBrowser bool     `json:"skipBrowser"`
	Targets     []string `json:"targets"`
}

type PluginRunRequest struct {
	Name              string   `json:"name"`
	DryRun            bool     `json:"dryRun"`
	Tag               string   `json:"tag"`
	ExcludeTags       []string `json:"excludeTags"`
	Scenario          string   `json:"scenario"`
	RerunFailedRunDir string   `json:"rerunFailedRunDir"`
	InstallEPF        bool     `json:"installEpf"`
	EPFURL            string   `json:"epfUrl"`
	EPFDest           string   `json:"epfDest"`
	PlatformExe       string   `json:"platformExe"`
	EPFPath           string   `json:"epfPath"`
	IBConnection      string   `json:"ibConnection"`
	ReportAllure      bool     `json:"reportAllure"`
	VaDir             string   `json:"vaDir"`
	VaFiles           string   `json:"vaFiles"`
}

type RunResult struct {
	Output string `json:"output"`
	Error  string `json:"error"`
}

type StepCatalogEntry struct {
	Label       string   `json:"label"`
	Action      string   `json:"action"`
	Category    string   `json:"category"`
	Description string   `json:"description"`
	Template    string   `json:"template"`
	Example     string   `json:"example"`
	Parameters  []string `json:"parameters"`
	Help        string   `json:"help"`
}

type StepCompletionSnippet struct {
	Label       string `json:"label"`
	Insert      string `json:"insert"`
	Description string `json:"description"`
}

type StepCompletionsDTO struct {
	Start int                     `json:"start"`
	End   int                     `json:"end"`
	Items []StepCompletionSnippet `json:"items"`
}

type AppSettingsDTO struct {
	Browser           string `json:"browser"`
	Headless          bool   `json:"headless"`
	ParallelWorkers   int    `json:"parallelWorkers"`
	SlowMo            int    `json:"slowMo"`
	MaxLoopIterations int    `json:"maxLoopIterations"`
	FilterRecording   bool   `json:"filterRecording"`
	NavOnlyRecording  bool   `json:"navOnlyRecording"`
	HoverRecord       bool   `json:"hoverRecord"`
	ToolbarCompact    bool     `json:"toolbarCompact"`
	StepsPanelVisible bool     `json:"stepsPanelVisible"`
	StepsPanelHeight  int      `json:"stepsPanelHeight"`
	SidebarWidth      int      `json:"sidebarWidth"`
	RecentProjects    []string `json:"recentProjects"`
	RecentFeatures    []string `json:"recentFeatures"`
	SessionProject    string   `json:"sessionProject"`
	OpenTabs          []string `json:"openTabs"`
	ActiveTab         string   `json:"activeTab"`
	ScrollBeforeClick bool     `json:"scrollBeforeClick"`
	HoverRecordMinMs        int      `json:"hoverRecordMinMs"`
	SelectorClickStrategies []string `json:"selectorClickStrategies"`
	SelectorInputStrategies []string `json:"selectorInputStrategies"`
	CheckUpdatesOnStartup bool `json:"checkUpdatesOnStartup"`
	Editor            settings.EditorSettings `json:"editor"`
}

type RunResultEntry struct {
	Path    string `json:"path"`
	Success bool   `json:"success"`
	Message string `json:"message"`
	Runner  string `json:"runner"`
	At      string `json:"at"`
}

func (s *Service) ListRunResults(limit int) ([]RunResultEntry, error) {
	path := s.ProjectPath()
	if path == "" {
		return []RunResultEntry{}, nil
	}
	store, err := runstatus.Open(path)
	if err != nil {
		return nil, err
	}
	entries, err := store.List(limit)
	if err != nil {
		return nil, err
	}
	out := make([]RunResultEntry, 0, len(entries))
	for _, e := range entries {
		out = append(out, RunResultEntry{
			Path:    e.Path,
			Success: e.Success,
			Message: e.Message,
			Runner:  e.Runner,
			At:      e.At,
		})
	}
	return out, nil
}

func (s *Service) BundledExamplesPath() string {
	candidates := []string{"examples"}
	if exe, err := os.Executable(); err == nil {
		candidates = append(candidates, filepath.Join(filepath.Dir(exe), "examples"))
	}
	for _, c := range candidates {
		abs, err := filepath.Abs(c)
		if err != nil {
			continue
		}
		if st, err := os.Stat(abs); err == nil && st.IsDir() {
			return abs
		}
	}
	return ""
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
	featureTags := collectFeatureTags(store, files)
	return ProjectInfo{Path: path, Features: files, Tags: tags, FeatureTags: featureTags}, nil
}

func collectFeatureTags(store *scenario.FeatureStore, files []string) map[string][]string {
	out := make(map[string][]string, len(files))
	for _, file := range files {
		feature, err := store.Load(file)
		if err != nil {
			continue
		}
		tags := gherkin.CollectFeatureTags(feature)
		if len(tags) > 0 {
			out[file] = tags
		}
	}
	return out
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

func (s *Service) WriteTempFeature(content string) (string, error) {
	dir, err := os.MkdirTemp("", "scenaria-run-")
	if err != nil {
		return "", fmt.Errorf("create temp dir: %w", err)
	}
	path := filepath.Join(dir, "scenario.feature")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return "", fmt.Errorf("write temp feature: %w", err)
	}
	return path, nil
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
	args := []string{}
	if len(req.Targets) > 0 {
		args = append(args, req.Targets...)
	} else if path != "" {
		args = append(args, path)
	} else {
		return RunResult{Error: "нет файлов для запуска — откройте сценарий или проект"}
	}
	if req.DryRun {
		args = append(args, "--dry-run")
	}
	if req.Tag != "" {
		args = append(args, "--tag", req.Tag)
	}
	if req.Scenario != "" {
		args = append(args, "--scenario", req.Scenario)
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
	if req.AllureDir != "" {
		args = append(args, "--allure", req.AllureDir)
	}
	if req.TraceDir != "" {
		args = append(args, "--trace", req.TraceDir)
	}
	if req.VideoDir != "" {
		args = append(args, "--video", req.VideoDir)
	}
	if req.HTMLPath != "" {
		args = append(args, "--html", req.HTMLPath)
	}
	if req.JUnitPath != "" {
		args = append(args, "--junit", req.JUnitPath)
	}
	if req.SummaryJSON != "" {
		args = append(args, "--summary-json", req.SummaryJSON)
	}
	if req.Browser != "" && !req.DryRun {
		args = append(args, "--browser", req.Browser)
	}
	if req.Workers > 0 {
		args = append(args, "--workers", fmt.Sprintf("%d", req.Workers))
	}
	if req.SlowMo > 0 && !req.DryRun {
		args = append(args, "--slow-mo", fmt.Sprintf("%d", req.SlowMo))
	}
	if req.BaseURL != "" && !req.DryRun {
		args = append(args, "--base-url", req.BaseURL)
	}
	if req.StartStep >= 0 {
		args = append(args, "--start-step", fmt.Sprintf("%d", req.StartStep))
	}
	if req.EndStep >= 0 {
		args = append(args, "--end-step", fmt.Sprintf("%d", req.EndStep))
	}
	out, err := captureCLI(func() error { return cli.RunRun(args) })
	if err != nil {
		return RunResult{Output: out, Error: err.Error()}
	}
	return RunResult{Output: out}
}

func (s *Service) Validate(req ValidateRequest) RunResult {
	path := s.ProjectPath()
	if path == "" {
		return RunResult{Error: "open a project folder first"}
	}
	args := []string{}
	if len(req.Targets) > 0 {
		args = append(args, req.Targets...)
	} else {
		args = append(args, path)
	}
	if req.SkipBrowser {
		args = append(args, "--no-browser")
	} else if req.Browser != "" {
		args = append(args, "--browser", req.Browser)
	}
	out, err := captureCLI(func() error { return cli.RunValidate(args) })
	if err != nil {
		return RunResult{Output: out, Error: err.Error()}
	}
	return RunResult{Output: out}
}

func (s *Service) CheckUpdate() RunResult {
	out, err := captureCLI(func() error { return cli.RunUpdate([]string{"--check"}) })
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

func (s *Service) ReadTestClientJSON(name string) (string, error) {
	path := s.ProjectPath()
	if path == "" {
		return "", fmt.Errorf("open a project folder first")
	}
	return settings.ReadTestClientJSON(path, name)
}

func (s *Service) SaveTestClientJSON(name, content string) error {
	path := s.ProjectPath()
	if path == "" {
		return fmt.Errorf("open a project folder first")
	}
	return settings.SaveTestClientFromJSON(path, name, content)
}

func (s *Service) DeleteTestClient(name string) error {
	path := s.ProjectPath()
	if path == "" {
		return fmt.Errorf("open a project folder first")
	}
	return settings.DeleteTestClient(path, name)
}

func (s *Service) SearchSteps(query string) []StepCatalogEntry {
	entries := stepcatalog.Search(query)
	out := make([]StepCatalogEntry, 0, len(entries))
	for _, entry := range entries {
		out = append(out, StepCatalogEntry{
			Label:       entry.Label,
			Action:      entry.Action,
			Category:    entry.Category,
			Description: entry.Description,
			Template:    entry.Template,
			Example:     entry.Example,
			Parameters:  entry.Parameters,
			Help:        entry.Help,
		})
	}
	return out
}

func (s *Service) CompletionsForLine(line string, column int) StepCompletionsDTO {
	result := stepcatalog.CompletionsForLine(line, column)
	out := StepCompletionsDTO{
		Start: result.Start,
		End:   result.End,
		Items: make([]StepCompletionSnippet, 0, len(result.Items)),
	}
	for _, item := range result.Items {
		out.Items = append(out.Items, StepCompletionSnippet{
			Label:       item.Label,
			Insert:      item.Insert,
			Description: item.Description,
		})
	}
	return out
}

func (s *Service) LoadSettings() (AppSettingsDTO, error) {
	cfg, err := settings.LoadDefaultAppSettings()
	if err != nil || cfg == nil {
		return defaultAppSettingsDTO(), nil
	}
	return appSettingsFromCfg(cfg), nil
}

func defaultAppSettingsDTO() AppSettingsDTO {
	return AppSettingsDTO{
		Browser:           "chromium",
		ParallelWorkers:   1,
		MaxLoopIterations: 100,
		StepsPanelVisible: true,
		StepsPanelHeight:  160,
		CheckUpdatesOnStartup: true,
		Editor:            settings.DefaultEditorSettings(),
	}
}

func appSettingsFromCfg(cfg *settings.AppSettings) AppSettingsDTO {
	height := cfg.StepsPanelHeight
	if height < 80 {
		height = 160
	}
	return AppSettingsDTO{
		Browser:           cfg.Browser,
		Headless:          cfg.Headless,
		ParallelWorkers:   maxInt(1, cfg.ParallelWorkers),
		SlowMo:            maxInt(0, cfg.SlowMo),
		MaxLoopIterations: maxInt(1, cfg.MaxLoopIterations),
		FilterRecording:   cfg.RecordingFilterMode,
		NavOnlyRecording:    cfg.NavOnlyRecording,
		HoverRecord:         cfg.RecordingHoverMode,
		ToolbarCompact:      cfg.ToolbarCompact,
		StepsPanelVisible:   cfg.StepsPanelVisible,
		StepsPanelHeight:    height,
		SidebarWidth:        clampSidebarWidth(cfg.SidebarWidth),
		RecentProjects:      trimRecents(cfg.RecentProjects),
		RecentFeatures:      trimRecents(cfg.RecentFeatures),
		SessionProject:      strings.TrimSpace(cfg.SessionProject),
		OpenTabs:            trimRecents(cfg.OpenTabs),
		ActiveTab:           strings.TrimSpace(cfg.ActiveTab),
		ScrollBeforeClick:   cfg.ScrollBeforeClick,
		HoverRecordMinMs:        maxInt(0, cfg.HoverRecordMinMs),
		SelectorClickStrategies:   selector.NormalizeClickStrategies(cfg.SelectorClickStrategies),
		SelectorInputStrategies: selector.NormalizeInputStrategies(cfg.SelectorInputStrategies),
		CheckUpdatesOnStartup: settings.CheckUpdatesOnStartupEnabled(cfg),
		Editor:              settings.NormalizeEditorSettings(cfg.Editor),
	}
}

func (s *Service) SaveSettings(dto AppSettingsDTO) error {
	path := settings.DefaultAppSettingsPath()
	if path == "" {
		return fmt.Errorf("settings path unavailable")
	}
	height := dto.StepsPanelHeight
	if height < 80 {
		height = 160
	}
	existing, _ := settings.LoadDefaultAppSettings()
	cfg := &settings.AppSettings{
		Browser:             dto.Browser,
		Headless:            dto.Headless,
		ParallelWorkers:     maxInt(1, dto.ParallelWorkers),
		SlowMo:              maxInt(0, dto.SlowMo),
		MaxLoopIterations:   maxInt(1, dto.MaxLoopIterations),
		RecordingFilterMode: dto.FilterRecording,
		NavOnlyRecording:      dto.NavOnlyRecording,
		RecordingHoverMode:    dto.HoverRecord,
		ToolbarCompact:        dto.ToolbarCompact,
		StepsPanelVisible:     dto.StepsPanelVisible,
		StepsPanelHeight:      height,
		SidebarWidth:          clampSidebarWidth(dto.SidebarWidth),
		RecentProjects:        trimRecents(dto.RecentProjects),
		RecentFeatures:        trimRecents(dto.RecentFeatures),
		SessionProject:          strings.TrimSpace(dto.SessionProject),
		OpenTabs:                trimRecents(dto.OpenTabs),
		ActiveTab:               strings.TrimSpace(dto.ActiveTab),
		ScrollBeforeClick:       dto.ScrollBeforeClick,
		HoverRecordMinMs:        normalizeHoverRecordMinMs(dto.HoverRecordMinMs),
		SelectorClickStrategies: selector.NormalizeClickStrategies(dto.SelectorClickStrategies),
		SelectorInputStrategies: selector.NormalizeInputStrategies(dto.SelectorInputStrategies),
	}
	checkUpdates := dto.CheckUpdatesOnStartup
	cfg.CheckUpdatesOnStartup = &checkUpdates
	cfg.Editor = settings.NormalizeEditorSettings(dto.Editor)
	if existing != nil {
		cfg.HTTPAuth = existing.HTTPAuth
		if len(cfg.RecentProjects) == 0 {
			cfg.RecentProjects = existing.RecentProjects
		}
		if len(cfg.RecentFeatures) == 0 {
			cfg.RecentFeatures = existing.RecentFeatures
		}
	}
	sidebarW := dto.SidebarWidth
	if sidebarW < 120 && existing != nil && existing.SidebarWidth >= 120 {
		sidebarW = existing.SidebarWidth
	}
	cfg.SidebarWidth = clampSidebarWidth(sidebarW)
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

func normalizeHoverRecordMinMs(ms int) int {
	if ms <= 0 {
		return 600
	}
	return ms
}

func runExport(args []string) error     { return cli.RunExport(args) }
func runImportJSON(args []string) error { return cli.RunImportJSON(args) }
func runRecord(args []string) error     { return cli.RunRecord(args) }
func runVA(args []string) error         { return cli.RunVA(args) }
