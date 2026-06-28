//go:build desktop

package desktop

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/bafgion/scenaria-golang/internal/cli"
	"github.com/bafgion/scenaria-golang/internal/player"
	"github.com/bafgion/scenaria-golang/internal/scenario"
	"github.com/bafgion/scenaria-golang/internal/stepcatalog"
	"github.com/bafgion/scenaria-golang/internal/version"
)

type App struct {
	window       fyne.Window
	projectPath  string
	featureList  *widget.List
	features     []string
	logOutput    *widget.Entry
	featureEditor *widget.Entry
	selectedPath string
}

func Run() {
	application := app.NewWithID("com.bafgion.scenaria")
	application.SetMetadata(fyne.AppMetadata{
		ID:      "com.bafgion.scenaria",
		Name:    version.AppName,
		Version: version.Version,
	})
	ui := &App{
		window: application.NewWindow(version.String()),
	}
	ui.window.Resize(fyne.NewSize(1100, 720))
	ui.build()
	ui.window.ShowAndRun()
}

func (a *App) build() {
	a.logOutput = widget.NewMultiLineEntry()
	a.logOutput.SetPlaceHolder("Лог выполнения…")
	a.logOutput.Disable()

	a.featureEditor = widget.NewMultiLineEntry()
	a.featureEditor.SetPlaceHolder("Выберите .feature файл слева…")

	a.featureList = widget.NewList(
		func() int { return len(a.features) },
		func() fyne.CanvasObject { return widget.NewLabel("feature") },
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*widget.Label).SetText(filepath.Base(a.features[id]))
		},
	)
	a.featureList.OnSelected = func(id widget.ListItemID) {
		if id < 0 || id >= len(a.features) {
			return
		}
		a.loadFeature(a.features[id])
	}

	openBtn := widget.NewButton("Открыть проект…", func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil || uri == nil {
				return
			}
			a.projectPath = uri.Path()
			a.refreshFeatures()
			a.appendLog(fmt.Sprintf("Проект: %s", a.projectPath))
		}, a.window)
	})

	runBtn := widget.NewButton("Dry-run", func() {
		a.runCLI("dry-run", []string{a.projectPath, "--dry-run"})
	})

	playwrightBtn := widget.NewButton("Playwright", func() {
		a.runPlaywrightInteractive()
	})

	vaBtn := widget.NewButton("Vanessa (dry)", func() {
		a.runCLI("vanessa", []string{"run", "--project", a.projectPath, "--dry-run"})
	})

	recordBtn := widget.NewButton("Запись…", func() {
		if a.projectPath == "" {
			dialog.ShowInformation("Scenaria", "Сначала откройте папку проекта", a.window)
			return
		}
		urlEntry := widget.NewEntry()
		urlEntry.SetPlaceHolder("https://example.com")
		dialog.ShowForm("Live-запись", "Начать", "Отмена", []*widget.FormItem{
			widget.NewFormItem("URL", urlEntry),
		}, func(ok bool) {
			if !ok || strings.TrimSpace(urlEntry.Text) == "" {
				return
			}
			out := filepath.Join(a.projectPath, "recorded.feature")
			a.appendLog("Запись " + out + " …")
			if err := cli.RunRecord([]string{"--live", "--url", urlEntry.Text, "--output", out, "--idle", "20"}); err != nil {
				a.appendLog("Ошибка: " + err.Error())
				return
			}
			a.refreshFeatures()
			a.appendLog("Запись сохранена.")
		}, a.window)
	})

	validateBtn := widget.NewButton("Проверить", func() {
		a.runCLI("validate", []string{a.projectPath})
	})

	saveBtn := widget.NewButton("Сохранить feature", func() {
		if a.selectedPath == "" {
			return
		}
		if err := os.WriteFile(a.selectedPath, []byte(a.featureEditor.Text), 0o644); err != nil {
			a.appendLog("Ошибка сохранения: " + err.Error())
			return
		}
		a.appendLog("Сохранено: " + a.selectedPath)
	})

	catalogBtn := widget.NewButton("Шаги…", func() {
		items := stepcatalog.Search("")
		lines := make([]string, 0, len(items))
		for _, entry := range items {
			lines = append(lines, entry.Template+" — "+entry.Help)
		}
		dialog.ShowInformation("Каталог шагов", strings.Join(lines, "\n"), a.window)
	})

	updateBtn := widget.NewButton("Обновления", func() {
		a.appendLog("Проверка обновлений…")
		if err := cli.RunUpdate([]string{"--check"}); err != nil {
			a.appendLog("Ошибка: " + err.Error())
			return
		}
		a.appendLog("Проверка завершена.")
	})

	initBtn := widget.NewButton("Init проекта", func() {
		if a.projectPath == "" {
			dialog.ShowInformation("Scenaria", "Сначала откройте папку проекта", a.window)
			return
		}
		if err := cli.RunInit([]string{a.projectPath}); err != nil {
			a.appendLog("Ошибка init: " + err.Error())
			return
		}
		a.appendLog("Проект инициализирован (.scenaria/).")
	})

	sidebar := container.NewBorder(
		container.NewVBox(
			widget.NewLabel("Scenaria Go"),
			openBtn,
			initBtn,
			runBtn,
			playwrightBtn,
			vaBtn,
			recordBtn,
			validateBtn,
			saveBtn,
			catalogBtn,
			updateBtn,
		),
		nil, nil, nil,
		a.featureList,
	)

	editorPane := container.NewBorder(nil, nil, nil, nil, a.featureEditor)
	right := container.NewVSplit(editorPane, a.logOutput)
	right.SetOffset(0.55)

	content := container.NewHSplit(sidebar, right)
	content.SetOffset(0.28)
	a.window.SetContent(content)
}

func (a *App) runPlaywrightInteractive() {
	if a.projectPath == "" {
		dialog.ShowInformation("Scenaria", "Сначала откройте папку проекта", a.window)
		return
	}
	done := make(chan struct{})
	player.EmailCodePrompt = func(email string) (string, error) {
		codeCh := make(chan string, 1)
		errCh := make(chan error, 1)
		label := "Код из почты"
		if email != "" {
			label = "Код для " + email
		}
		entry := widget.NewEntry()
		dialog.ShowForm(label, "OK", "Отмена", []*widget.FormItem{
			widget.NewFormItem("Код", entry),
		}, func(ok bool) {
			if !ok {
				errCh <- fmt.Errorf("otp prompt cancelled")
				return
			}
			codeCh <- entry.Text
		}, a.window)
		select {
		case code := <-codeCh:
			return code, nil
		case err := <-errCh:
			return "", err
		case <-done:
			return "", fmt.Errorf("run finished")
		}
	}
	defer func() {
		close(done)
		player.EmailCodePrompt = nil
	}()
	a.runCLI("playwright", []string{a.projectPath, "--engine", "playwright", "--headed", "--install-playwright"})
}

func (a *App) runCLI(name string, args []string) {
	if a.projectPath == "" {
		dialog.ShowInformation("Scenaria", "Сначала откройте папку проекта", a.window)
		return
	}
	a.appendLog("Запуск " + name + "…")
	var err error
	switch name {
	case "validate":
		err = cli.RunValidate(args)
	case "vanessa":
		err = cli.RunVA(args)
	default:
		err = cli.RunRun(args)
	}
	if err != nil {
		a.appendLog("Ошибка: " + err.Error())
		return
	}
	a.appendLog(name + " завершён.")
}

func (a *App) loadFeature(path string) {
	payload, err := os.ReadFile(path)
	if err != nil {
		a.appendLog("Не удалось открыть: " + err.Error())
		return
	}
	a.selectedPath = path
	a.featureEditor.SetText(string(payload))
}

func (a *App) refreshFeatures() {
	store := scenario.NewFeatureStore()
	files, err := store.Discover(a.projectPath)
	if err != nil {
		a.appendLog("Не удалось прочитать проект: " + err.Error())
		return
	}
	a.features = files
	a.featureList.Refresh()
}

func (a *App) appendLog(line string) {
	a.logOutput.SetText(a.logOutput.Text + line + "\n")
}
