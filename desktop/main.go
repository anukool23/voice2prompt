// Voice2Prompt desktop app (Wails): the settings/onboarding surface plus the
// in-process dictation engine. Shares internal/{engine,config,inject,hotkey} with
// the CLI.
package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()
	err := wails.Run(&options.App{
		Title:  "Voice2Prompt",
		Width:  560,
		Height: 720,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup:         app.startup,
		OnShutdown:        app.shutdown,
		HideWindowOnClose: true, // closing the window keeps the app alive in the menu bar
		Bind:              []interface{}{app},
		Mac: &mac.Options{
			TitleBar:   mac.TitleBarHiddenInset(),
			Appearance: mac.NSAppearanceNameDarkAqua,
			About: &mac.AboutInfo{
				Title:   "Voice2Prompt",
				Message: "On-device voice dictation.",
			},
		},
	})
	if err != nil {
		println("error:", err.Error())
	}
}
