package models

import (
	"time"
)

type PageStats struct {
	Pathname    string    `db:"pathname"`
	Pageviews   int64     `db:"pageviews"`
	Visitors    int64     `db:"visitors"`
	Entries     int64     `db:"entries"`
	BounceRate  float64   `db:"bounce_rate"`
	AvgDuration int64     `db:"avg_duration"`
	Date        time.Time `db:"date" json:"omitempty"`
}
