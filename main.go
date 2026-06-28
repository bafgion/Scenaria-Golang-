package main

import (
	"embed"

	"github.com/bafgion/scenaria-golang/internal/logx"
	"github.com/bafgion/scenaria-golang/internal/wailsapp"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	logx.Init()
	app := wailsapp.NewApp()
	err := wails.Run(&options.App{
		Title:  "Scenaria",
		Width:  1280,
		Height: 800,
		MinWidth:  960,
		MinHeight: 640,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 24, G: 26, B: 32, A: 1},
		OnStartup:        app.Startup,
		Bind: []interface{}{
			app,
		},
	})
	if err != nil {
		logx.Error("wails run failed", "err", err)
	}
}
