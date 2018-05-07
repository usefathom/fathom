package models

import (
	"time"
)

type SiteStats struct {
	Visitors    int64     `db:"visitors"`
	Pageviews   int64     `db:"pageviews"`
	Sessions    int64     `db:"sessions"`
	Bounces     int64     `db:"bounces"`
	AvgDuration int64     `db:"avg_duration"`
	Date        time.Time `db:"date"`
}
