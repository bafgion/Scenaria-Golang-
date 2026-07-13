package wailsapp

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bafgion/scenaria-golang/internal/gui"
	"github.com/bafgion/scenaria-golang/internal/paths"
	"github.com/bafgion/scenaria-golang/internal/player"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App exposes methods to the Wails frontend (Svelte).
type App struct {
	ctx context.Context
	svc *gui.Service

	otpMu          sync.Mutex
	otpCode        chan string
	otpErr         chan error
	jobSeq         atomic.Uint64
	closeConfirmed atomic.Bool
}

type projectEventEnvelope struct {
	ProjectVersion uint64 `json:"projectVersion"`
	Payload        any    `json:"payload,omitempty"`
}

func NewApp() *App {
	return &App{svc: gui.NewService()}
}

func panicRunResult(op string, r any) gui.RunResult {
	return gui.RunResult{Error: fmt.Sprintf("%s panic: %v", op, r)}
}

func (a *App) safeRunResult(op string, fn func() gui.RunResult) (result gui.RunResult) {
	defer func() {
		if r := recover(); r != nil {
			result = panicRunResult(op, r)
		}
	}()
	return fn()
}

func (a *App) safeGoRunResult(event string, fn func() gui.RunResult) {
	projectVersion := a.svc.ProjectVersion()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				a.emitProjectEventWithVersion(projectVersion, event, panicRunResult(event, r))
			}
		}()
		a.emitProjectEventWithVersion(projectVersion, event, fn())
	}()
}

func (a *App) nextJobID(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, a.jobSeq.Add(1))
}

func (a *App) startRunResultJob(prefix, finishedEvent string, fn func() gui.RunResult) string {
	projectVersion := a.svc.ProjectVersion()
	jobID := a.nextJobID(prefix)
	go func() {
		result := a.safeRunResult(prefix, fn)
		a.emitProjectEventWithVersion(projectVersion, finishedEvent, gui.AsyncRunResultDTO{JobID: jobID, Result: result})
	}()
	return jobID
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	paths.ConfigurePlaywrightBrowsers()
	player.SetEmailCodePrompt(a.promptEmailCode)
	player.SetOTPCancelHook(a.CancelOTP)
	_ = a.svc.EnsureReportBridge(a.emitReportGoto, a.emitReportRerun, a.emitReportTrace)
	a.svc.CleanupGlobalStartupTemps(os.TempDir())
}

func (a *App) emitReportGoto(req gui.ReportGotoRequest) {
	if req.ReportID == "" {
		req.ReportID = fmt.Sprintf("report:%s:%s", req.FeaturePath, req.Scenario)
	}
	a.emitProjectEvent("report-goto", req)
}

func (a *App) emitReportRerun(req gui.ReportRerunRequest) {
	if req.ReportID == "" {
		req.ReportID = fmt.Sprintf("report:%s:%s", req.FeaturePath, req.Scenario)
	}
	a.emitProjectEvent("report-rerun", req)
}

func (a *App) emitReportTrace(req gui.ReportTraceRequest) {
	if req.ReportID == "" {
		base := req.ReportDir
		if base == "" {
			base = filepath.Dir(req.TracePath)
		}
		req.ReportID = "report:" + strings.ReplaceAll(base, "\\", "/")
	}
	a.emitProjectEvent("report-trace", req)
}

// Shutdown tears down in-flight runs and browser sessions when the app exits.
func (a *App) Shutdown(ctx context.Context) {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), gui.DefaultShutdownTimeout)
	defer cancel()
	a.svc.Shutdown(shutdownCtx)
}

// BeforeClose blocks native close until the frontend flushes session state and calls ConfirmAppClose.
func (a *App) BeforeClose(ctx context.Context) bool {
	if a.closeConfirmed.Load() {
		return false
	}
	reasons := a.svc.CloseGuardReasons()
	a.emitEvent("app-close-requested", map[string]any{"reasons": reasons})
	return true
}

// ConfirmAppClose bypasses the close guard and exits the application.
func (a *App) ConfirmAppClose() {
	a.closeConfirmed.Store(true)
	if a.ctx != nil {
		runtime.Quit(a.ctx)
	}
}

func (a *App) promptEmailCode(email string) (string, error) {
	a.otpMu.Lock()
	if a.otpCode != nil || a.otpErr != nil {
		a.otpMu.Unlock()
		return "", fmt.Errorf("otp prompt already active")
	}
	a.otpCode = make(chan string, 1)
	a.otpErr = make(chan error, 1)
	a.otpMu.Unlock()

	if a.ctx != nil {
		runtime.WindowUnminimise(a.ctx)
		runtime.WindowShow(a.ctx)
	}
	a.emitProjectEvent("otp-prompt", email)

	defer a.clearOTPChannels()

	select {
	case code := <-a.otpCode:
		return code, nil
	case err := <-a.otpErr:
		return "", err
	case <-a.ctx.Done():
		return "", a.ctx.Err()
	}
}

func (a *App) clearOTPChannels() {
	a.otpMu.Lock()
	a.otpCode = nil
	a.otpErr = nil
	a.otpMu.Unlock()
}

func (a *App) SubmitOTPCode(code string) bool {
	a.otpMu.Lock()
	defer a.otpMu.Unlock()
	if a.otpCode == nil {
		return false
	}
	select {
	case a.otpCode <- code:
		return true
	default:
		return false
	}
}

func (a *App) CancelOTP() {
	a.otpMu.Lock()
	defer a.otpMu.Unlock()
	if a.otpErr == nil {
		return
	}
	select {
	case a.otpErr <- fmt.Errorf("otp cancelled"):
	default:
	}
	a.otpCode = nil
	a.otpErr = nil
}

func (a *App) Version() string {
	return a.svc.Version()
}

func (a *App) OpenProject(path string) (gui.ProjectInfo, error) {
	return a.svc.OpenProject(path)
}

func (a *App) RefreshProject() (gui.ProjectInfo, error) {
	return a.svc.RefreshProject()
}

func (a *App) ProjectPath() string {
	return a.svc.ProjectPath()
}

func (a *App) ScenariaArtifactPath(sub string) string {
	return a.svc.ScenariaArtifactPath(sub)
}

func (a *App) ReadFeature(path string) (string, error) {
	return a.svc.ReadFeature(path)
}

func (a *App) SaveFeature(path, content string) error {
	return a.svc.SaveFeature(path, content)
}

func (a *App) WriteTempFeature(content string) (string, error) {
	return a.svc.WriteTempFeature(content)
}

func (a *App) InitProject() (string, error) {
	return a.svc.InitProject()
}

func (a *App) InitProjectAt(path string) (string, error) {
	return a.svc.InitProjectAt(path)
}

func (a *App) Run(req gui.RunRequest) gui.RunResult {
	projectVersion := a.svc.ProjectVersion()
	return a.safeRunResult("run", func() gui.RunResult {
		return a.svc.Run(req, a.emitterForProjectVersion(projectVersion))
	})
}

func (a *App) StartRun(req gui.RunRequest) string {
	projectVersion := a.svc.ProjectVersion()
	return a.startRunResultJob("run", "run-finished", func() gui.RunResult {
		return a.svc.Run(req, a.emitterForProjectVersion(projectVersion))
	})
}

func (a *App) CancelRun() {
	a.svc.CancelRun()
}

func (a *App) Validate(req gui.ValidateRequest) gui.RunResult {
	return a.safeRunResult("validate", func() gui.RunResult {
		return a.svc.Validate(req)
	})
}

func (a *App) StartValidate(req gui.ValidateRequest) string {
	return a.startRunResultJob("validate", "validate-finished", func() gui.RunResult {
		return a.svc.Validate(req)
	})
}

func (a *App) ValidateFeature(text string) []gui.ValidationIssue {
	return a.svc.ValidateFeature(text)
}

func (a *App) ListTestClients() ([]string, error) {
	return a.svc.ListTestClients()
}

func (a *App) TestClientDetails(name string) (string, error) {
	return a.svc.TestClientDetails(name)
}

func (a *App) ReadTestClientJSON(name string) (string, error) {
	return a.svc.ReadTestClientJSON(name)
}

func (a *App) SaveTestClientJSON(name, content string) error {
	return a.svc.SaveTestClientJSON(name, content)
}

func (a *App) DeleteTestClient(name string) error {
	return a.svc.DeleteTestClient(name)
}

func (a *App) CaptureBrowserSession(name string) (string, error) {
	return a.svc.CaptureBrowserSession(name)
}

func (a *App) ListVanessaRunDirs(limit int) ([]string, error) {
	return a.svc.ListVanessaRunDirs(limit)
}

func (a *App) ListScenarioTitles() ([]string, error) {
	return a.svc.ListScenarioTitles()
}

func (a *App) ReadVanessaSettingsJSON() (string, error) {
	return a.svc.ReadVanessaSettingsJSON()
}

func (a *App) SaveVanessaSettingsJSON(content string) error {
	return a.svc.SaveVanessaSettingsJSON(content)
}

func (a *App) LoadProjectConfig() (gui.ProjectConfigDTO, error) {
	return a.svc.LoadProjectConfig()
}

func (a *App) SaveProjectConfig(dto gui.ProjectConfigDTO) error {
	return a.svc.SaveProjectConfig(dto)
}

func (a *App) SearchSteps(query string) []gui.StepCatalogEntry {
	return a.svc.SearchSteps(query)
}

func (a *App) DescribeEditorLine(line string) gui.StepCatalogEntry {
	entry, ok := gui.DescribeEditorLine(line)
	if !ok {
		return gui.StepCatalogEntry{}
	}
	return entry
}

func (a *App) CompletionsForLine(line string, column int, language string) gui.StepCompletionsDTO {
	return a.svc.CompletionsForLine(line, column, language)
}

func (a *App) CheckUpdate() gui.RunResult {
	return a.safeRunResult("check update", func() gui.RunResult {
		return a.svc.CheckUpdate()
	})
}

func (a *App) CheckUpdateInfo() (gui.UpdateInfoDTO, error) {
	return a.svc.CheckUpdateInfo()
}

// EventBindingTypes exposes DTOs used only in runtime.EventsEmit so wails generate keeps them in models.ts.
func (a *App) EventBindingTypes() (gui.UpdateProgressDTO, gui.VanessaRunResultDTO, player.RunProgressEvent, gui.AsyncRunResultDTO) {
	return gui.UpdateProgressDTO{}, gui.VanessaRunResultDTO{}, player.RunProgressEvent{}, gui.AsyncRunResultDTO{}
}

func (a *App) DownloadUpdate() {
	if a.ctx == nil {
		return
	}
	go func() {
		path, err := a.svc.DownloadUpdateProgress(func(p gui.UpdateProgressDTO) {
			a.emitProjectEvent("update-progress", p)
		})
		if err != nil {
			a.emitProjectEvent("update-finished", gui.RunResult{Error: err.Error()})
			return
		}
		a.emitProjectEvent("update-finished", gui.RunResult{Output: path})
	}()
}

func (a *App) ApplyUpdate() {
	if a.ctx == nil {
		return
	}
	go func() {
		err := a.svc.ApplyUpdateProgress(func(p gui.UpdateProgressDTO) {
			a.emitProjectEvent("update-progress", p)
		})
		if err != nil {
			a.emitProjectEvent("update-finished", gui.RunResult{Error: err.Error()})
			return
		}
		a.emitProjectEvent("update-finished", gui.RunResult{Output: "restart"})
		time.Sleep(900 * time.Millisecond)
		runtime.Quit(a.ctx)
	}()
}

func (a *App) OpenExternalURL(url string) error {
	return a.svc.OpenExternalURL(url)
}

func (a *App) ValidateBrowser(req gui.ValidateRequest) (issues []gui.ValidationIssue, err error) {
	defer func() {
		if r := recover(); r != nil {
			issues = nil
			err = fmt.Errorf("validate browser panic: %v", r)
		}
	}()
	return a.svc.ValidateBrowser(req)
}

func (a *App) BrowserInstallStatus(engine string) gui.BrowserInstallStatusDTO {
	return a.svc.BrowserInstallStatus(engine)
}

func (a *App) InstallBrowserEngine(engine string) gui.RunResult {
	return a.safeRunResult("install browser engine", func() gui.RunResult {
		return a.svc.InstallBrowserEngine(engine)
	})
}

func (a *App) StartInstallBrowserEngine(engine string) string {
	return a.startRunResultJob("install-browser-engine", "browser-install-finished", func() gui.RunResult {
		return a.svc.InstallBrowserEngine(engine)
	})
}

func (a *App) ListRunResults(limit int) ([]gui.RunResultEntry, error) {
	return a.svc.ListRunResults(limit)
}

func (a *App) FlakyMetrics(historyLimit int) (gui.FlakyMetricsDTO, error) {
	return a.svc.FlakyMetrics(historyLimit)
}

func (a *App) BundledExamplesPath() string {
	return a.svc.BundledExamplesPath()
}

func (a *App) ProjectArtifacts() gui.ProjectArtifacts {
	return a.svc.ProjectArtifacts()
}

func (a *App) ParseEditorSteps(text string) []gui.EditorStepRow {
	return a.svc.ParseEditorSteps(text)
}

func (a *App) AnalyzeEditorContent(text string, includeHints bool) gui.EditorAnalysisDTO {
	return a.svc.AnalyzeEditorContent(text, includeHints)
}

func (a *App) ArtifactExists(path string) bool {
	return a.svc.ArtifactExists(path)
}

func (a *App) LoadRecents() gui.RecentsDTO {
	return a.svc.LoadRecents()
}

func (a *App) RememberRecentProject(path string) error {
	return a.svc.RememberRecentProject(path)
}

func (a *App) RememberRecentFeature(path string) error {
	return a.svc.RememberRecentFeature(path)
}

func (a *App) HighlightFeature(text string) []gui.HighlightSpan {
	return a.svc.HighlightFeature(text)
}

func (a *App) RefactorUpdateStartURLs(text, newURL string) gui.RefactorResult {
	return a.svc.RefactorUpdateStartURLs(text, newURL)
}

func (a *App) RefactorNormalizeIndents(text string) string {
	return a.svc.RefactorNormalizeIndents(text)
}

func (a *App) RefactorCollapseBlankLines(text string) string {
	return a.svc.RefactorCollapseBlankLines(text)
}

func (a *App) FormatFeature(text string) string {
	return a.svc.RefactorFormatFeature(text)
}

func (a *App) RefactorReplaceInText(text, find, replace string, caseSensitive bool) gui.RefactorResult {
	return a.svc.RefactorReplaceInText(text, find, replace, caseSensitive)
}

func (a *App) AnalyzeScenarioHints(text string) []gui.ScenarioHintDTO {
	return a.svc.AnalyzeScenarioHints(text)
}

func (a *App) ApplyScenarioHintFix(req gui.ScenarioHintFixRequest) gui.RefactorResult {
	return a.svc.ApplyScenarioHintFix(req)
}

func (a *App) ResolveRunFromLine(text string, line int) (gui.RunFromLineDTO, error) {
	return a.svc.ResolveRunFromLine(text, line)
}

func (a *App) ResolveRunToLine(text string, line int) (gui.RunFromLineDTO, error) {
	return a.svc.ResolveRunToLine(text, line)
}

func (a *App) SaveFeatureDraft(featurePath, content string) error {
	return a.svc.SaveFeatureDraft(featurePath, content)
}

func (a *App) LoadFeatureDraft(featurePath string) (string, error) {
	return a.svc.LoadFeatureDraft(featurePath)
}

func (a *App) ClearFeatureDraft(featurePath string) error {
	return a.svc.ClearFeatureDraft(featurePath)
}

func (a *App) ListPlugins() ([]gui.PluginEntryDTO, error) {
	return a.svc.ListPlugins()
}

func (a *App) InstallPlugin(name, source string) error {
	return a.svc.InstallPlugin(name, source)
}

func (a *App) UninstallPlugin(name string) error {
	return a.svc.UninstallPlugin(name)
}

func (a *App) ReplaceInProject(req gui.ProjectReplaceRequest) (gui.ProjectReplaceResult, error) {
	return a.svc.ReplaceInProject(req)
}

func (a *App) DeleteFeature(path string) error {
	return a.svc.DeleteFeature(path)
}

func (a *App) DuplicateFeature(path, newName string) (string, error) {
	return a.svc.DuplicateFeature(path, newName)
}

func (a *App) MoveFeature(src, destDir string) (string, error) {
	return a.svc.MoveFeature(src, destDir)
}

func (a *App) RenameFeature(path, newName string) (string, error) {
	return a.svc.RenameFeature(path, newName)
}

func (a *App) ImportFeatures(destDir string, paths []string) ([]string, error) {
	return a.svc.ImportFeatures(destDir, paths)
}

func (a *App) LoadSettings() (gui.AppSettingsDTO, error) {
	return a.svc.LoadSettings()
}

func (a *App) SaveSettings(dto gui.AppSettingsDTO) error {
	return a.svc.SaveSettings(dto)
}

func (a *App) UpdateDirtyTabsState(dirty bool) {
	a.svc.UpdateDirtyTabsState(dirty)
}

func (a *App) Export(req gui.ExportRequest) gui.RunResult {
	return a.safeRunResult("export", func() gui.RunResult {
		return a.svc.Export(req)
	})
}

func (a *App) StartExport(req gui.ExportRequest) string {
	return a.startRunResultJob("export", "export-finished", func() gui.RunResult {
		return a.svc.Export(req)
	})
}

func (a *App) PreviewExport(text string) gui.ExportPreview {
	return a.svc.PreviewExport(text)
}

func (a *App) ImportJSON(req gui.ImportRequest) gui.RunResult {
	return a.safeRunResult("import json", func() gui.RunResult {
		return a.svc.ImportJSON(req)
	})
}

func (a *App) StartImportJSON(req gui.ImportRequest) string {
	return a.startRunResultJob("import-json", "import-json-finished", func() gui.RunResult {
		return a.svc.ImportJSON(req)
	})
}

func (a *App) RunVanessa(dryRun bool) gui.RunResult {
	return a.safeRunResult("run vanessa", func() gui.RunResult {
		return a.svc.RunVanessa(dryRun)
	})
}

func (a *App) RunPlugin(req gui.PluginRunRequest) gui.RunResult {
	return a.safeRunResult("run plugin", func() gui.RunResult {
		return a.svc.RunPlugin(req)
	})
}

func (a *App) StartRunPlugin(req gui.PluginRunRequest) string {
	return a.startRunResultJob("run-plugin", "plugin-run-finished", func() gui.RunResult {
		return a.svc.RunPlugin(req)
	})
}

func (a *App) emitEvent(name string, payload any) {
	a.emitProjectEvent(name, payload)
}

func (a *App) emitterForProjectVersion(projectVersion uint64) func(string, any) {
	return func(name string, payload any) {
		a.emitProjectEventWithVersion(projectVersion, name, payload)
	}
}

func (a *App) emitProjectEvent(name string, payload any) {
	a.emitProjectEventWithVersion(a.svc.ProjectVersion(), name, payload)
}

func (a *App) emitProjectEventWithVersion(projectVersion uint64, name string, payload any) {
	if a.ctx == nil {
		return
	}
	runtime.EventsEmit(a.ctx, name, projectEventEnvelope{
		ProjectVersion: projectVersion,
		Payload:        payload,
	})
}

func (a *App) StartVanessaRun(req gui.PluginRunRequest) {
	projectVersion := a.svc.ProjectVersion()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				a.emitProjectEventWithVersion(projectVersion, "vanessa-run-finished", gui.VanessaRunResultDTO{Error: fmt.Sprintf("vanessa run panic: %v", r)})
			}
		}()
		a.emitProjectEventWithVersion(projectVersion, "vanessa-run-started", nil)
		result := a.svc.RunVanessaPlugin(req)
		a.emitProjectEventWithVersion(projectVersion, "vanessa-run-finished", result)
	}()
}

func (a *App) PollVanessaRun(runDir string, totalPlanned int) gui.VanessaRunSnapshotDTO {
	return a.svc.PollVanessaRun(runDir, totalPlanned)
}

func (a *App) PollBrowserSession() gui.BrowserSessionDTO {
	return a.svc.PollBrowserSession()
}

func (a *App) OpenBrowser(req gui.OpenBrowserRequest) {
	projectVersion := a.svc.ProjectVersion()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				a.emitProjectEventWithVersion(projectVersion, "browser-closed", map[string]any{
					"result":           panicRunResult("open browser", r),
					"browserSessionId": a.svc.CurrentBrowserSessionID(),
				})
			}
		}()
		emit := func(name string, payload any) {
			a.emitProjectEventWithVersion(projectVersion, name, payload)
		}
		result := a.svc.OpenBrowser(req, emit)
		emit("browser-closed", map[string]any{
			"result":           result,
			"browserSessionId": firstNonEmpty(a.svc.LastClosedBrowserSessionID(), a.svc.CurrentBrowserSessionID()),
		})
	}()
}

func (a *App) StartRecord(req gui.RecordRequest) {
	projectVersion := a.svc.ProjectVersion()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				a.emitProjectEventWithVersion(projectVersion, "record-finished", map[string]any{
					"result":           panicRunResult("record", r),
					"recordSessionId":  a.svc.CurrentRecordSessionID(),
					"browserSessionId": a.svc.CurrentBrowserSessionID(),
				})
			}
		}()
		emit := func(name string, payload any) {
			a.emitProjectEventWithVersion(projectVersion, name, payload)
		}
		req.BrowseOnly = false
		result := a.svc.RecordLive(req, emit)
		if !a.svc.HasLiveBrowser() {
			emit("record-finished", map[string]any{
				"result":           result,
				"recordSessionId":  firstNonEmpty(a.svc.LastClosedRecordSessionID(), a.svc.CurrentRecordSessionID()),
				"browserSessionId": firstNonEmpty(a.svc.LastClosedBrowserSessionID(), a.svc.CurrentBrowserSessionID()),
			})
		} else if result.Error != "" {
			a.emitProjectEventWithVersion(projectVersion, "record-error", map[string]any{
				"message":          result.Error,
				"recordSessionId":  a.svc.CurrentRecordSessionID(),
				"browserSessionId": a.svc.CurrentBrowserSessionID(),
			})
		}
	}()
}

func (a *App) BeginRecordingCapture() (started bool, err error) {
	defer func() {
		if r := recover(); r != nil {
			started = false
			err = fmt.Errorf("begin recording capture panic: %v", r)
		}
	}()
	started, err = a.svc.BeginRecordingCapture()
	if err != nil {
		return false, err
	}
	payload := map[string]any{
		"append":           true,
		"recordSessionId":  a.svc.CurrentRecordSessionID(),
		"browserSessionId": a.svc.CurrentBrowserSessionID(),
		"targetPath":       a.svc.CurrentRecordTargetPath(),
	}
	if !started {
		payload["sync"] = true
	}
	a.emitEvent("record-started", payload)
	return started, nil
}

func (a *App) RecordBaseline(req gui.BaselineRecordRequest) gui.RunResult {
	return a.safeRunResult("record baseline", func() gui.RunResult {
		return a.svc.RecordBaseline(req)
	})
}

func (a *App) StartRecordBaseline(req gui.BaselineRecordRequest) string {
	return a.startRunResultJob("record-baseline", "record-baseline-finished", func() gui.RunResult {
		return a.svc.RecordBaseline(req)
	})
}

func (a *App) PauseRecording()  { a.svc.PauseRecording() }
func (a *App) ResumeRecording() { a.svc.ResumeRecording() }
func (a *App) CancelRecording() { a.svc.CancelRecording() }
func (a *App) CloseBrowser()    { a.svc.CloseBrowser() }
func (a *App) StopRecordingCapture() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("stop recording capture panic: %v", r)
		}
	}()
	stopped, err := a.svc.StopRecordingCapture()
	if err != nil {
		return err
	}
	if stopped {
		a.emitEvent("record-stopped", map[string]any{
			"reason":           "manual",
			"recordSessionId":  a.svc.CurrentRecordSessionID(),
			"browserSessionId": a.svc.CurrentBrowserSessionID(),
			"targetPath":       a.svc.CurrentRecordTargetPath(),
		})
	}
	return nil
}

func (a *App) OpenTrace(path string) gui.RunResult {
	return a.safeRunResult("open trace", func() gui.RunResult {
		return a.svc.OpenTrace(path)
	})
}

func (a *App) FailedStepLine(featurePath, scenarioName string, leafIndex int) (int, error) {
	return a.svc.FailedStepLine(featurePath, scenarioName, leafIndex)
}
func (a *App) IsRecordingPaused() bool {
	return a.svc.IsRecordingPaused()
}

func (a *App) FocusBrowser() error {
	return a.svc.FocusBrowser()
}

func (a *App) UpdateRecordingOptions(filterRecording, navOnlyRecording, hoverRecord, headless, scrollBeforeClick bool, hoverRecordMinMs int, recordURLWaitAfterClick bool) error {
	return a.svc.UpdateRecordingOptions(filterRecording, navOnlyRecording, hoverRecord, headless, scrollBeforeClick, hoverRecordMinMs, recordURLWaitAfterClick)
}

func (a *App) UndoRecordedStep() bool {
	return a.svc.UndoRecordedStep()
}

func (a *App) PickSelector() gui.PickSelectorResult {
	return a.svc.PickSelector()
}

func (a *App) PickerStepChoices(selector, keyword string) []gui.PickerStepChoice {
	return gui.PickerStepChoices(selector, keyword)
}

func (a *App) ListHTTPAuthHosts() ([]string, error) {
	return a.svc.ListHTTPAuthHosts()
}

func (a *App) HTTPAuthForHost(host string) (gui.HTTPAuthCredentials, error) {
	return a.svc.HTTPAuthForHost(host)
}

func (a *App) SaveHTTPAuth(req gui.HTTPAuthRequest) error {
	return a.svc.SaveHTTPAuth(req)
}

func (a *App) RemoveHTTPAuth(host string) error {
	return a.svc.RemoveHTTPAuth(host)
}

func (a *App) PickProjectFolder() (string, error) {
	if a.ctx == nil {
		return "", fmt.Errorf("application not ready")
	}
	return runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Выберите папку проекта Scenaria",
	})
}

func (a *App) PickSaveFile(title, defaultName string) (string, error) {
	if a.ctx == nil {
		return "", fmt.Errorf("application not ready")
	}
	return runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           title,
		DefaultFilename: defaultName,
		Filters: []runtime.FileFilter{
			{DisplayName: "Feature", Pattern: "*.feature"},
			{DisplayName: "JSON", Pattern: "*.json"},
		},
	})
}

func (a *App) PickOpenFile(title string) (string, error) {
	if a.ctx == nil {
		return "", fmt.Errorf("application not ready")
	}
	return runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: title,
		Filters: []runtime.FileFilter{
			{DisplayName: "JSON", Pattern: "*.json"},
			{DisplayName: "Feature", Pattern: "*.feature"},
		},
	})
}

func (a *App) PickOpenFiles(title string) ([]string, error) {
	if a.ctx == nil {
		return nil, fmt.Errorf("application not ready")
	}
	return runtime.OpenMultipleFilesDialog(a.ctx, runtime.OpenDialogOptions{
		Title: title,
		Filters: []runtime.FileFilter{
			{DisplayName: "Feature", Pattern: "*.feature"},
		},
	})
}

func (a *App) OpenFolder(path string) error {
	return paths.OpenWithDefaultApp(path)
}

func (a *App) ServeAllure(dir string) gui.RunResult {
	return a.safeRunResult("serve allure", func() gui.RunResult {
		return a.svc.ServeAllure(dir)
	})
}

func (a *App) StartServeAllure(dir string) string {
	return a.startRunResultJob("serve-allure", "allure-serve-finished", func() gui.RunResult {
		return a.svc.ServeAllure(dir)
	})
}

func (a *App) AllureStatus(dir string) gui.AllureStatusDTO {
	return a.svc.AllureStatus(dir)
}

func (a *App) OpenHTMLReport(path string) gui.RunResult {
	return a.safeRunResult("open html report", func() gui.RunResult {
		result := a.svc.OpenHTMLReport(path)
		if result.Error != "" {
			return result
		}
		absPath := strings.TrimSpace(result.Output)
		if absPath == "" {
			return gui.RunResult{Error: "report path is empty"}
		}
		if os.Getenv("SCENARIA_DESKTOP_SMOKE") != "" {
			return gui.RunResult{Output: absPath}
		}
		if err := paths.OpenWithDefaultApp(absPath); err != nil {
			return gui.RunResult{Error: fmt.Sprintf("open report: %v", err)}
		}
		return gui.RunResult{Output: absPath}
	})
}

// BeginSplashWindowChrome removes native title bar and system buttons during splash (Windows).
func (a *App) BeginSplashWindowChrome() {
	applySplashChrome()
}

// OpenMainWindowChrome restores the normal window frame after splash (Windows).
func (a *App) OpenMainWindowChrome() {
	applyMainChrome()
}

// CenterAppWindow centers the native app window on the current monitor.
func (a *App) CenterAppWindow() {
	if a.ctx != nil {
		runtime.WindowCenter(a.ctx)
	}
	centerAppWindow()
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}
