package services

const (
	// Запрос для быстрого подсчета участников
	countParticipantsQuery = `SELECT COUNT(*) FROM participants;`

	// Запрос для получения квот по индексу
	fetchQuotaQuery = `
		SELECT winners_quota, winners_and_prizers_quota, min_winners_points_percent
		FROM quota_rules
		WHERE ? >= min_participants
		  AND ? <= COALESCE(max_participants, 999999)
		LIMIT 1;`

	// Основной аналитический запрос с агрегацией и оконными функциями
	fetchRankedParticipantsQuery = `
		WITH ordered_tasks AS (
			SELECT id, participant_id, value, number
			FROM tasks
			ORDER BY number ASC
		), tasks_count AS (
			SELECT CAST(value AS INTEGER) AS value FROM settings WHERE name = 'tasks_count' LIMIT 1
		), participant_aggregates AS (
			SELECT
				p.id, p.full_name, p.cipher, p.school_id, p.class_name, p.created_at,
				COALESCE(SUM(t.value), 0) AS total,
				'[' || COALESCE(GROUP_CONCAT(json_object('id', t.id, 'number', t.number, 'value', t.value)), '') || ']' AS tasks_json
			FROM participants p
			LEFT JOIN ordered_tasks t ON t.participant_id = p.id
			WHERE t.number <= (SELECT value FROM tasks_count)
			GROUP BY p.id
		),
		max_limit AS (
			SELECT CAST(value AS REAL) AS max_val FROM settings WHERE name = 'max_points' LIMIT 1
		)
		SELECT
			a.id, a.full_name, a.cipher, a.school_id, a.class_name, a.created_at, a.tasks_json, a.total,
			ROUND((a.total / (SELECT max_val FROM max_limit)) * 100, 2) AS percent,
			RANK() OVER (ORDER BY a.total DESC) AS raw_rank,
			ROW_NUMBER() OVER (ORDER BY a.total DESC, a.full_name ASC) AS strict_rank,
			COUNT(*) OVER (PARTITION BY a.total) AS tie_count
		FROM participant_aggregates a
		ORDER BY raw_rank ASC, a.full_name ASC;`
)
