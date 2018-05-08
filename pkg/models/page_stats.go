package models

import (
	"time"
)

type PageStats struct {
	Hostname    string    `db:"hostname"`
	Pathname    string    `db:"pathname"`
	Pageviews   int64     `db:"pageviews"`
	Visitors    int64     `db:"visitors"`
	Entries     int64     `db:"entries"`
	BounceRate  float64   `db:"bounce_rate"`
	AvgDuration float64   `db:"avg_duration"`
	Date        time.Time `db:"date" json:"omitempty"`
}
