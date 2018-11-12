package models

import (
	"time"
)

type ReferrerStats struct {
	New            bool      `db:"-" json:"-"`
	SiteID         int64     `db:"site_id"`
	HostnameID     int64     `db:"hostname_id" json:"-"`
	PathnameID     int64     `db:"pathname_id" json:"-"`
	Hostname       string    `db:"hostname"`
	Pathname       string    `db:"pathname"`
	Group          string    `db:"groupname"`
	Visitors       int64     `db:"visitors"`
	Pageviews      int64     `db:"pageviews"`
	BounceRate     float64   `db:"bounce_rate"`
	AvgDuration    float64   `db:"avg_duration"`
	KnownDurations int64     `db:"known_durations"`
	Date           time.Time `db:"date" json:",omitempty"`
}

func (s *ReferrerStats) HandlePageview(p *Pageview) {
	s.Pageviews += 1

	if p.IsNewVisitor {
		s.Visitors += 1
	}

	if p.IsBounce {
		s.BounceRate = ((float64(s.Pageviews-1) * s.BounceRate) + 1.00) / (float64(s.Pageviews))
	} else {
		s.BounceRate = ((float64(s.Pageviews-1) * s.BounceRate) + 0.00) / (float64(s.Pageviews))
	}

	if p.Duration > 0.00 {
		s.KnownDurations += 1
		s.AvgDuration = s.AvgDuration + ((float64(p.Duration) - s.AvgDuration) * 1 / float64(s.KnownDurations))
	}
}
