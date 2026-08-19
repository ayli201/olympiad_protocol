package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"protocol/backend/models"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
)

// type BulkSavePayload struct {
// 	Create []models.Participant `json:"create"`
// 	Update []models.Participant `json:"update"`
// 	Delete []int                `json:"delete"` // Массив ID для удаления
// }

// type BulkSaveResponse struct {
// 	Success bool                       `json:"success"`
// 	Message string                     `json:"message"`
// 	Data    []models.ParticipantReport `json:"data"` // Возвращаем обновленный массив строк с бэка
// }

type ParticipantService struct {
	db              *sqlx.DB
	taskService     *TaskService
	settingsService *SettingsService
}

func NewParticipantService(db *sqlx.DB) *ParticipantService {
	taskService := NewTaskService(db)
	settingsService := NewSettingsService(db)
	return &ParticipantService{db: db, taskService: taskService, settingsService: settingsService}
}

// Add добавляет нового участника.
func (s *ParticipantService) add(ctx context.Context, tx *sqlx.Tx, participant *models.Participant) error {
	query := `INSERT INTO participants (full_name, cipher, school_id, class_name)
	VALUES (:full_name, :cipher, :school_id, :class_name)`
	// participant.CreatedAt = time.Now()
	result, err := tx.NamedExecContext(ctx, query, participant)
	if err != nil {
		return fmt.Errorf("failed to add participant: %w", err)
	}

	parentID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	participant.Id = int(parentID)

	// Настраиваем внешний ключ у всех дочерних элементов
	for i := range participant.Tasks {
		participant.Tasks[i].ParticipantID = int(parentID)
		participant.Tasks[i].CreatedAt = time.Now()
	}

	err = s.taskService.BulkAdd(ctx, tx, participant.Tasks)
	if err != nil {
		return fmt.Errorf("failed to add tasks: %w", err)
	}
	// tx.Commit()
	return nil
}

func (s *ParticipantService) update(ctx context.Context, tx *sqlx.Tx, participant *models.Participant) error {
	query := "UPDATE participants SET full_name = :full_name, cipher = :cipher, school_id = :school_id, class_name = :class_name WHERE id = :id"
	_, err := tx.NamedExecContext(ctx, query, participant)
	if err != nil {
		return fmt.Errorf("failed to update participant: %w", err)
	}

	for i := range participant.Tasks {
		participant.Tasks[i].ParticipantID = int(participant.Id)
	}

	err = s.taskService.BulkUpdate(ctx, tx, participant.Tasks)
	if err != nil {
		return fmt.Errorf("failed to update tasks: %w", err)
	}

	// tx.Commit()
	return nil
}

func (s *ParticipantService) bulkDelete(ctx context.Context, tx *sqlx.Tx, ids []int) error {
	if len(ids) == 0 {
		return nil
	}

	query := "DELETE FROM participants WHERE id IN (?)"
	query, args, err := sqlx.In(query, ids)
	if err != nil {
		return fmt.Errorf("failed to expand sqlx.In: %w", err)
	}
	query = tx.Rebind(query)

	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to bulk delete via sqlx: %w", err)
	}
	return nil
}

func (s *ParticipantService) bulkUpdate(ctx context.Context, tx *sqlx.Tx, items []models.Participant) error {
	if len(items) == 0 {
		return nil
	}

	for _, item := range items {
		err := s.update(ctx, tx, &item)
		if err != nil {
			return fmt.Errorf("failed to update participant: %w", err)
		}
	}
	return nil
}

func (s *ParticipantService) bulkCreate(ctx context.Context, tx *sqlx.Tx, items []models.Participant) error {
	if len(items) == 0 {
		return nil
	}

	for _, item := range items {
		err := s.add(ctx, tx, &item)
		if err != nil {
			return fmt.Errorf("failed to add participant: %w", err)
		}
	}

	return nil
}

// Get возвращает одну запись по её ID.
func (s *ParticipantService) Get(ctx context.Context, id int64) (*models.Participant, error) {
	var item models.Participant
	query := `SELECT id, name, price FROM participants WHERE id = $1`

	err := s.db.GetContext(ctx, &item, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("participant with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get participant: %w", err)
	}

	return &item, nil
}

// GetAll возвращает список всех записей.
func (s *ParticipantService) GetAll(ctx context.Context) ([]models.ParticipantReport, error) {
	// 1. Получаем общее количество людей
	total, err := s.countParticipants(ctx)
	if err != nil || total == 0 {
		return []models.ParticipantReport{}, err
	}

	// 2. Получаем правила квот из БД
	quota := s.getQuotaRule(ctx, total)

	// 3. Делаем выборку агрегированных данных из БД
	items, err := s.fetchRankedReports(ctx, total)
	if err != nil {
		return nil, err
	}

	// 4. Проводим постобработку (бизнес-расчеты) в памяти Go
	s.calculateFinalMetrics(ctx, items, quota)

	return items, nil
}

// Update обновляет существующую запись в базе данных.
func (s *ParticipantService) Update(ctx context.Context, item *models.Participant) error {
	query := `UPDATE participants SET full_name = :full_name, cipher = :cipher, school_id = :school_id, group = :group, result = :result, percent = :percent, rating = :rating, status = :status, created_at = :created_at WHERE id = :id`

	result, err := s.db.NamedExecContext(ctx, query, item)
	if err != nil {
		return fmt.Errorf("failed to update participant: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("participant with id %d not found for update", item.Id)
	}

	return nil
}

// Delete удаляет запись по её ID.
func (s *ParticipantService) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM participants WHERE id = $1`

	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete participant: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("participant with id %d not found for deletion", id)
	}

	return nil
}

func (s *ParticipantService) calculateFinalMetrics(ctx context.Context, items []models.ParticipantReport, quota models.QuotaRule) error {
	maxWinnersCount := quota.WinnersQuota
	maxWinnersAndPrizersCount := quota.WinnersQuota + quota.WinnersAndPrizersQuota
	minWinnersPointsPercent := quota.MinWinnersPointsPercent

	participantCount := len(items)
	if quota.MaxParticipants == nil && participantCount > quota.MinParticipants {
		maxWinnersPercentSetting, err := s.settingsService.GetByName(ctx, "winners_percent")
		if err != nil {
			return err
		}
		maxWinnersPercent, err := strconv.ParseFloat(maxWinnersPercentSetting.Value, 32)
		if err != nil {
			return err
		}
		maxWinnersAndPrizersPercentSetting, err := s.settingsService.GetByName(ctx, "winners_and_prizers_percent")
		if err != nil {
			return err
		}
		maxWinnersAndPrizersPercent, err := strconv.ParseFloat(maxWinnersAndPrizersPercentSetting.Value, 32)
		if err != nil {
			return err
		}
		maxWinnersCount = int(math.Round(float64(participantCount) * maxWinnersPercent / 100))
		maxWinnersAndPrizersCount = int(math.Round(float64(participantCount) * maxWinnersAndPrizersPercent / 100))
	}

	// Трекер для деления мест (например, "1-1", "1-2")
	tieTracker := make(map[int]int, participantCount/2)

	for i := range items {
		// Парсим задачи из JSON-строки
		if items[i].RawTasks != "" && items[i].RawTasks != "[]" {
			_ = json.Unmarshal([]byte(items[i].RawTasks), &items[i].Tasks)
		} else {
			items[i].Tasks = []models.Task{}
		}

		rank := items[i].RawRank
		ties := items[i].TieCount
		strictRank := items[i].StrictRank

		// Вычисляем строку рейтинга
		if ties > 1 {
			tieTracker[rank]++
			ratingStr := fmt.Sprintf("%d-%d", rank, rank+tieTracker[rank]-1)
			items[i].Rating = &ratingStr
		} else {
			ratingStr := fmt.Sprintf("%d", rank)
			items[i].Rating = &ratingStr
		}

		// Вычисляем текстовый статус
		if strictRank <= maxWinnersCount &&
			strictRank <= maxWinnersAndPrizersCount &&
			items[i].Percent != nil &&
			*items[i].Percent >= float32(minWinnersPointsPercent) {
			status := "Победитель"
			items[i].Status = &status
		} else if strictRank <= maxWinnersAndPrizersCount {
			status := "Призер"
			items[i].Status = &status
		} else {
			status := "Участник"
			items[i].Status = &status
		}
	}
	return nil
}

// Вспомогательный метод: подсчет количества строк
func (s *ParticipantService) countParticipants(ctx context.Context) (int, error) {
	var count int
	err := s.db.GetContext(ctx, &count, countParticipantsQuery)
	return count, err
}

// Вспомогательный метод: безопасное получение квоты
func (s *ParticipantService) getQuotaRule(ctx context.Context, totalParticipants int) models.QuotaRule {
	var quota models.QuotaRule
	err := s.db.GetContext(ctx, &quota, fetchQuotaQuery, totalParticipants, totalParticipants)
	if err != nil {
		// Дефолтные значения, если таблица квот пустая
		return models.QuotaRule{WinnersQuota: 1, WinnersAndPrizersQuota: 2}
	}
	return quota
}

// Вспомогательный метод: загрузка основного списка с выделением памяти
func (s *ParticipantService) fetchRankedReports(ctx context.Context, total int) ([]models.ParticipantReport, error) {
	items := make([]models.ParticipantReport, 0, total) // Аллокация памяти
	err := s.db.SelectContext(ctx, &items, fetchRankedParticipantsQuery)
	if err != nil {
		return nil, fmt.Errorf("failed fetching reports from db: %w", err)
	}
	return items, nil
}

func (s *ParticipantService) BulkSave(ctx context.Context, payload GenericBulkSavePayload[models.Participant]) (*GenericBulkSaveResponse[models.ParticipantReport], error) {
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

	return &GenericBulkSaveResponse[models.ParticipantReport]{
		Success: true,
		Message: "Changes saved successfully",
		Data:    freshData,
	}, nil
}

func (s *ParticipantService) Clear(ctx context.Context) error {
	query := `DELETE FROM participants`

	_, err := s.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("Ошибка удаления участников: %w", err)
	}
	return nil
}
