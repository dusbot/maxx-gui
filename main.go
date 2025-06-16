package main

import (
	"embed"
	"maxxgui/backend"
	"maxxgui/backend/consts"

	"github.com/dusbot/maxx/libs/slog"
	"github.com/kbinani/screenshot"
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
	const (
		MIN_WIDTH  = 1285
		MIN_HEIGHT = 850
	)
	width, height := determinScreenResolution(MIN_WIDTH, MIN_HEIGHT)
	err := wails.Run(&options.App{
		Title:       "maxx-gui",
		Width:       width,
		Height:      height,
		MinWidth:    MIN_WIDTH,
		MinHeight:   MIN_HEIGHT,
		StartHidden: false,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		Bind:      app.Bind(),
		EnumBind:  app.Enums(),
		OnStartup: app.OnStartup,
		Windows: &windows.Options{
			DisableFramelessWindowDecorations: false,
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
		slog.Printf(slog.ERROR, "An error occurred: %s", err.Error())
	}
}

func determinScreenResolution(minWidth, minHeight int) (width, height int) {
	rate := 0.8
	n := screenshot.NumActiveDisplays()
	if n > 0 {
		bounds := screenshot.GetDisplayBounds(0)
		fullWidth, fullHeight := bounds.Dx(), bounds.Dy()
		w, h := int(float64(fullWidth)*rate), int(float64(fullHeight)*rate)
		if w < minWidth {
			w = fullWidth
		}
		if h < minHeight {
			h = fullHeight
		}
		return w, h
	}
	// fallback
	return minWidth, minHeight
}
