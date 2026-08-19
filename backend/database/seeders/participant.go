package seeders

import (
	"database/sql"
)

func SeedParticipants(db *sql.DB) error {
	_, err := db.Exec(`
		INSERT INTO participants
		(full_name, cipher, school_id, class_name)
		VALUES(?, ?, ?, ?);
	`)
	return err
}
