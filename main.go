package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

// Version info - injected at build time via ldflags
var version = "dev"

//go:embed all:frontend/dist
var assets embed.FS

//go:embed assets/fonts/*
var fonts embed.FS

// Fonts returns the embedded fonts filesystem for use in other packages.
func Fonts() embed.FS {
	return fonts
}

func main() {
	app := NewApp(fonts)

	err := wails.Run(&options.App{
		Title:     "Image Border Tool",
		Width:     1200,
		Height:    800,
		MinWidth:  800,
		MinHeight: 600,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 245, G: 245, B: 245, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		log.Fatalf("Failed to start application: %v", err)
	}
}
