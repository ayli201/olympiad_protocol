package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type DraftService struct {
	db *sqlx.DB
}

func NewDraftService(db *sqlx.DB) *DraftService {
	return &DraftService{db: db}
}

func (s *DraftService) Get(ctx context.Context, entityName string) (any, error) {
	query := `
		SELECT json_changes
		FROM drafts
		WHERE entity_name = $1
	`

	var jsonStr string
	err := s.db.QueryRowContext(ctx, query, entityName).Scan(&jsonStr)
	if err == sql.ErrNoRows {
		return nil, nil // Черновик не найден
	}
	if err != nil {
		return nil, err
	}

	var target any
	if err := json.Unmarshal([]byte(jsonStr), &target); err != nil {
		return nil, fmt.Errorf("ошибка парсинга JSON черновика: %w", err)
	}

	return target, nil
}

func (s *DraftService) Save(ctx context.Context, entityName string, entityID int, payload any) (*GenericBulkSaveResponse[any], error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("ошибка сериализации черновика: %w", err)
	}

	query := `
		INSERT INTO drafts (entity_name, json_changes)
		VALUES ($1, $2)
		ON CONFLICT(entity_name) DO UPDATE SET
			json_changes = excluded.json_changes,
			updated_at = CURRENT_TIMESTAMP
	`
	_, err = s.db.ExecContext(ctx, query, entityName, string(jsonData))

	return &GenericBulkSaveResponse[any]{
		Success: true,
		Message: "Draft saved successfully",
		Data:    nil,
	}, nil
}

func (s *DraftService) Delete(ctx context.Context, entityName string) error {
	query := `DELETE FROM drafts WHERE entity_name = $1`

	result, err := s.db.ExecContext(ctx, query, entityName)
	if err != nil {
		return fmt.Errorf("failed to delete draft: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("draft with entity name %s not found for deletion", entityName)
	}

	return nil
}
