//go:build desktop

package desktop

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/bafgion/scenaria-golang/internal/version"
)

func showSplash(application fyne.App) fyne.Window {
	win := application.NewWindow("Scenaria")
	win.SetFixedSize(true)
	win.Resize(fyne.NewSize(420, 180))
	win.SetContent(container.NewCenter(
		container.NewVBox(
			widget.NewLabelWithStyle("Scenaria", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			widget.NewLabel(version.String()),
			widget.NewLabel("Загрузка…"),
		),
	))
	win.CenterOnScreen()
	return win
}

func revealMainAfterSplash(splash, main fyne.Window, delay time.Duration) {
	time.AfterFunc(delay, func() {
		splash.Close()
		main.Show()
	})
}
