package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"protocol/backend/models"

	"github.com/jmoiron/sqlx"
)

type SettingsService struct {
	db *sqlx.DB
}

func NewSettingsService(db *sqlx.DB) *SettingsService {
	return &SettingsService{db: db}
}

// Get возвращает одну запись по её ID.
func (s *SettingsService) Get(ctx context.Context, id int64) (*models.Setting, error) {
	var item models.Setting
	query := `SELECT id, name, title, value, hidden FROM settings WHERE id = $1`

	err := s.db.GetContext(ctx, &item, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("product with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return &item, nil
}

func (s *SettingsService) GetByName(ctx context.Context, name string) (*models.Setting, error) {
	var item models.Setting
	query := `SELECT id, name, title, value, hidden FROM settings WHERE name = $1`

	err := s.db.GetContext(ctx, &item, query, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("setting with name %s not found", name)
		}
		return nil, fmt.Errorf("failed to get setting: %w", err)
	}

	return &item, nil
}

// GetAll возвращает список всех записей.
func (s *SettingsService) GetAllDefaults(ctx context.Context) ([]models.Setting, error) {
	var items []models.Setting
	query := `SELECT id, name, title, value FROM settings ORDER BY id DESC`

	err := s.db.SelectContext(ctx, &items, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all products: %w", err)
	}

	return items, nil
}

// GetAll возвращает список всех записей.
func (s *SettingsService) GetAllVisible(ctx context.Context) ([]models.Setting, error) {
	var items []models.Setting
	query := `SELECT id, name, title, value FROM settings WHERE hidden = 0 ORDER BY id DESC`

	err := s.db.SelectContext(ctx, &items, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all products: %w", err)
	}

	return items, nil
}

func (s *SettingsService) GetAll(ctx context.Context) ([]models.Setting, error) {
	var items []models.Setting
	query := `SELECT id, name, title, value FROM settings ORDER BY id DESC`

	err := s.db.SelectContext(ctx, &items, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all products: %w", err)
	}

	return items, nil
}

// Update обновляет существующую запись в базе данных.
func (s *SettingsService) Update(ctx context.Context, item *models.Setting) error {
	query := `UPDATE settings SET name = :name, title = :title, value = :value,
	hidden = :hidden WHERE id = :id`
	_, err := s.db.NamedExecContext(ctx, query, item)
	if err != nil {
		return fmt.Errorf("failed to update settings: %w", err)
	}

	return nil
}

// Delete удаляет запись по её ID.
func (s *SettingsService) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM settings WHERE id = $1`

	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("product with id %d not found for deletion", id)
	}

	return nil
}

func (s *SettingsService) bulkUpdate(ctx context.Context, tx *sqlx.Tx, items []models.Setting) error {
	if len(items) == 0 {
		return nil
	}
	query := `UPDATE settings SET name = :name, title = :title, value = :value,
	hidden = :hidden WHERE id = :id`
	for _, item := range items {
		_, err := tx.NamedExecContext(ctx, query, item)
		if err != nil {
			return fmt.Errorf("failed to update settings: %w", err)
		}
	}
	return nil
}

func (s *SettingsService) BulkSave(ctx context.Context,
	payload GenericBulkSavePayload[models.Setting]) (*GenericBulkSaveResponse[models.Setting], error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("failed starting transaction: %w", err)
	}
	defer tx.Rollback()

	// Обработка ОБНОВЛЕНИЙ (Update)
	err = s.bulkUpdate(ctx, tx, payload.Update)
	if err != nil {
		return nil, fmt.Errorf("failed to bulk update: %w", err)
	}

	// Коммитим изменения в базу, если всё прошло успешно
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit tx: %w", err)
	}

	// Получаем свежие данные, чтобы вернуть их фронтенду для обновления Svelte $state
	freshData, err := s.GetAllVisible(ctx)
	if err != nil {
		return nil, fmt.Errorf("saved, but failed to fetch fresh data: %w", err)
	}

	return &GenericBulkSaveResponse[models.Setting]{
		Success: true,
		Message: "Changes saved successfully",
		Data:    freshData,
	}, nil
}
