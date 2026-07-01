package wailsapp

import (
	"context"
	"fmt"
	"strings"
	"sync"
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

	otpMu   sync.Mutex
	otpCode chan string
	otpErr  chan error
}

func NewApp() *App {
	return &App{svc: gui.NewService()}
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	player.SetEmailCodePrompt(a.promptEmailCode)
	player.SetOTPCancelHook(a.CancelOTP)
}

func (a *App) promptEmailCode(email string) (string, error) {
	a.otpMu.Lock()
	a.otpCode = make(chan string, 1)
	a.otpErr = make(chan error, 1)
	a.otpMu.Unlock()

	runtime.EventsEmit(a.ctx, "otp-prompt", email)

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

func (a *App) Run(req gui.RunRequest) gui.RunResult {
	return a.svc.Run(req)
}

func (a *App) CancelRun() {
	a.svc.CancelRun()
}

func (a *App) Validate(req gui.ValidateRequest) gui.RunResult {
	return a.svc.Validate(req)
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

func (a *App) CompletionsForLine(line string, column int) gui.StepCompletionsDTO {
	return a.svc.CompletionsForLine(line, column)
}

func (a *App) CheckUpdate() gui.RunResult {
	return a.svc.CheckUpdate()
}

func (a *App) CheckUpdateInfo() (gui.UpdateInfoDTO, error) {
	return a.svc.CheckUpdateInfo()
}

// EventBindingTypes exposes DTOs used only in runtime.EventsEmit so wails generate keeps them in models.ts.
func (a *App) EventBindingTypes() (gui.UpdateProgressDTO, gui.VanessaRunResultDTO) {
	return gui.UpdateProgressDTO{}, gui.VanessaRunResultDTO{}
}

func (a *App) DownloadUpdate() {
	if a.ctx == nil {
		return
	}
	go func() {
		path, err := a.svc.DownloadUpdateProgress(func(p gui.UpdateProgressDTO) {
			runtime.EventsEmit(a.ctx, "update-progress", p)
		})
		if err != nil {
			runtime.EventsEmit(a.ctx, "update-finished", gui.RunResult{Error: err.Error()})
			return
		}
		runtime.EventsEmit(a.ctx, "update-finished", gui.RunResult{Output: path})
	}()
}

func (a *App) ApplyUpdate() {
	if a.ctx == nil {
		return
	}
	go func() {
		err := a.svc.ApplyUpdateProgress(func(p gui.UpdateProgressDTO) {
			runtime.EventsEmit(a.ctx, "update-progress", p)
		})
		if err != nil {
			runtime.EventsEmit(a.ctx, "update-finished", gui.RunResult{Error: err.Error()})
			return
		}
		runtime.EventsEmit(a.ctx, "update-finished", gui.RunResult{Output: "restart"})
		time.Sleep(900 * time.Millisecond)
		runtime.Quit(a.ctx)
	}()
}

func (a *App) OpenExternalURL(url string) error {
	return a.svc.OpenExternalURL(url)
}

func (a *App) ValidateBrowser(req gui.ValidateRequest) ([]gui.ValidationIssue, error) {
	return a.svc.ValidateBrowser(req)
}

func (a *App) BrowserInstallStatus(engine string) gui.BrowserInstallStatusDTO {
	return a.svc.BrowserInstallStatus(engine)
}

func (a *App) InstallBrowserEngine(engine string) gui.RunResult {
	return a.svc.InstallBrowserEngine(engine)
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
	return gui.AnalyzeScenarioHints(text)
}

func (a *App) ApplyScenarioHintFix(req gui.ScenarioHintFixRequest) gui.RefactorResult {
	return gui.ApplyScenarioHintFix(req)
}

func (a *App) ResolveRunFromLine(text string, line int) gui.RunFromLineDTO {
	result, err := gui.ResolveRunFromLine(text, line)
	if err != nil {
		return gui.RunFromLineDTO{StartStep: -1, EndStep: -1}
	}
	return result
}

func (a *App) ResolveRunToLine(text string, line int) gui.RunFromLineDTO {
	result, err := gui.ResolveRunToLine(text, line)
	if err != nil {
		return gui.RunFromLineDTO{StartStep: -1, EndStep: -1}
	}
	return result
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

func (a *App) Export(req gui.ExportRequest) gui.RunResult {
	return a.svc.Export(req)
}

func (a *App) PreviewExport(text string) gui.ExportPreview {
	return a.svc.PreviewExport(text)
}

func (a *App) ImportJSON(req gui.ImportRequest) gui.RunResult {
	return a.svc.ImportJSON(req)
}

func (a *App) RunVanessa(dryRun bool) gui.RunResult {
	return a.svc.RunVanessa(dryRun)
}

func (a *App) RunPlugin(req gui.PluginRunRequest) gui.RunResult {
	return a.svc.RunPlugin(req)
}

func (a *App) emitEvent(name string, payload any) {
	if a.ctx == nil {
		return
	}
	runtime.EventsEmit(a.ctx, name, payload)
}

func (a *App) StartVanessaRun(req gui.PluginRunRequest) {
	go func() {
		a.emitEvent("vanessa-run-started", nil)
		result := a.svc.RunVanessaPlugin(req)
		a.emitEvent("vanessa-run-finished", result)
	}()
}

func (a *App) PollVanessaRun(runDir string, totalPlanned int) gui.VanessaRunSnapshotDTO {
	return a.svc.PollVanessaRun(runDir, totalPlanned)
}

func (a *App) PollBrowserSession() gui.BrowserSessionDTO {
	return a.svc.PollBrowserSession()
}

func (a *App) OpenBrowser(req gui.OpenBrowserRequest) {
	go func() {
		emit := func(name string, payload any) {
			a.emitEvent(name, payload)
		}
		emit("browser-opened", nil)
		result := a.svc.OpenBrowser(req, emit)
		emit("browser-closed", result)
	}()
}

func (a *App) StartRecord(req gui.RecordRequest) {
	go func() {
		emit := func(name string, payload any) {
			a.emitEvent(name, payload)
		}
		req.BrowseOnly = false
		if !a.svc.HasLiveBrowser() {
			emit("record-started", map[string]any{"resume": false, "output": req.Output})
		}
		result := a.svc.RecordLive(req, emit)
		if !a.svc.HasLiveBrowser() {
			emit("record-finished", result)
		} else if result.Error != "" {
			a.emitEvent("record-error", result.Error)
		}
	}()
}

func (a *App) BeginRecordingCapture() error {
	started, err := a.svc.BeginRecordingCapture()
	if err != nil {
		return err
	}
	payload := map[string]any{"append": true}
	if !started {
		payload["sync"] = true
	}
	runtime.EventsEmit(a.ctx, "record-started", payload)
	return nil
}

func (a *App) RecordBaseline(req gui.BaselineRecordRequest) gui.RunResult {
	return a.svc.RecordBaseline(req)
}

func (a *App) PauseRecording()  { a.svc.PauseRecording() }
func (a *App) ResumeRecording() { a.svc.ResumeRecording() }
func (a *App) CancelRecording() { a.svc.CancelRecording() }
func (a *App) CloseBrowser()    { a.svc.CloseBrowser() }
func (a *App) StopRecordingCapture() error {
	err := a.svc.StopRecordingCapture()
	if err == nil {
		runtime.EventsEmit(a.ctx, "record-stopped", nil)
	}
	return err
}
func (a *App) IsRecordingPaused() bool {
	return a.svc.IsRecordingPaused()
}

func (a *App) FocusBrowser() error {
	return a.svc.FocusBrowser()
}

func (a *App) UpdateRecordingOptions(filterRecording, navOnlyRecording, hoverRecord, headless, scrollBeforeClick bool, hoverRecordMinMs int) error {
	return a.svc.UpdateRecordingOptions(filterRecording, navOnlyRecording, hoverRecord, headless, scrollBeforeClick, hoverRecordMinMs)
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
	return a.svc.ServeAllure(dir)
}

func (a *App) OpenHTMLReport(path string) gui.RunResult {
	result := a.svc.OpenHTMLReport(path)
	if result.Error != "" {
		return result
	}
	absPath := strings.TrimSpace(result.Output)
	if absPath == "" {
		return gui.RunResult{Error: "report path is empty"}
	}
	if err := paths.OpenWithDefaultApp(absPath); err != nil {
		return gui.RunResult{Error: fmt.Sprintf("open report: %v", err)}
	}
	return gui.RunResult{Output: absPath}
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
