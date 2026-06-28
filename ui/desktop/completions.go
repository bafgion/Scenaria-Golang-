//go:build desktop

package desktop

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/bafgion/scenaria-golang/internal/stepcatalog"
)

func (a *App) showStepCompletions() {
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Поиск шага…")

	filtered := stepcatalog.Search("")
	list := widget.NewList(
		func() int { return len(filtered) },
		func() fyne.CanvasObject {
			return container.NewVBox(widget.NewLabel("template"), widget.NewLabel("help"))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			box := item.(*fyne.Container)
			entry := filtered[id]
			box.Objects[0].(*widget.Label).SetText(entry.Template)
			box.Objects[1].(*widget.Label).SetText(entry.Category + " — " + entry.Help)
		},
	)

	refresh := func() {
		filtered = stepcatalog.Search(searchEntry.Text)
		list.Refresh()
	}
	searchEntry.OnChanged = func(string) { refresh() }

	list.OnSelected = func(id widget.ListItemID) {
		if id < 0 || id >= len(filtered) {
			return
		}
		a.tabs.InsertLine(filtered[id].Template)
	}

	content := container.NewBorder(searchEntry, nil, nil, nil, list)
	d := dialog.NewCustom("Вставить шаг", "Закрыть", content, a.window)
	d.Resize(fyne.NewSize(560, 420))
	d.Show()
}
