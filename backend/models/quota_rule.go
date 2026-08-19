package models

import "time"

type QuotaRule struct {
	Id                      int       `db:"id" json:"id"`
	MinParticipants         int       `db:"min_participants" json:"minParticipants"`
	MaxParticipants         *int      `db:"max_participants" json:"maxParticipants"`
	WinnersQuota            int       `db:"winners_quota" json:"winnersQuota"`
	WinnersAndPrizersQuota  int       `db:"winners_and_prizers_quota" json:"winnersAndPrizersQuota"`
	MinWinnersPointsPercent int       `db:"min_winners_points_percent" json:"minWinnersPointsPercent"`
	CreatedAt               time.Time `db:"created_at" json:"createdAt"`
}
