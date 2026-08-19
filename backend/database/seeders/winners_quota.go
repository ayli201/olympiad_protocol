package seeders

import (
	"github.com/jmoiron/sqlx"
)

func SeedWinnersQuota(db *sqlx.DB) error {
	_, err := db.Exec(`
		INSERT INTO quota_rules (id, min_participants, max_participants, winners_quota,
			winners_and_prizers_quota, min_winners_points_percent)
		VALUES (1, 1, 3, 1, 2, 35),
		(2, 4, 7, 1, 3, 50),
		(3, 8, 11, 2, 4, 50),
		(4, 12, 15, 2, 5, 50),
		(5, 16, 19, 3, 7, 50),
		(6, 20, 25, 4, 9, 50),
		(7, 26, NULL, 4, 10, 50)
		ON CONFLICT(id) DO NOTHING;
	`)
	return err
}
