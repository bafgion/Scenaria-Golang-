package main

import (
	"embed"

	"github.com/bafgion/scenaria-golang/internal/brand"
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
		Title:  brand.Name,
		Width:  560,
		Height: 500,
		MinWidth:  800,
		MinHeight: 520,
		StartHidden: true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 20, G: 22, B: 24, A: 1},
		OnStartup:        app.Startup,
		OnShutdown:       app.Shutdown,
		Bind: []interface{}{
			app,
		},
	})
	if err != nil {
		logx.Error("wails run failed", "err", err)
	}
}
