package aggregator

import (
	"github.com/usefathom/fathom/pkg/datastore"
	"github.com/usefathom/fathom/pkg/models"
	"net/url"

	log "github.com/sirupsen/logrus"
)

type aggregator struct {
	database datastore.Datastore
}

// New returns a new aggregator instance with the database dependency injected.
func New(db datastore.Datastore) *aggregator {
	return &aggregator{
		database: db,
	}
}

// Run processes the pageviews which are ready to be processed and adds them to daily aggregation
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

// Process processes the given pageviews and returns the (aggregated) results per metric per day
func (agg *aggregator) Process(pageviews []*models.Pageview) *results {
	log.Debugf("processing %d pageviews", len(pageviews))
	results := newResults()

	for _, p := range pageviews {
		site, err := agg.getSiteStats(results, p.Timestamp)
		if err != nil {
			log.Error(err)
			continue
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

		// referrer stats
		if p.Referrer != "" {
			hostname, pathname, _ := parseUrlParts(p.Referrer)
			referrerStats, err := agg.getReferrerStats(results, p.Timestamp, hostname, pathname)
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
				referrerStats.KnownDurations += 1
				referrerStats.AvgDuration = referrerStats.AvgDuration + ((float64(p.Duration) - referrerStats.AvgDuration) * 1 / float64(referrerStats.KnownDurations))
			}

		}

	}

	return results
}

func parseUrlParts(s string) (string, string, error) {
	u, err := url.Parse(s)
	if err != nil {
		return "", "", err
	}

	return u.Scheme + "://" + u.Host, u.Path, nil
}
