package services

import (
	"fmt"

	"context"

	"github.com/jmoiron/sqlx"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type ExportService struct {
	db                 *sqlx.DB
	participantService *ParticipantService
	settingsService    *SettingsService
	schoolService      *SchoolService
	builder            *ExcelStreamBuilder
}

func NewExportService(db *sqlx.DB) *ExportService {
	participantService := NewParticipantService(db)
	settingsService := NewSettingsService(db)
	schoolService := NewSchoolService(db)
	return &ExportService{db: db, participantService: participantService,
		settingsService: settingsService, schoolService: schoolService}
}

func (s *ExportService) ExportData(ctx context.Context, fileName string) error {
	dialog := application.Get().Dialog.SaveFile()
	dialog.SetOptions(&application.SaveFileDialogOptions{
		Title: "Сохранить как ..",
	})
	dialog.AddFilter("Excel Files", "*.xlsx")
	dialog.SetFilename(fileName + ".xlsx")

	path, err := dialog.PromptForSingleSelection()

	if err != nil {
		return err
	}

	if path == "" {
		return fmt.Errorf("Сохранение отменено")
	}

	err = s.buildReport(ctx, path)
	if err != nil {
		return err
	}

	return nil
}
