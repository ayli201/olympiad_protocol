package services

import (
	"context"
	"fmt"
	"protocol/backend/models"

	"github.com/jmoiron/sqlx"
)

type QuotaRuleService struct {
	db *sqlx.DB
}

func NewQuotaRuleService(db *sqlx.DB) *QuotaRuleService {
	return &QuotaRuleService{db: db}
}

func (s *QuotaRuleService) GetAll(ctx context.Context) ([]models.QuotaRule, error) {
	var items []models.QuotaRule
	query := `SELECT id, min_participants, max_participants, winners_quota,
	winners_and_prizers_quota, min_winners_points_percent, created_at
	FROM quota_rules ORDER BY min_participants ASC`

	err := s.db.SelectContext(ctx, &items, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all products: %w", err)
	}

	return items, nil
}

func (s *QuotaRuleService) bulkDelete(ctx context.Context, tx *sqlx.Tx, ids []int) error {
	if len(ids) == 0 {
		return nil
	}
	query := "DELETE FROM quota_rules WHERE id IN (?)"
	query, args, err := sqlx.In(query, ids)
	if err != nil {
		return fmt.Errorf("failed to expand sqlx.In: %w", err)
	}
	query = tx.Rebind(query)

	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to bulk delete quota rules: %w", err)
	}

	return nil
}

func (s *QuotaRuleService) bulkUpdate(ctx context.Context, tx *sqlx.Tx, items []models.QuotaRule) error {
	if len(items) == 0 {
		return nil
	}
	query := `UPDATE quota_rules SET min_participants = :min_participants,
	max_participants = :max_participants, winners_quota = :winners_quota,
	winners_and_prizers_quota = :winners_and_prizers_quota WHERE id = :id`
	for _, item := range items {
		_, err := tx.NamedExecContext(ctx, query, item)
		if err != nil {
			return fmt.Errorf("failed to update school: %w", err)
		}
	}
	return nil
}

func (s *QuotaRuleService) bulkCreate(ctx context.Context, tx *sqlx.Tx, items []models.QuotaRule) error {
	if len(items) == 0 {
		return nil
	}
	query := `INSERT INTO quota_rules (min_participants, max_participants, winners_quota,
	winners_and_prizers_quota) VALUES
	(:min_participants, :max_participants, :winners_quota, :winners_and_prizers_quota)`
	_, err := tx.NamedExecContext(ctx, query, items)
	if err != nil {
		return fmt.Errorf("failed to add quota rule: %w", err)
	}
	return nil
}

func (s *QuotaRuleService) BulkSave(ctx context.Context,
	payload GenericBulkSavePayload[models.QuotaRule]) (*GenericBulkSaveResponse[models.QuotaRule], error) {
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

	return &GenericBulkSaveResponse[models.QuotaRule]{
		Success: true,
		Message: "Changes saved successfully",
		Data:    freshData,
	}, nil
}
