package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"

	"takeout-md-fixer/internal/service"
)

// Desktop app entrypoint using Wails v3 (application, services, assets). See:
// https://v3.wails.io/

//go:embed all:frontend/dist
var assets embed.FS

// Default window size fits the stepper + review summary + options without manual resizing.
// Minimum size keeps the 720px content column and stepper usable (horizontal scroll avoided).
const (
	windowWidth     = 1000
	windowHeight    = 820
	windowMinWidth  = 720
	windowMinHeight = 640
)

func main() {
	app := application.New(application.Options{
		Name:        "Takeout Metadata Fixer",
		Description: "Fix metadata on Google Takeout media files",
		Services: []application.Service{
			application.NewService(&service.MetadataService{}),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:     "Takeout Metadata Fixer",
		Width:     windowWidth,
		Height:    windowHeight,
		MinWidth:  windowMinWidth,
		MinHeight: windowMinHeight,
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(253, 246, 236),
		URL:              "/",
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
