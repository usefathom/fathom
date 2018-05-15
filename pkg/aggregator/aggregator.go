package aggregator

import (
	"github.com/usefathom/fathom/pkg/datastore"
	"github.com/usefathom/fathom/pkg/models"

	log "github.com/sirupsen/logrus"
)

type aggregator struct {
	database datastore.Datastore
}

func New(db datastore.Datastore) *aggregator {
	return &aggregator{
		database: db,
	}
}

func (agg *aggregator) Run() {
	// Get unprocessed pageviews
	pageviews, err := agg.database.GetProcessablePageviews()
	if err != nil && err != datastore.ErrNoResults {
		log.Error(err)
		return
	}

	//  Do we have anything to process?
	if len(pageviews) == 0 {
		return
	}

	results := agg.Process(pageviews)

	// update stats
	for _, site := range results.Sites {
		err = agg.database.UpdateSiteStats(site)
		if err != nil {
			log.Error(err)
		}
	}

	for _, pageStats := range results.Pages {
		err = agg.database.UpdatePageStats(pageStats)
		if err != nil {
			log.Error(err)
		}
	}

	for _, referrerStats := range results.Referrers {
		err = agg.database.UpdateReferrerStats(referrerStats)
		if err != nil {
			log.Error(err)
		}
	}

	// finally, remove pageviews that we just processed
	err = agg.database.DeletePageviews(pageviews)
	if err != nil {
		log.Error(err)
	}
}

func (agg *aggregator) Process(pageviews []*models.Pageview) *Results {
	log.Debugf("processing %d pageviews", len(pageviews))
	results := NewResults()

	for _, p := range pageviews {
		site, err := agg.getSiteStats(results, p.Timestamp)
		if err != nil {
			log.Error(err)
			continue
		}

		site.Pageviews += 1

		// TODO: Weight isn't right here because we need the number of pageview with a known time of page, not all pageviews
		if p.Duration > 0.00 {
			site.AvgDuration = site.AvgDuration + ((float64(p.Duration) - site.AvgDuration) * 1 / float64(site.Pageviews))
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

		pageStats, err := agg.getPageStats(results, p.Timestamp, p.Hostname, p.Pathname)
		if err != nil {
			log.Error(err)
			continue
		}

		pageStats.Pageviews += 1
		if p.IsUnique {
			pageStats.Visitors += 1
		}

		if p.Duration > 0.00 {
			pageStats.AvgDuration = pageStats.AvgDuration + ((float64(p.Duration) - pageStats.AvgDuration) * 1 / float64(pageStats.Pageviews))
		}

		if p.IsNewSession {
			pageStats.Entries += 1

			if p.IsBounce {
				pageStats.BounceRate = ((float64(pageStats.Entries-1) * pageStats.BounceRate) + 1.00) / (float64(pageStats.Entries))
			} else {
				pageStats.BounceRate = ((float64(pageStats.Entries-1) * pageStats.BounceRate) + 0.00) / (float64(pageStats.Entries))
			}
		}

		// referrer stats
		if p.Referrer != "" {
			referrerStats, err := agg.getReferrerStats(results, p.Timestamp, p.Referrer)
			if err != nil {
				log.Error(err)
				continue
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
				referrerStats.AvgDuration = referrerStats.AvgDuration + ((float64(p.Duration) - referrerStats.AvgDuration) * 1 / float64(referrerStats.Pageviews))
			}

		}

	}

	return results
}
