package wailsapp

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	goruntime "runtime"
	"sync"

	"github.com/bafgion/scenaria-golang/internal/gui"
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
	player.EmailCodePrompt = a.promptEmailCode
}

func (a *App) promptEmailCode(email string) (string, error) {
	a.otpMu.Lock()
	a.otpCode = make(chan string, 1)
	a.otpErr = make(chan error, 1)
	a.otpMu.Unlock()

	runtime.EventsEmit(a.ctx, "otp-prompt", email)

	select {
	case code := <-a.otpCode:
		return code, nil
	case err := <-a.otpErr:
		return "", err
	case <-a.ctx.Done():
		return "", a.ctx.Err()
	}
}

func (a *App) SubmitOTPCode(code string) {
	a.otpMu.Lock()
	defer a.otpMu.Unlock()
	if a.otpCode != nil {
		a.otpCode <- code
	}
}

func (a *App) CancelOTP() {
	a.otpMu.Lock()
	defer a.otpMu.Unlock()
	if a.otpErr != nil {
		a.otpErr <- fmt.Errorf("otp cancelled")
	}
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

func (a *App) ReadFeature(path string) (string, error) {
	return a.svc.ReadFeature(path)
}

func (a *App) SaveFeature(path, content string) error {
	return a.svc.SaveFeature(path, content)
}

func (a *App) InitProject() (string, error) {
	return a.svc.InitProject()
}

func (a *App) Run(req gui.RunRequest) gui.RunResult {
	return a.svc.Run(req)
}

func (a *App) Validate(browser string, skipBrowser bool) gui.RunResult {
	return a.svc.Validate(browser, skipBrowser)
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

func (a *App) SearchSteps(query string) []gui.StepCatalogEntry {
	return a.svc.SearchSteps(query)
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

func (a *App) ImportJSON(req gui.ImportRequest) gui.RunResult {
	return a.svc.ImportJSON(req)
}

func (a *App) RunVanessa(dryRun bool) gui.RunResult {
	return a.svc.RunVanessa(dryRun)
}

func (a *App) StartRecord(req gui.RecordRequest) {
	go func() {
		runtime.EventsEmit(a.ctx, "record-started", req.Output)
		result := a.svc.RecordLive(req)
		runtime.EventsEmit(a.ctx, "record-finished", result)
	}()
}

func (a *App) PauseRecording()  { a.svc.PauseRecording() }
func (a *App) ResumeRecording() { a.svc.ResumeRecording() }
func (a *App) CancelRecording() { a.svc.CancelRecording() }
func (a *App) IsRecordingPaused() bool {
	return a.svc.IsRecordingPaused()
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

func (a *App) OpenFolder(path string) error {
	path = filepath.Clean(path)
	if path == "" {
		return fmt.Errorf("path is required")
	}
	switch goruntime.GOOS {
	case "windows":
		return exec.Command("explorer", path).Start()
	case "darwin":
		return exec.Command("open", path).Start()
	default:
		return exec.Command("xdg-open", path).Start()
	}
}
