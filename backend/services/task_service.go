package services

import (
	"context"
	"fmt"
	"protocol/backend/models"

	"github.com/jmoiron/sqlx"
)

type TaskService struct {
	DB *sqlx.DB
}

func NewTaskService(db *sqlx.DB) *TaskService {
	return &TaskService{DB: db}
}

func (s *TaskService) BulkAdd(ctx context.Context, tx *sqlx.Tx, items []models.Task) error {
	if len(items) > 0 {
		queryItem := `INSERT INTO tasks (participant_id, value, number) VALUES (:participant_id, :value, :number)`

		// Массовая вставка через NamedExec внутри транзакции
		_, err := tx.NamedExecContext(ctx, queryItem, items)
		if err != nil {
			return fmt.Errorf("failed to add tasks: %w", err)
		}
	}
	return nil
}

func (s *TaskService) BulkUpdate(ctx context.Context, tx *sqlx.Tx, items []models.Task) error {
	// query := `UPDATE tasks SET value = :value WHERE id = :id`
	query := `INSERT INTO tasks (participant_id, value, number)
			VALUES (:participant_id, :value, :number)
			ON CONFLICT(participant_id, number) DO UPDATE SET
				value = excluded.value`
	for _, item := range items {
		_, err := tx.NamedExecContext(ctx, query, item)
		if err != nil {
			return fmt.Errorf("failed to update task: %w", err)
		}
	}
	return nil
}

func (s *TaskService) BulkDelete(ctx context.Context, tx *sqlx.Tx, ids []int) error {
	query := `DELETE FROM tasks WHERE id IN (:ids)`
	_, err := tx.NamedExecContext(ctx, query, ids)
	if err != nil {
		return fmt.Errorf("failed to delete tasks: %w", err)
	}
	return nil
}
