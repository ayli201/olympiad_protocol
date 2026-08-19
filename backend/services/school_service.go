package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"protocol/backend/models"

	"github.com/jmoiron/sqlx"
)

type SchoolService struct {
	db *sqlx.DB
}

func NewSchoolService(db *sqlx.DB) *SchoolService {
	return &SchoolService{db: db}
}

// Get возвращает одну запись по её ID.
func (s *SchoolService) Get(ctx context.Context, id int64) (*models.School, error) {
	var item models.School
	query := `SELECT id, title FROM schools WHERE id = $1`

	err := s.db.GetContext(ctx, &item, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("School with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get school: %w", err)
	}

	return &item, nil
}

// GetAll возвращает список всех записей.
func (s *SchoolService) GetAll(ctx context.Context) ([]models.School, error) {
	var items []models.School
	query := `SELECT
			id,
			title,
			created_at
		FROM schools;`

	err := s.db.SelectContext(ctx, &items, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all schools: %w", err)
	}

	return items, nil
}

func (s *SchoolService) FindByTitle(ctx context.Context, title string) (*models.School, error) {
	var item models.School
	query := `SELECT
			id,
			title,
			created_at
		FROM schools WHERE title LIKE $1;`

	err := s.db.GetContext(ctx, &item, query, title)
	if err != nil {
		return nil, fmt.Errorf("failed to find schools by title: %w", err)
	}

	return &item, nil
}

// Update обновляет существующую запись в базе данных.
func (s *SchoolService) Update(ctx context.Context, item *models.School) error {
	query := `UPDATE schools SET title = :title WHERE id = :id`

	result, err := s.db.NamedExecContext(ctx, query, item)
	if err != nil {
		return fmt.Errorf("failed to update schools: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("School with id %d not found for update", item.Value)
	}

	return nil
}

// Delete удаляет запись по её ID.
func (s *SchoolService) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM schools WHERE id = $1`

	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete schools: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("School with id %d not found for deletion", id)
	}

	return nil
}

func (s *SchoolService) bulkDelete(ctx context.Context, tx *sqlx.Tx, ids []int) error {
	if len(ids) == 0 {
		return nil
	}
	query := "DELETE FROM schools WHERE id IN (?)"
	query, args, err := sqlx.In(query, ids)
	if err != nil {
		return fmt.Errorf("failed to expand sqlx.In: %w", err)
	}
	query = tx.Rebind(query)

	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to bulk delete schools: %w", err)
	}

	return nil
}

func (s *SchoolService) bulkUpdate(ctx context.Context, tx *sqlx.Tx, items []models.School) error {
	if len(items) == 0 {
		return nil
	}
	query := `UPDATE schools SET title = :title WHERE id = :id`
	for _, item := range items {
		_, err := tx.NamedExecContext(ctx, query, item)
		if err != nil {
			return fmt.Errorf("failed to update school: %w", err)
		}
	}
	return nil
}

func (s *SchoolService) bulkCreate(ctx context.Context, tx *sqlx.Tx, items []models.School) error {
	if len(items) == 0 {
		return nil
	}
	query := `INSERT INTO schools (title) VALUES (:title)`
	// participant.CreatedAt = time.Now()
	_, err := tx.NamedExecContext(ctx, query, items)
	if err != nil {
		return fmt.Errorf("failed to add school: %w", err)
	}
	return nil
}

func (s *SchoolService) BulkSave(ctx context.Context,
	payload GenericBulkSavePayload[models.School]) (*GenericBulkSaveResponse[models.School], error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("failed starting transaction: %w", err)
	}
	defer tx.Rollback()

	// Обработка УДАЛЕНИЙ (Delete)
	err = s.bulkDelete(ctx, tx, payload.Delete)
	if err != nil {
		return nil, fmt.Errorf("failed to bulk delete: %w", err)
	}

	// Обработка ОБНОВЛЕНИЙ (Update)
	err = s.bulkUpdate(ctx, tx, payload.Update)
	if err != nil {
		return nil, fmt.Errorf("failed to bulk update: %w", err)
	}

	// Обработка ДОБАВЛЕНИЙ (Create)
	err = s.bulkCreate(ctx, tx, payload.Create)
	if err != nil {
		return nil, fmt.Errorf("failed to bulk create: %w", err)
	}

	// Коммитим изменения в базу, если всё прошло успешно
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit tx: %w", err)
	}

	// Получаем свежие данные, чтобы вернуть их фронтенду для обновления Svelte $state
	freshData, err := s.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("saved, but failed to fetch fresh data: %w", err)
	}

	return &GenericBulkSaveResponse[models.School]{
		Success: true,
		Message: "Changes saved successfully",
		Data:    freshData,
	}, nil
}
