package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"protocol/backend/models"

	"github.com/jmoiron/sqlx"
)

type HistoryService struct {
	db *sqlx.DB
}

func NewHistoryService(db *sqlx.DB) *HistoryService {
	return &HistoryService{db: db}
}

// Get возвращает одну запись по её ID.
func (s *HistoryService) Get(ctx context.Context, id int64) (*models.History, error) {
	var item models.History
	query := `SELECT id, item_id, payload, changed_at FROM history WHERE id = $1`

	err := s.db.GetContext(ctx, &item, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("product with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return &item, nil
}

// GetAll возвращает список всех записей.
func (s *HistoryService) GetAll(ctx context.Context) ([]models.History, error) {
	var items []models.History
	query := `SELECT id, item_id, payload, changed_at FROM history ORDER BY id DESC`

	err := s.db.SelectContext(ctx, &items, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all products: %w", err)
	}

	return items, nil
}

// Update обновляет существующую запись в базе данных.
func (s *HistoryService) Update(ctx context.Context, item *models.History) error {
	query := `UPDATE history SET item_id = :item_id, payload = :payload, changed_at = :changed_at WHERE id = :id`

	result, err := s.db.NamedExecContext(ctx, query, item)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("product with id %d not found for update", item.Id)
	}

	return nil
}

// Delete удаляет запись по её ID.
func (s *HistoryService) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM history WHERE id = $1`

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
