//go:build desktop

package desktop

import (
	"strconv"

	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/bafgion/scenaria-golang/internal/player"
	"github.com/bafgion/scenaria-golang/internal/settings"
)

func (a *App) showSettingsDialog() {
	cfg, err := settings.LoadDefaultAppSettings()
	if err != nil || cfg == nil {
		cfg = &settings.AppSettings{Browser: "chromium", ParallelWorkers: 1, MaxLoopIterations: 100}
	}

	browser := widget.NewSelect([]string{"chromium", "firefox", "webkit"}, nil)
	browser.SetSelected(cfg.Browser)
	headless := widget.NewCheck("Headless", nil)
	headless.SetChecked(cfg.Headless)
	workers := widget.NewEntry()
	workers.SetText(strconv.Itoa(max(1, cfg.ParallelWorkers)))
	loops := widget.NewEntry()
	loops.SetText(strconv.Itoa(max(1, cfg.MaxLoopIterations)))

	dialog.ShowForm("Настройки Scenaria", "Сохранить", "Отмена", []*widget.FormItem{
		widget.NewFormItem("Браузер", browser),
		widget.NewFormItem("", headless),
		widget.NewFormItem("Параллельные воркеры", workers),
		widget.NewFormItem("Лимит циклов", loops),
	}, func(ok bool) {
		if !ok {
			return
		}
		parallel, _ := strconv.Atoi(workers.Text)
		maxLoops, _ := strconv.Atoi(loops.Text)
		out := &settings.AppSettings{
			Browser:             browser.Selected,
			Headless:            headless.Checked,
			ParallelWorkers:     max(1, parallel),
			MaxLoopIterations:   max(1, maxLoops),
			RecordingHoverMode:  cfg.RecordingHoverMode,
			RecordingFilterMode: cfg.RecordingFilterMode,
		}
		path := settings.DefaultAppSettingsPath()
		if path == "" {
			a.appendLog("Не удалось определить путь настроек")
			return
		}
		if err := settings.SaveAppSettings(path, out); err != nil {
			a.appendLog("Ошибка сохранения настроек: " + err.Error())
			return
		}
		player.SetMaxLoopIterations(out.MaxLoopIterations)
		a.appendLog("Настройки сохранены.")
	}, a.window)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
