package models

import (
	"time"
)

type SiteStats struct {
	Visitors     int64     `db:"visitors"`
	Pageviews    int64     `db:"pageviews"`
	Bounced      int64     `db:"bounced"`
	BouncedN     int64     `db:"bounced_n"`
	AvgDuration  int64     `db:"avg_duration"`
	AvgDurationN int64     `db:"avg_duration_n"`
	Date         time.Time `db:"date"`
}
