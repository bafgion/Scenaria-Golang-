package gui

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bafgion/scenaria-golang/internal/logx"
	"github.com/bafgion/scenaria-golang/internal/player"
	"github.com/bafgion/scenaria-golang/internal/report"
	"github.com/bafgion/scenaria-golang/internal/selector"
	"github.com/bafgion/scenaria-golang/internal/settings"
	"github.com/bafgion/scenaria-golang/internal/version"
)

// Service exposes project and runner operations without UI framework dependencies.
type Service struct {
	mu                    sync.RWMutex
	projectPath           string
	projectSession        *ProjectSession
	projectService        *ProjectService
	editorAnalysisService *EditorAnalysisService
	reportService         *ReportService
	runService            *RunService
	fileOps               *FileOperationService
	recorderService       *RecorderService
	settingsService       *SettingsService
	testClientService     *TestClientService
	pluginService         *PluginService
	catalogService        *CatalogService
	cliOps                *CLIOps
	projectVersion        uint64
	validateCancel        context.CancelFunc
	validateGen           uint64
	tempFeatureMu         sync.Mutex
	tempFeatureDirs       []string
	reportBridgeMu        sync.Mutex
	reportBridge          *reportBridge
	allureServe           allureServeState
	activePlaywright      sync.WaitGroup
	activeBackground      sync.WaitGroup
	settingsStore         *settings.Store
	projectFSMu           sync.RWMutex
	activePlaywrightCount atomic.Int32
	liveDirtyTabs         atomic.Bool
}

func NewService() *Service {
	svc := &Service{settingsStore: settings.DefaultStore()}
	svc.projectService = NewProjectService(svc.withProjectFSReadLock)
	svc.editorAnalysisService = NewEditorAnalysisService()
	svc.reportService = svc.buildReportService()
	svc.runService = NewRunService(svc.ProjectPath)
	svc.fileOps = NewFileOperationService(svc.confineFeaturePath, svc.ProjectPath, svc.withProjectFSReadLock, svc.withProjectFSWriteLock)
	svc.recorderService = NewRecorderService()
	svc.settingsService = NewSettingsService(svc.settingsStore)
	svc.testClientService = NewTestClientService(svc.ProjectPath)
	svc.pluginService = NewPluginService(svc.ProjectPath, func() pluginCLIRunner { return svc.cliRunner() })
	svc.catalogService = NewCatalogService()
	svc.cliOps = NewCLIOps()
	return svc
}

// ProjectSession identifies active project ownership for long-running operations.
type ProjectSession struct {
	ID      string          `json:"id"`
	Root    string          `json:"root"`
	Version uint64          `json:"version"`
	Context context.Context `json:"-"`
	cancel  context.CancelFunc
}

// RunSession identifies one in-flight run and its immutable snapshot.
type RunSession struct {
	RunID           string          `json:"runId"`
	ProjectVersion  uint64          `json:"projectVersion"`
	RequestSnapshot RunRequest      `json:"requestSnapshot"`
	Context         context.Context `json:"-"`
	TempResources   []string        `json:"tempResources,omitempty"`
}

type ProjectInfo struct {
	Path        string              `json:"path"`
	Features    []string            `json:"features"`
	Tags        []string            `json:"tags"`
	FeatureTags map[string][]string `json:"featureTags"`
	Version     uint64              `json:"version"`
}

type RunRequest struct {
	Tag              string            `json:"tag"`
	Scenario         string            `json:"scenario"`
	TestClient       string            `json:"testClient"`
	Vars             map[string]string `json:"vars"`
	DryRun           bool              `json:"dryRun"`
	Headed           bool              `json:"headed"`
	Engine           string            `json:"engine"`
	InstallPW        bool              `json:"installPlaywright"`
	AllureDir        string            `json:"allureDir"`
	TraceDir         string            `json:"traceDir"`
	VideoDir         string            `json:"videoDir"`
	HTMLPath         string            `json:"htmlPath"`
	JUnitPath        string            `json:"junitPath"`
	SummaryJSON      string            `json:"summaryJson"`
	Targets          []string          `json:"targets"`
	Browser          string            `json:"browser"`
	Workers          int               `json:"workers"`
	SlowMo           int               `json:"slowMo"`
	BaseURL          string            `json:"baseUrl"`
	StartStep        int               `json:"startStep"`
	EndStep          int               `json:"endStep"`
	ContinueOnFail   bool              `json:"continueOnFail"`
	HTMLLightMode    bool              `json:"htmlLightMode"`
	ReuseLiveBrowser bool              `json:"reuseLiveBrowser"`
	ReportLocale     string            `json:"reportLocale"`
}

type ValidateRequest struct {
	Browser     string   `json:"browser"`
	SkipBrowser bool     `json:"skipBrowser"`
	Targets     []string `json:"targets"`
	Mode        string   `json:"mode"`
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
	Output      string           `json:"output"`
	Error       string           `json:"error"`
	Entries     []RunResultEntry `json:"entries,omitempty"`
	ReportPath  string           `json:"reportPath,omitempty"`
	HTMLPath    string           `json:"htmlPath,omitempty"`
	JUnitPath   string           `json:"junitPath,omitempty"`
	SummaryJSON string           `json:"summaryJson,omitempty"`
	AllureDir   string           `json:"allureDir,omitempty"`
	TraceDir    string           `json:"traceDir,omitempty"`
	VideoDir    string           `json:"videoDir,omitempty"`
}

type AsyncRunResultDTO struct {
	JobID  string    `json:"jobId"`
	Result RunResult `json:"result"`
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
	Browser                 string                  `json:"browser"`
	Headless                bool                    `json:"headless"`
	ParallelWorkers         int                     `json:"parallelWorkers"`
	SlowMo                  int                     `json:"slowMo"`
	MaxLoopIterations       int                     `json:"maxLoopIterations"`
	NavWaitUntil            string                  `json:"navWaitUntil"`
	FilterRecording         bool                    `json:"filterRecording"`
	NavOnlyRecording        bool                    `json:"navOnlyRecording"`
	DisableRecordURLWait    bool                    `json:"disableRecordUrlWait"`
	HoverRecord             bool                    `json:"hoverRecord"`
	ToolbarCompact          bool                    `json:"toolbarCompact"`
	StepsPanelVisible       bool                    `json:"stepsPanelVisible"`
	StepsPanelHeight        int                     `json:"stepsPanelHeight"`
	SidebarWidth            int                     `json:"sidebarWidth"`
	RecentProjects          []string                `json:"recentProjects"`
	RecentFeatures          []string                `json:"recentFeatures"`
	SessionProject          string                  `json:"sessionProject"`
	OpenTabs                []string                `json:"openTabs"`
	ActiveTab               string                  `json:"activeTab"`
	UntitledTabs            []UntitledTabDTO        `json:"untitledTabs"`
	ScrollBeforeClick       bool                    `json:"scrollBeforeClick"`
	HoverRecordMinMs        int                     `json:"hoverRecordMinMs"`
	SelectorClickStrategies []string                `json:"selectorClickStrategies"`
	SelectorInputStrategies []string                `json:"selectorInputStrategies"`
	LibraryHeuristicsMUI    bool                    `json:"libraryHeuristicsMui"`
	LibraryHeuristicsAnt    bool                    `json:"libraryHeuristicsAnt"`
	CheckUpdatesOnStartup   bool                    `json:"checkUpdatesOnStartup"`
	Editor                  settings.EditorSettings `json:"editor"`
	ChecklistDismissed      bool                    `json:"checklistDismissed"`
	WelcomePlayedSuccess    bool                    `json:"welcomePlayedSuccess"`
	OnboardingCompleted     bool                    `json:"onboardingCompleted"`
	OnboardingDismissed     bool                    `json:"onboardingDismissed"`
	OnboardingVersion       int                     `json:"onboardingVersion"`
	StartURL                string                  `json:"startUrl"`
	RunDialogConfirmed      bool                    `json:"runDialogConfirmed"`
	PickerDuringRecording   bool                    `json:"pickerDuringRecording"`
	UILocale                string                  `json:"uiLocale"`
}

type UntitledTabDTO struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

type RunResultEntry struct {
	Path       string `json:"path"`
	Success    bool   `json:"success"`
	Status     string `json:"status,omitempty"`
	Message    string `json:"message"`
	Runner     string `json:"runner"`
	At         string `json:"at"`
	FailedStep *int   `json:"failed_step,omitempty"`
}

type FlakyScenarioDTO struct {
	Path       string `json:"path"`
	Failures   int    `json:"failures"`
	Passes     int    `json:"passes"`
	Total      int    `json:"total"`
	Flaky      bool   `json:"flaky"`
	LastFailed string `json:"last_failed_at,omitempty"`
}

type FlakyStepDTO struct {
	Path       string `json:"path"`
	Step       int    `json:"step"`
	Failures   int    `json:"failures"`
	LastFailed string `json:"last_failed_at,omitempty"`
}

type FlakyMetricsDTO struct {
	Scenarios []FlakyScenarioDTO `json:"scenarios"`
	Steps     []FlakyStepDTO     `json:"steps"`
}

func (s *Service) ListRunResults(limit int) ([]RunResultEntry, error) {
	return s.runner().ListRunResults(limit)
}

func (s *Service) FlakyMetrics(historyLimit int) (FlakyMetricsDTO, error) {
	return s.runner().FlakyMetrics(historyLimit)
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

func (s *Service) ProjectVersion() uint64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.projectVersion
}

func (s *Service) CurrentRunSession() *RunSession {
	return s.runner().CurrentSession()
}

func (s *Service) HasActiveRun() bool {
	return s.runner().HasActiveRun()
}

func (s *Service) UpdateDirtyTabsState(dirty bool) {
	if s == nil {
		return
	}
	s.liveDirtyTabs.Store(dirty)
}

func (s *Service) HasDirtyTabs() bool {
	if s == nil {
		return false
	}
	if s.liveDirtyTabs.Load() {
		return true
	}
	cfg, err := s.settingOps().LoadSettings()
	if err != nil {
		return false
	}
	if len(cfg.UntitledTabs) > 0 {
		return true
	}
	for _, path := range cfg.OpenTabs {
		draftPath, err := featureDraftPath(path)
		if err != nil {
			continue
		}
		if _, err := os.Stat(draftPath); err == nil {
			return true
		}
	}
	return false
}

func (s *Service) CloseGuardReasons() []string {
	if s == nil {
		return nil
	}
	reasons := make([]string, 0, 3)
	if s.HasActiveRun() {
		reasons = append(reasons, "active_run")
	}
	if s.HasLiveBrowser() {
		reasons = append(reasons, "active_recorder")
	}
	if s.HasDirtyTabs() {
		reasons = append(reasons, "unsaved_tabs")
	}
	return reasons
}

func (s *Service) editorAnalyzer() *EditorAnalysisService {
	s.mu.RLock()
	analysis := s.editorAnalysisService
	s.mu.RUnlock()
	if analysis != nil {
		return analysis
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.editorAnalysisService == nil {
		s.editorAnalysisService = NewEditorAnalysisService()
	}
	return s.editorAnalysisService
}

func (s *Service) buildReportService() *ReportService {
	return NewReportService(
		s.ProjectPath,
		func() string {
			if session := s.CurrentRunSession(); session != nil {
				return session.RunID
			}
			return ""
		},
		s.ReportBridgeCredentials,
	)
}

func (s *Service) reporter() *ReportService {
	s.mu.RLock()
	reporter := s.reportService
	s.mu.RUnlock()
	if reporter != nil {
		return reporter
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.reportService == nil {
		s.reportService = s.buildReportService()
	}
	return s.reportService
}

func (s *Service) runner() *RunService {
	s.mu.RLock()
	run := s.runService
	s.mu.RUnlock()
	if run != nil {
		return run
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.runService == nil {
		s.runService = NewRunService(s.ProjectPath)
	}
	return s.runService
}

func (s *Service) fileOperator() *FileOperationService {
	s.mu.RLock()
	ops := s.fileOps
	s.mu.RUnlock()
	if ops != nil {
		return ops
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.fileOps == nil {
		s.fileOps = NewFileOperationService(s.confineFeaturePath, s.ProjectPath, s.withProjectFSReadLock, s.withProjectFSWriteLock)
	}
	return s.fileOps
}

func (s *Service) recorderOps() *RecorderService {
	s.mu.RLock()
	rec := s.recorderService
	s.mu.RUnlock()
	if rec != nil {
		return rec
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.recorderService == nil {
		s.recorderService = NewRecorderService()
	}
	return s.recorderService
}

func (s *Service) settingOps() *SettingsService {
	s.mu.RLock()
	settingsSvc := s.settingsService
	s.mu.RUnlock()
	if settingsSvc != nil {
		return settingsSvc
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.settingsService == nil {
		s.settingsService = NewSettingsService(s.settingsStore)
	}
	return s.settingsService
}

func (s *Service) cliRunner() *CLIOps {
	s.mu.RLock()
	ops := s.cliOps
	s.mu.RUnlock()
	if ops != nil {
		return ops
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.cliOps == nil {
		s.cliOps = NewCLIOps()
	}
	return s.cliOps
}

func (s *Service) rotateProjectSessionLocked(path string) {
	if s.projectSession != nil && s.projectSession.cancel != nil {
		s.projectSession.cancel()
	}
	s.projectVersion++
	projectCtx, projectCancel := context.WithCancel(context.Background())
	s.projectSession = &ProjectSession{
		ID:      fmt.Sprintf("project-%d", s.projectVersion),
		Root:    path,
		Version: s.projectVersion,
		Context: projectCtx,
		cancel:  projectCancel,
	}
	s.projectPath = path
}

func (s *Service) withActivePlaywright(fn func()) {
	if s == nil || fn == nil {
		return
	}
	s.activePlaywrightCount.Add(1)
	s.activePlaywright.Add(1)
	defer s.activePlaywright.Done()
	defer s.activePlaywrightCount.Add(-1)
	fn()
}

func (s *Service) ActivePlaywrightSessions() int {
	if s == nil {
		return 0
	}
	return int(s.activePlaywrightCount.Load())
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
	// Project switch is a lifecycle boundary: retire project-scoped async work.
	if s.runService != nil {
		s.runService.CancelIfActive()
	}
	if s.validateCancel != nil {
		s.validateCancel()
		s.validateCancel = nil
	}
	if s.recorderService != nil {
		s.recorderService.Session().CancelContextOnProjectSwitch()
	}
	s.rotateProjectSessionLocked(path)
	s.mu.Unlock()
	s.CleanupStartupTempArtifacts(path)
	return s.projectInfo()
}

// RefreshProject rescans .feature files for the opened project without changing workspace state.
func (s *Service) RefreshProject() (ProjectInfo, error) {
	started := time.Now()
	defer logWailsTiming("RefreshProject", started)
	return s.projectInfo()
}

func (s *Service) projectInfo() (ProjectInfo, error) {
	s.mu.RLock()
	path := s.projectPath
	version := s.projectVersion
	projectService := s.projectService
	s.mu.RUnlock()
	if projectService == nil {
		projectService = NewProjectService(s.withProjectFSReadLock)
		s.mu.Lock()
		if s.projectService == nil {
			s.projectService = projectService
		}
		s.mu.Unlock()
	}
	return projectService.ProjectInfo(path, version)
}

func (s *Service) ReadFeature(path string) (string, error) {
	return s.fileOperator().ReadFeature(path)
}

func (s *Service) SaveFeature(path, content string) error {
	return s.fileOperator().SaveFeature(path, content)
}

func (s *Service) WriteTempFeature(content string) (string, error) {
	root := s.ProjectPath()
	var dir string
	if root != "" {
		dir = filepath.Join(root, ".scenaria", "temp", fmt.Sprintf("run-%d", time.Now().UnixNano()))
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return "", fmt.Errorf("create temp dir: %w", err)
		}
	} else {
		var err error
		dir, err = os.MkdirTemp("", "scenaria-run-")
		if err != nil {
			return "", fmt.Errorf("create temp dir: %w", err)
		}
	}
	path := filepath.Join(dir, "scenario.feature")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		_ = os.RemoveAll(dir)
		return "", fmt.Errorf("write temp feature: %w", err)
	}
	s.tempFeatureMu.Lock()
	s.tempFeatureDirs = append(s.tempFeatureDirs, dir)
	s.tempFeatureMu.Unlock()
	return path, nil
}

func (s *Service) cleanupTempFeatureDirs() {
	s.tempFeatureMu.Lock()
	dirs := s.tempFeatureDirs
	s.tempFeatureDirs = nil
	s.tempFeatureMu.Unlock()
	for _, dir := range dirs {
		_ = os.RemoveAll(dir)
	}
}

func (s *Service) InitProject() (string, error) {
	path := s.ProjectPath()
	if path == "" {
		return "", fmt.Errorf("open a project folder first")
	}
	return s.InitProjectAt(path)
}

func (s *Service) InitProjectAt(path string) (string, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return "", fmt.Errorf("project path is required")
	}
	info, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("init project: %w", err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("project path must be a directory")
	}
	return s.cliRunner().InitProject(path)
}

// EventEmitter sends Wails runtime events during long GUI operations.
type EventEmitter func(name string, payload any)

func cloneRunRequest(req RunRequest) RunRequest {
	cloned := req
	if len(req.Targets) > 0 {
		cloned.Targets = append([]string(nil), req.Targets...)
	}
	if len(req.Vars) > 0 {
		cloned.Vars = make(map[string]string, len(req.Vars))
		for k, v := range req.Vars {
			cloned.Vars[k] = v
		}
	}
	return cloned
}

func (s *Service) Run(req RunRequest, emit EventEmitter) RunResult {
	s.activePlaywrightCount.Add(1)
	s.activePlaywright.Add(1)
	defer s.activePlaywright.Done()
	defer s.activePlaywrightCount.Add(-1)
	req = cloneRunRequest(req)
	logx.Debug("run request received", "targets", len(req.Targets), "dry_run", req.DryRun)
	if len(req.Targets) == 0 && s.ProjectPath() == "" {
		return RunResult{Error: "нет файлов для запуска — откройте сценарий или проект"}
	}

	begin, err := s.runner().TryBegin(s.ProjectVersion(), req, s.tempFeatureDirsForTargets(req.Targets), DefaultRunTimeout)
	if err != nil {
		return RunResult{Error: err.Error()}
	}
	ctx := begin.Context
	runID := begin.RunID
	myGen := begin.Gen
	defer begin.Cancel()
	defer func() {
		s.cleanupTempFeatureResources(begin.TempResources)
		s.runner().Finish(myGen)
	}()

	emitRunEvent := func(name string, payload any) {
		if emit == nil {
			return
		}
		switch name {
		case "run-log-line":
			switch line := payload.(type) {
			case string:
				emit(name, map[string]any{"line": line, "runId": runID})
			case map[string]any:
				line["runId"] = runID
				emit(name, line)
			default:
				emit(name, map[string]any{"runId": runID})
			}
		case "run-progress":
			if ev, ok := payload.(player.RunProgressEvent); ok {
				emit(name, runProgressPayload(ev, runID))
				return
			}
			if m, ok := payload.(map[string]any); ok {
				m["runId"] = runID
				emit(name, m)
				return
			}
			emit(name, map[string]any{"runId": runID})
		case "run-results-changed":
			emit(name, map[string]any{"runId": runID})
		default:
			emit(name, payload)
		}
	}

	result, artifacts, err := s.runInProcess(ctx, req, emitRunEvent)
	logx.Debug("run completed", "targets", len(req.Targets), "cases", result.Scenarios, "executed", len(result.ScenarioResults))
	out := s.formatRunOutput(result, err)
	runner := resolveGUIEngine(req, req.Targets)
	if req.DryRun {
		runner = "dry-run"
	}
	entries := scenarioResultsToEntries(result.ScenarioResults, runner)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return runResultWithArtifacts(RunResult{Output: out, Error: "превышен лимит времени прогона (20 мин) — нажмите «Стоп» или упростите сценарий", Entries: entries}, artifacts)
		}
		return runResultWithArtifacts(RunResult{Output: out, Error: err.Error(), Entries: entries}, artifacts)
	}
	return runResultWithArtifacts(RunResult{Output: out, Entries: entries}, artifacts)
}

func runProgressPayload(ev player.RunProgressEvent, runID string) map[string]any {
	return map[string]any{
		"phase":       ev.Phase,
		"index":       ev.Index,
		"total":       ev.Total,
		"caseId":      ev.CaseID,
		"featurePath": ev.FeaturePath,
		"scenario":    ev.Scenario,
		"success":     ev.Success,
		"message":     ev.Message,
		"runId":       runID,
	}
}

func runResultWithArtifacts(result RunResult, artifacts report.RunArtifactLayout) RunResult {
	result.ReportPath = firstNonEmpty(artifacts.HTMLPath, result.ReportPath)
	result.HTMLPath = firstNonEmpty(artifacts.HTMLPath, result.HTMLPath)
	result.JUnitPath = firstNonEmpty(artifacts.JUnitPath, result.JUnitPath)
	result.SummaryJSON = firstNonEmpty(artifacts.SummaryJSON, result.SummaryJSON)
	result.AllureDir = firstNonEmpty(artifacts.AllureDir, result.AllureDir)
	result.TraceDir = firstNonEmpty(artifacts.TraceDir, result.TraceDir)
	result.VideoDir = firstNonEmpty(artifacts.VideoDir, result.VideoDir)
	return result
}

func (s *Service) tempFeatureDirsForTargets(targets []string) []string {
	if len(targets) == 0 {
		return nil
	}
	s.tempFeatureMu.Lock()
	defer s.tempFeatureMu.Unlock()
	known := make(map[string]struct{}, len(s.tempFeatureDirs))
	for _, dir := range s.tempFeatureDirs {
		abs, err := filepath.Abs(dir)
		if err != nil {
			continue
		}
		known[abs] = struct{}{}
	}
	seen := make(map[string]struct{})
	out := make([]string, 0)
	for _, target := range targets {
		targetAbs, err := filepath.Abs(target)
		if err != nil {
			continue
		}
		for dirAbs := range known {
			rel, err := filepath.Rel(dirAbs, targetAbs)
			if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
				continue
			}
			if _, exists := seen[dirAbs]; !exists {
				seen[dirAbs] = struct{}{}
				out = append(out, dirAbs)
			}
		}
	}
	return out
}

func (s *Service) cleanupTempFeatureResources(resources []string) {
	if len(resources) == 0 {
		return
	}
	set := make(map[string]struct{}, len(resources))
	for _, resource := range resources {
		if abs, err := filepath.Abs(resource); err == nil {
			set[abs] = struct{}{}
		}
	}
	for dir := range set {
		_ = os.RemoveAll(dir)
	}
	s.tempFeatureMu.Lock()
	kept := make([]string, 0, len(s.tempFeatureDirs))
	for _, dir := range s.tempFeatureDirs {
		abs, err := filepath.Abs(dir)
		if err != nil {
			continue
		}
		if _, remove := set[abs]; !remove {
			kept = append(kept, dir)
		}
	}
	s.tempFeatureDirs = kept
	s.tempFeatureMu.Unlock()
}

// CancelRun aborts an in-progress GUI scenario run.
func (s *Service) CancelRun() {
	s.runner().Cancel()
	s.mu.Lock()
	validateCancel := s.validateCancel
	s.mu.Unlock()
	if validateCancel != nil {
		validateCancel()
	}
}

func (s *Service) Validate(req ValidateRequest) RunResult {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultValidateTimeout)
	s.mu.Lock()
	if s.validateCancel != nil {
		s.validateCancel()
	}
	s.validateGen++
	myGen := s.validateGen
	s.validateCancel = cancel
	s.mu.Unlock()
	defer func() {
		s.mu.Lock()
		if s.validateGen == myGen {
			s.validateCancel = nil
		} else {
			logx.Debug("stale validate cleanup skipped", "expected_gen", myGen, "current_gen", s.validateGen)
		}
		s.mu.Unlock()
		cancel()
	}()

	path := s.ProjectPath()
	out, err := s.cliRunner().ValidateProject(ctx, path, req)
	if err != nil {
		return RunResult{Output: out, Error: err.Error()}
	}
	return RunResult{Output: out}
}

func (s *Service) CheckUpdate() RunResult {
	out, err := s.cliRunner().CheckUpdates()
	if err != nil {
		return RunResult{Output: out, Error: err.Error()}
	}
	return RunResult{Output: out}
}

func (s *Service) ListTestClients() ([]string, error) {
	return s.testClientOps().List()
}

func (s *Service) TestClientDetails(name string) (string, error) {
	return s.testClientOps().Details(name)
}

func (s *Service) ReadTestClientJSON(name string) (string, error) {
	return s.testClientOps().ReadJSON(name)
}

func (s *Service) SaveTestClientJSON(name, content string) error {
	return s.testClientOps().SaveJSON(name, content)
}

func (s *Service) DeleteTestClient(name string) error {
	return s.testClientOps().Delete(name)
}

func (s *Service) testClientOps() *TestClientService {
	s.mu.RLock()
	ops := s.testClientService
	s.mu.RUnlock()
	if ops != nil {
		return ops
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.testClientService == nil {
		s.testClientService = NewTestClientService(s.ProjectPath)
	}
	return s.testClientService
}

func (s *Service) LoadSettings() (AppSettingsDTO, error) {
	return s.settingOps().LoadSettings()
}

func defaultAppSettingsDTO() AppSettingsDTO {
	return AppSettingsDTO{
		Browser:               "chromium",
		ParallelWorkers:       1,
		MaxLoopIterations:     100,
		StepsPanelVisible:     true,
		StepsPanelHeight:      160,
		CheckUpdatesOnStartup: true,
		Editor:                settings.DefaultEditorSettings(),
	}
}

func appSettingsFromCfg(cfg *settings.AppSettings) AppSettingsDTO {
	height := cfg.StepsPanelHeight
	if height < 80 {
		height = 160
	}
	return AppSettingsDTO{
		Browser:                 cfg.Browser,
		Headless:                cfg.Headless,
		ParallelWorkers:         maxInt(1, cfg.ParallelWorkers),
		SlowMo:                  maxInt(0, cfg.SlowMo),
		MaxLoopIterations:       maxInt(1, cfg.MaxLoopIterations),
		NavWaitUntil:            strings.TrimSpace(cfg.NavWaitUntil),
		FilterRecording:         cfg.RecordingFilterMode,
		NavOnlyRecording:        cfg.NavOnlyRecording,
		DisableRecordURLWait:    cfg.DisableRecordURLWait,
		HoverRecord:             cfg.RecordingHoverMode,
		ToolbarCompact:          cfg.ToolbarCompact,
		StepsPanelVisible:       cfg.StepsPanelVisible,
		StepsPanelHeight:        height,
		SidebarWidth:            clampSidebarWidth(cfg.SidebarWidth),
		RecentProjects:          trimRecents(cfg.RecentProjects),
		RecentFeatures:          trimRecents(cfg.RecentFeatures),
		SessionProject:          strings.TrimSpace(cfg.SessionProject),
		OpenTabs:                trimRecents(cfg.OpenTabs),
		ActiveTab:               strings.TrimSpace(cfg.ActiveTab),
		UntitledTabs:            untitledTabsFromCfg(cfg.UntitledTabs),
		ScrollBeforeClick:       cfg.ScrollBeforeClick,
		HoverRecordMinMs:        maxInt(0, cfg.HoverRecordMinMs),
		SelectorClickStrategies: selector.NormalizeClickStrategies(cfg.SelectorClickStrategies),
		SelectorInputStrategies: selector.NormalizeInputStrategies(cfg.SelectorInputStrategies),
		LibraryHeuristicsMUI:    selector.LibraryHeuristicsMUIEnabled(cfg),
		LibraryHeuristicsAnt:    selector.LibraryHeuristicsAntEnabled(cfg),
		CheckUpdatesOnStartup:   settings.CheckUpdatesOnStartupEnabled(cfg),
		Editor:                  settings.NormalizeEditorSettings(cfg.Editor),
		ChecklistDismissed:      cfg.ChecklistDismissed,
		WelcomePlayedSuccess:    cfg.WelcomePlayedSuccess,
		OnboardingCompleted:     cfg.OnboardingCompleted,
		OnboardingDismissed:     cfg.OnboardingDismissed,
		OnboardingVersion:       cfg.OnboardingVersion,
		StartURL:                strings.TrimSpace(cfg.StartURL),
		RunDialogConfirmed:      cfg.RunDialogConfirmed,
		PickerDuringRecording:   cfg.PickerDuringRecording,
		UILocale:                normalizeUILocale(cfg.UILocale),
	}
}

func (s *Service) SaveSettings(dto AppSettingsDTO) error {
	return s.settingOps().SaveSettings(dto)
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func normalizeUILocale(raw string) string {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "en":
		return "en"
	default:
		return "ru"
	}
}

func normalizeHoverRecordMinMs(ms int) int {
	if ms <= 0 {
		return 600
	}
	return ms
}

func untitledTabsFromCfg(in []settings.UntitledTabSession) []UntitledTabDTO {
	if len(in) == 0 {
		return nil
	}
	out := make([]UntitledTabDTO, 0, len(in))
	for _, tab := range in {
		path := strings.TrimSpace(tab.Path)
		if path == "" {
			continue
		}
		out = append(out, UntitledTabDTO{Path: path, Content: tab.Content})
	}
	return out
}

func untitledTabsToCfg(in []UntitledTabDTO) []settings.UntitledTabSession {
	if len(in) == 0 {
		return nil
	}
	out := make([]settings.UntitledTabSession, 0, len(in))
	for _, tab := range in {
		path := strings.TrimSpace(tab.Path)
		if path == "" {
			continue
		}
		out = append(out, settings.UntitledTabSession{Path: path, Content: tab.Content})
	}
	return out
}
