package seeders

import (
	"github.com/jmoiron/sqlx"
)

func SeedSettings(db *sqlx.DB) error {
	_, err := db.Exec(`
		INSERT INTO settings (id, name, title, value, hidden)
		VALUES (1, 'max_points', 'Максимум баллов', '35', 0),
		(2, 'tasks_count', 'Количество заданий', '10', 0),
		(3, 'winners_percent', 'Максимальный процент победителей', '15', 0),
		(4, 'winners_and_prizers_percent', 'Максимальный процент победителей и призеров', '35', 0),
		(5, 'discipline', 'Дисциплина', '', 1),
		(6, 'year_start', 'Год начала', '2026', 1),
		(7, 'year_end', 'Год окончания', '2027', 1)
		ON CONFLICT(id) DO NOTHING;
	`)
	return err
}
