package models

import (
	"time"
)

type PageStats struct {
	Hostname       string    `db:"hostname"`
	Pathname       string    `db:"pathname"`
	Pageviews      int64     `db:"pageviews"`
	Visitors       int64     `db:"visitors"`
	Entries        int64     `db:"entries"`
	BounceRate     float64   `db:"bounce_rate"`
	AvgDuration    float64   `db:"avg_duration"`
	KnownDurations int64     `db:"known_durations"`
	Date           time.Time `db:"date" json:",omitempty"`
}

func (s *PageStats) HandlePageview(p *Pageview) {

	s.Pageviews += 1
	if p.IsUnique {
		s.Visitors += 1
	}

	if p.Duration > 0.00 {
		s.KnownDurations += 1
		s.AvgDuration = s.AvgDuration + ((float64(p.Duration) - s.AvgDuration) * 1 / float64(s.KnownDurations))
	}

	if p.IsNewSession {
		s.Entries += 1

		if p.IsBounce {
			s.BounceRate = ((float64(s.Entries-1) * s.BounceRate) + 1.00) / (float64(s.Entries))
		} else {
			s.BounceRate = ((float64(s.Entries-1) * s.BounceRate) + 0.00) / (float64(s.Entries))
		}
	}

}
