package models

import (
	"time"
)

type PageStats struct {
	Pathname    string    `db:"pathname"`
	Views       int64     `db:"views"`
	UniqueViews int64     `db:"unique_views"`
	Entries     int64     `db:"entries"`
	Bounces     int64     `db:"bounces"`
	AvgDuration int64     `db:"avg_duration"`
	Date        time.Time `db:"date" json:"omitempty"`
}
