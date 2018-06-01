package aggregator

import (
	log "github.com/sirupsen/logrus"
	"github.com/usefathom/fathom/pkg/models"
)

func (agg *aggregator) handlePageview(results *results, p *models.Pageview) error {
	pageStats, err := agg.getPageStats(results, p.Timestamp, p.Hostname, p.Pathname)
	if err != nil {
		log.Error(err)
		return err
	}

	pageStats.Pageviews += 1
	if p.IsUnique {
		pageStats.Visitors += 1
	}

	if p.Duration > 0.00 {
		pageStats.KnownDurations += 1
		pageStats.AvgDuration = pageStats.AvgDuration + ((float64(p.Duration) - pageStats.AvgDuration) * 1 / float64(pageStats.KnownDurations))
	}

	if p.IsNewSession {
		pageStats.Entries += 1

		if p.IsBounce {
			pageStats.BounceRate = ((float64(pageStats.Entries-1) * pageStats.BounceRate) + 1.00) / (float64(pageStats.Entries))
		} else {
			pageStats.BounceRate = ((float64(pageStats.Entries-1) * pageStats.BounceRate) + 0.00) / (float64(pageStats.Entries))
		}
	}

	return nil
}
