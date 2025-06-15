package main

import (
	"embed"
	"maxxgui/backend"
	"maxxgui/backend/consts"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := backend.NewApp()
	err := wails.Run(&options.App{
		Title:       "maxx-gui",
		Width:       1920,
		Height:      1080,
		MinWidth:    1285,
		MinHeight:   850,
		StartHidden: false,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		Bind:      app.Bind(),
		EnumBind:  app.Enums(),
		OnStartup: app.OnStartup,
		Windows: &windows.Options{
			DisableFramelessWindowDecorations: true,
			BackdropType:                      windows.None,
		},
		Mac: &mac.Options{
			TitleBar: mac.TitleBarHiddenInset(),
			About: &mac.AboutInfo{
				Title:   consts.APP_NAME,
				Message: app.Copyright(),
			},
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
