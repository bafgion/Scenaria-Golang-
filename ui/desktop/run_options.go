//go:build desktop

package desktop

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
	"github.com/bafgion/scenaria-golang/internal/scenario"
	"github.com/bafgion/scenaria-golang/internal/settings"
)

type runOptions struct {
	tag        string
	vars       map[string]string
	testClient string
	dryRun     bool
	headed     bool
	engine     string
	installPW  bool
}

func (a *App) showRunOptionsDialog(title string, defaults runOptions, onRun func(runOptions)) {
	tagEntry := widget.NewEntry()
	tagEntry.SetPlaceHolder("@smoke")
	tagEntry.SetText(strings.TrimSpace(defaults.tag))

	varsEntry := widget.NewMultiLineEntry()
	varsEntry.SetPlaceHolder("BASE_URL=https://example.com\nemail_code=1234")
	varsEntry.SetMinRowsVisible(3)
	if len(defaults.vars) > 0 {
		lines := make([]string, 0, len(defaults.vars))
		for key, value := range defaults.vars {
			lines = append(lines, key+"="+value)
		}
		varsEntry.SetText(strings.Join(lines, "\n"))
	}

	clients, _ := settings.ListTestClientNames(a.projectPath)
	clientOptions := append([]string{"(из feature / не задан)"}, clients...)
	clientSelect := widget.NewSelect(clientOptions, nil)
	if defaults.testClient != "" {
		clientSelect.SetSelected(defaults.testClient)
	} else {
		clientSelect.SetSelected(clientOptions[0])
	}

	dryRun := widget.NewCheck("Dry-run (без браузера)", nil)
	dryRun.SetChecked(defaults.dryRun)
	headed := widget.NewCheck("Headed (видимый браузер)", nil)
	headed.SetChecked(defaults.headed)
	installPW := widget.NewCheck("Установить Playwright при необходимости", nil)
	installPW.SetChecked(defaults.installPW)

	engineSelect := widget.NewSelect([]string{"playwright", "stub"}, nil)
	if defaults.engine != "" {
		engineSelect.SetSelected(defaults.engine)
	} else {
		engineSelect.SetSelected("playwright")
	}

	tagHints := widget.NewLabel(a.collectProjectTagsHint())

	form := container.NewVBox(
		widget.NewForm(
			widget.NewFormItem("Тег", tagEntry),
			widget.NewFormItem("TestClient", clientSelect),
			widget.NewFormItem("Движок", engineSelect),
		),
		widget.NewLabel("Переменные (NAME=VALUE, по строке):"),
		varsEntry,
		tagHints,
		dryRun,
		headed,
		installPW,
	)

	d := dialog.NewCustomConfirm(title, "Запустить", "Отмена", form, func(ok bool) {
		if !ok {
			return
		}
		opts := runOptions{
			tag:        strings.TrimSpace(tagEntry.Text),
			vars:       parseVarLines(varsEntry.Text),
			testClient: selectedTestClient(clientSelect.Selected),
			dryRun:     dryRun.Checked,
			headed:     headed.Checked,
			engine:     engineSelect.Selected,
			installPW:  installPW.Checked,
		}
		a.lastRunOptions = opts
		onRun(opts)
	}, a.window)
	d.Resize(fyne.NewSize(520, 420))
	d.Show()
}

func (a *App) collectProjectTagsHint() string {
	if a.projectPath == "" {
		return ""
	}
	store := scenario.NewFeatureStore()
	files, err := store.Discover(a.projectPath)
	if err != nil || len(files) == 0 {
		return ""
	}
	seen := map[string]struct{}{}
	tags := make([]string, 0, 8)
	for _, path := range files {
		feature, err := store.Load(path)
		if err != nil {
			continue
		}
		for _, tag := range gherkin.CollectFeatureTags(feature) {
			if _, ok := seen[tag]; ok {
				continue
			}
			seen[tag] = struct{}{}
			tags = append(tags, tag)
			if len(tags) >= 8 {
				break
			}
		}
		if len(tags) >= 8 {
			break
		}
	}
	if len(tags) == 0 {
		return "Теги в проекте: нет"
	}
	return "Теги в проекте: " + strings.Join(tags, ", ")
}

func parseVarLines(text string) map[string]string {
	out := map[string]string{}
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		eq := strings.Index(line, "=")
		if eq <= 0 {
			continue
		}
		out[strings.TrimSpace(line[:eq])] = strings.TrimSpace(line[eq+1:])
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func selectedTestClient(value string) string {
	if value == "" || value == "(из feature / не задан)" {
		return ""
	}
	return value
}

func (a *App) buildRunArgs(opts runOptions) []string {
	args := []string{a.projectPath}
	if opts.dryRun {
		args = append(args, "--dry-run")
	}
	if opts.tag != "" {
		args = append(args, "--tag", opts.tag)
	}
	if opts.testClient != "" {
		args = append(args, "--test-client", opts.testClient)
	}
	for key, value := range opts.vars {
		args = append(args, "--var", key+"="+value)
	}
	if opts.engine != "" {
		args = append(args, "--engine", opts.engine)
	}
	if opts.headed {
		args = append(args, "--headed")
	}
	if opts.installPW {
		args = append(args, "--install-playwright")
	}
	return args
}

func (a *App) showTestClientDialog() {
	if a.projectPath == "" {
		dialog.ShowInformation("Scenaria", "Сначала откройте папку проекта", a.window)
		return
	}
	names, err := settings.ListTestClientNames(a.projectPath)
	if err != nil {
		a.appendLog("TestClient: " + err.Error())
		return
	}
	if len(names) == 0 {
		dialog.ShowInformation("TestClient", "Нет файлов в .scenaria/test_clients/*.json\nСкопируйте DemoUser.json.example → DemoUser.json", a.window)
		return
	}

	details := widget.NewLabel("Выберите клиент")
	selectedID := widget.ListItemID(-1)
	list := widget.NewList(
		func() int { return len(names) },
		func() fyne.CanvasObject { return widget.NewLabel("client") },
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*widget.Label).SetText(names[id])
		},
	)
	list.OnSelected = func(id widget.ListItemID) {
		selectedID = id
		if id < 0 || id >= len(names) {
			return
		}
		client, err := settings.LoadTestClientByName(a.projectPath, names[id])
		if err != nil {
			details.SetText(err.Error())
			return
		}
		details.SetText(fmt.Sprintf("%s\nbase_url: %s\ncookies: %d\nlocal_storage: %d",
			client.Name, client.BaseURL, len(client.Cookies), len(client.LocalStorage)))
	}

	useBtn := widget.NewButton("Использовать при запуске", func() {
		if selectedID < 0 || int(selectedID) >= len(names) {
			return
		}
		a.lastRunOptions.testClient = names[selectedID]
		a.appendLog("TestClient для запуска: " + names[selectedID])
	})

	content := container.NewBorder(nil, useBtn, nil, nil, container.NewBorder(nil, details, nil, nil, list))
	d := dialog.NewCustom("TestClient", "Закрыть", content, a.window)
	d.Resize(fyne.NewSize(460, 360))
	d.Show()
}
