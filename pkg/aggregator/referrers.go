package aggregator

import (
	log "github.com/sirupsen/logrus"
	"github.com/usefathom/fathom/pkg/models"
)

func (agg *aggregator) handleReferral(results *results, p *models.Pageview) error {
	hostname, pathname, err := parseUrlParts(p.Referrer)
	if err != nil {
		return err
	}

	referrerStats, err := agg.getReferrerStats(results, p.Timestamp, hostname, pathname)
	if err != nil {
		log.Error(err)
		return err
	}

	referrerStats.Pageviews += 1

	if p.IsNewVisitor {
		referrerStats.Visitors += 1
	}

	if p.IsBounce {
		referrerStats.BounceRate = ((float64(referrerStats.Pageviews-1) * referrerStats.BounceRate) + 1.00) / (float64(referrerStats.Pageviews))
	} else {
		referrerStats.BounceRate = ((float64(referrerStats.Pageviews-1) * referrerStats.BounceRate) + 0.00) / (float64(referrerStats.Pageviews))
	}

	if p.Duration > 0.00 {
		referrerStats.KnownDurations += 1
		referrerStats.AvgDuration = referrerStats.AvgDuration + ((float64(p.Duration) - referrerStats.AvgDuration) * 1 / float64(referrerStats.KnownDurations))
	}

	return nil
}
