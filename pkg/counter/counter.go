package counter

import (
	"time"

	"github.com/usefathom/fathom/pkg/datastore"
	"github.com/usefathom/fathom/pkg/models"
)

func Aggregate() error {

	// Get unprocessed pageviews
	pageviews, err := datastore.GetRawPageviews()
	if err != nil && err != datastore.ErrNoResults {
		return err
	}

	//  Do we have anything to process?
	if len(pageviews) == 0 {
		return nil
	}

	// site stats
	date := time.Now()
	siteStats, err := getSiteStats(date)
	if err != nil {
		return err
	}

	for _, p := range pageviews {
		siteStats.Pageviews += 1

		if p.IsNewVisitor {
			siteStats.Visitors += 1
			siteStats.BouncedN += 1

			if p.IsBounce {
				siteStats.Bounced += 1
			}

			// TODO: duration
		}
	}

	// TODO: page stats

	// TODO: referrer stats

	err = datastore.SaveSiteStats(siteStats)
	if err != nil {
		return err
	}

	// TODO: delete data
	return nil
}

func getSiteStats(date time.Time) (*models.SiteStats, error) {
	stats, err := datastore.GetSiteStats(date)
	if err != nil && err != datastore.ErrNoResults {
		return nil, err
	}

	if stats == nil {
		return &models.SiteStats{
			Date: date,
		}, nil
	}

	return stats, nil
}
