package backend

import (
	"context"
	"embed"
	"log"
	"protocol/backend/database"
	"protocol/backend/services"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

func CreateApp(assets embed.FS) *application.App {
	// 1. Инициализируем БД и применяем миграции
	db, err := database.InitDB("Olympiad's protocol") // Имя папки в AppData/Config
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	// defer db.Close()

	importService := services.NewImportService(db)

	// Create a new Wails application by providing the necessary options.
	// Variables 'Name' and 'Description' are for application metadata.
	// 'Assets' configures the asset server with the 'FS' variable pointing to the frontend files.
	// 'Bind' is a list of Go struct instances. The frontend has access to the methods of these instances.
	// 'Mac' options tailor the application when running an macOS.
	app := application.New(application.Options{
		Name:        "protocol_wails3",
		Description: "Protocol",
		Services: []application.Service{
			application.NewService(services.NewParticipantService(db)),
			application.NewService(services.NewSchoolService(db)),
			application.NewService(importService),
			application.NewService(services.NewExportService(db)),
			application.NewService(services.NewSettingsService(db)),
			application.NewService(services.NewQuotaRuleService(db)),
			application.NewService(services.NewDraftService(db)),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	// Create a new window with the necessary options.
	// 'Title' is the title of the window.
	// 'Mac' options tailor the window when running on macOS.
	// 'BackgroundColour' is the background colour of the window.
	// 'URL' is the URL that will be loaded into the webview.
	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:          "Протокол олимпиады",
		Width:          1400,
		Height:         768,
		EnableFileDrop: true,
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/",
	})

	window.OnWindowEvent(events.Common.WindowFilesDropped, func(event *application.WindowEvent) {
		files := event.Context().DroppedFiles()
		if len(files) > 0 {
			ctx := context.Background()
			participants, err := importService.ImportDataFromFile(ctx, files[0])
			if err != nil {
				app.Event.Emit("protocol:data-imported-event", participants)
				return
			}
			app := application.Get()
			app.Event.Emit("protocol:data-imported-event", participants)
		}
	})

	return app
}
