package models

import "time"

type School struct {
	Value     int       `db:"id" json:"value"`
	Label     string    `db:"title" json:"label"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}
