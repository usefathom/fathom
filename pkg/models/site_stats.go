package models

import (
	"time"
)

type SiteStats struct {
	Visitors    int64     `db:"visitors"`
	Pageviews   int64     `db:"pageviews"`
	Sessions    int64     `db:"sessions"`
	BounceRate  float64   `db:"bounce_rate"`
	AvgDuration int64     `db:"avg_duration"`
	Date        time.Time `db:"date" json:"omitempty"`
}
