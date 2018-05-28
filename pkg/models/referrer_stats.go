package models

import (
	"time"
)

type ReferrerStats struct {
	URL            string    `db:"url"`
	Visitors       int64     `db:"visitors"`
	Pageviews      int64     `db:"pageviews"`
	BounceRate     float64   `db:"bounce_rate"`
	AvgDuration    float64   `db:"avg_duration"`
	KnownDurations int64     `db:"known_durations"`
	Date           time.Time `db:"date" json:",omitempty"`
}
