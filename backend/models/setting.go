package models

import "time"

type Setting struct {
	Id        int       `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Title     string    `db:"title" json:"title"`
	Value     string    `db:"value" json:"value"`
	Hidden    bool      `db:"hidden" json:"hidden"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}
