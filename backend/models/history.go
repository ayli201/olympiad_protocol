package models

import "time"

type History struct {
	Id        int       `db:"id" json:"id"`
	ItemId    int       `db:"item_id" json:"item_id"`
	JsonData  string    `db:"json_data" json:"json_data"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}
