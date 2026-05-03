//go:build !headless

package main

import (
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

// main serves as the entry point for the application.  It parses command
// line arguments to determine if an INI file was passed (for example when
// opening the application via a file association) and creates the Wails
// application with appropriate options.  The backend App is bound so that
// its exported methods are available to the frontend.  The frontend
// resources are served from the embedded assets defined in embed.go.
func main() {
	// Determine if a file path has been provided via command line args.
	var initialFile string
	if len(os.Args) > 1 {
		initialFile = os.Args[1]
	}
	app := NewApp(initialFile)

	// Create the application with options.  We set a sensible window size
	// and title.  The AssetServer is configured to serve our embedded
	// frontend files from the assets variable in embed.go.
	err := wails.Run(&options.App{
		Title:  "INIreader",
		Width:  900,
		Height: 600,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: app.startup,
		// Windows specific options.  No special configuration is needed here
		// but we keep the struct to enable future tweaks like disabling
		// console windows.
		Windows: &windows.Options{},
		Bind:    []interface{}{app},
	})
	if err != nil {
		println("Error:", err.Error())
	}
}
