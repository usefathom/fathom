package models

import (
	"time"
)

type ReferrerStats struct {
	URL         string    `db:"url"`
	Visitors    int64     `db:"visitors"`
	Pageviews   int64     `db:"pageviews"`
	Bounces     int64     `db:"bounces"`
	AvgDuration int64     `db:"avg_duration"`
	Date        time.Time `db:"date" json:"omitempty"`
}
