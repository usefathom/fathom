package models

import (
	"fmt"
	"time"
)

type SiteStats struct {
	New            bool      `db:"-" json:"-" `
	SiteID         int64     `db:"site_id"`
	Visitors       int64     `db:"visitors"`
	Pageviews      int64     `db:"pageviews"`
	Sessions       int64     `db:"sessions"`
	BounceRate     float64   `db:"bounce_rate"`
	AvgDuration    float64   `db:"avg_duration"`
	KnownDurations int64     `db:"known_durations" json:",omitempty"`
	Date           time.Time `db:"date" json:",omitempty"`
}

func (s *SiteStats) FormattedDuration() string {
	return fmt.Sprintf("%d:%d", int(s.AvgDuration/60.00), (int(s.AvgDuration) % 60))
}

func (s *SiteStats) HandlePageview(p *Pageview) {
	s.Pageviews += 1

	if p.Duration > 0.00 {
		s.KnownDurations += 1
		s.AvgDuration = s.AvgDuration + ((float64(p.Duration) - s.AvgDuration) * 1 / float64(s.KnownDurations))
	}

	if p.IsNewVisitor {
		s.Visitors += 1
	}

	if p.IsNewSession {
		s.Sessions += 1

		if p.IsBounce {
			s.BounceRate = ((float64(s.Sessions-1) * s.BounceRate) + 1) / (float64(s.Sessions))
		} else {
			s.BounceRate = ((float64(s.Sessions-1) * s.BounceRate) + 0) / (float64(s.Sessions))
		}
	}
}
