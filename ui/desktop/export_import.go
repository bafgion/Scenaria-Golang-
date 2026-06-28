//go:build desktop

package desktop

import (
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/bafgion/scenaria-golang/internal/cli"
)

func (a *App) showExportDialog() {
	path := a.tabs.CurrentPath()
	if path == "" {
		dialog.ShowInformation("Scenaria", "Откройте feature-файл для экспорта", a.window)
		return
	}
	formatSelect := widget.NewSelect([]string{"json", "feature", "ts", "python"}, nil)
	formatSelect.SetSelected("json")
	outEntry := widget.NewEntry()
	outEntry.SetPlaceHolder("путь к выходному файлу")
	baseEntry := widget.NewEntry()
	baseEntry.SetPlaceHolder("https://example.com")
	dialog.ShowForm("Экспорт сценария", "Экспорт", "Отмена", []*widget.FormItem{
		widget.NewFormItem("Формат", formatSelect),
		widget.NewFormItem("Файл", outEntry),
		widget.NewFormItem("Base URL", baseEntry),
	}, func(ok bool) {
		if !ok || strings.TrimSpace(outEntry.Text) == "" {
			return
		}
		args := []string{path, "--output", outEntry.Text, "--format", formatSelect.Selected}
		if strings.TrimSpace(baseEntry.Text) != "" {
			args = append(args, "--base-url", baseEntry.Text)
		}
		a.appendLog("Экспорт " + path + " …")
		if err := cli.RunExport(args); err != nil {
			a.appendLog("Ошибка: " + err.Error())
			return
		}
		a.appendLog("Экспортировано: " + outEntry.Text)
	}, a.window)
}

func (a *App) showImportJSONDialog() {
	if a.projectPath == "" {
		dialog.ShowInformation("Scenaria", "Сначала откройте папку проекта", a.window)
		return
	}
	dialog.ShowFileOpen(func(uri fyne.URIReadCloser, err error) {
		if err != nil || uri == nil {
			return
		}
		jsonPath := uri.URI().Path()
		_ = uri.Close()
		outEntry := widget.NewEntry()
		outEntry.SetText(filepath.Join(a.projectPath, strings.TrimSuffix(filepath.Base(jsonPath), ".json")+".feature"))
		dialog.ShowForm("Импорт JSON", "Импорт", "Отмена", []*widget.FormItem{
			widget.NewFormItem("Feature", outEntry),
		}, func(ok bool) {
			if !ok || strings.TrimSpace(outEntry.Text) == "" {
				return
			}
			a.appendLog("Импорт " + jsonPath + " …")
			if err := cli.RunImportJSON([]string{jsonPath, "--output", outEntry.Text}); err != nil {
				a.appendLog("Ошибка: " + err.Error())
				return
			}
			a.refreshFeatures()
			a.loadFeature(outEntry.Text)
			a.appendLog("Импортировано: " + outEntry.Text)
		}, a.window)
	}, a.window)
}
