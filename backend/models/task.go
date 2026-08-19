package models

import "time"

type Task struct {
	Id            int       `db:"id" json:"id"`
	Value         float32   `db:"value" json:"value"`
	Number        int       `db:"number" json:"number"`
	ParticipantID int       `db:"participant_id" json:"participantId"`
	CreatedAt     time.Time `db:"created_at" json:"createdAt"`
}
