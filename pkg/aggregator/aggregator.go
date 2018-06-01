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
		err := agg.handleSiteview(results, p)
		if err != nil {
			continue
		}

		err = agg.handlePageview(results, p)
		if err != nil {
			continue
		}

		// referrer stats
		if p.Referrer != "" {
			err := agg.handleReferral(results, p)
			if err != nil {
				continue
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
