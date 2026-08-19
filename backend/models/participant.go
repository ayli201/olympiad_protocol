package models

import "time"

type Participant struct {
	Id        int       `db:"id" json:"id"`
	FullName  string    `db:"full_name" json:"fullName"`
	Cipher    string    `db:"cipher" json:"cipher"`
	SchoolId  int       `db:"school_id" json:"schoolId"`
	ClassName string    `db:"class_name" json:"className"`
	Tasks     []Task    `db:"-" json:"tasks"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}
