package models

type ParticipantReport struct {
	Participant
	RawTasks   string   `db:"tasks_json" json:"-"`
	RawRank    int      `db:"raw_rank" json:"-"`
	StrictRank int      `db:"strict_rank" json:"-"`
	TieCount   int      `db:"tie_count" json:"-"`
	Total      *int     `db:"total" json:"total"`
	Percent    *float32 `db:"percent" json:"percent"`
	Rating     *string  `db:"rating" json:"rating"`
	Status     *string  `db:"status" json:"status"`
}
