package wailsapp

import (
	"context"
	"fmt"

	"github.com/bafgion/scenaria-golang/internal/gui"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App exposes methods to the Wails frontend (Svelte).
type App struct {
	ctx context.Context
	svc *gui.Service
}

func NewApp() *App {
	return &App{svc: gui.NewService()}
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
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

// PickProjectFolder opens a native directory picker.
func (a *App) PickProjectFolder() (string, error) {
	if a.ctx == nil {
		return "", fmt.Errorf("application not ready")
	}
	return runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Выберите папку проекта Scenaria",
	})
}
