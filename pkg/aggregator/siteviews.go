package aggregator

import (
	log "github.com/sirupsen/logrus"
	"github.com/usefathom/fathom/pkg/models"
)

func (agg *aggregator) handleSiteview(results *results, p *models.Pageview) error {
	site, err := agg.getSiteStats(results, p.Timestamp)
	if err != nil {
		log.Error(err)
		return err
	}

	site.Pageviews += 1

	if p.Duration > 0.00 {
		site.KnownDurations += 1
		site.AvgDuration = site.AvgDuration + ((float64(p.Duration) - site.AvgDuration) * 1 / float64(site.KnownDurations))
	}

	if p.IsNewVisitor {
		site.Visitors += 1
	}

	if p.IsNewSession {
		site.Sessions += 1

		if p.IsBounce {
			site.BounceRate = ((float64(site.Sessions-1) * site.BounceRate) + 1) / (float64(site.Sessions))
		} else {
			site.BounceRate = ((float64(site.Sessions-1) * site.BounceRate) + 0) / (float64(site.Sessions))
		}
	}

	return nil
}
