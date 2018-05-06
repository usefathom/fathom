package models

import (
	"time"
)

type PageStats struct {
	Pathname     string    `db:"pathname"`
	Views        int64     `db:"views"`
	UniqueViews  int64     `db:"unique_views"`
	Bounced      int64     `db:"bounced"`
	BouncedN     int64     `db:"bounced_n"`
	AvgDuration  int64     `db:"avg_duration"`
	AvgDurationN int64     `db:"avg_duration_n"`
	Date         time.Time `db:"date"`
}
