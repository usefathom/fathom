package aggregator

import (
	"net/url"

	"github.com/usefathom/fathom/pkg/datastore"
	"github.com/usefathom/fathom/pkg/models"

	log "github.com/sirupsen/logrus"
)

type Aggregator struct {
	database datastore.Datastore
}

type results struct {
	Sites     map[string]*models.SiteStats
	Pages     map[string]*models.PageStats
	Referrers map[string]*models.ReferrerStats
}

// New returns a new aggregator instance with the database dependency injected.
func New(db datastore.Datastore) *Aggregator {
	return &Aggregator{
		database: db,
	}
}

// Run processes the pageviews which are ready to be processed and adds them to daily aggregation
func (agg *Aggregator) Run() int {
	// Get unprocessed pageviews
	pageviews, err := agg.database.GetProcessablePageviews()
	if err != nil && err != datastore.ErrNoResults {
		log.Error(err)
		return 0
	}

	//  Do we have anything to process?
	n := len(pageviews)
	if n == 0 {
		return 0
	}

	results := &results{
		Sites:     map[string]*models.SiteStats{},
		Pages:     map[string]*models.PageStats{},
		Referrers: map[string]*models.ReferrerStats{},
	}

	log.Debugf("processing %d pageviews", len(pageviews))

	sites, err := agg.database.GetSites()
	if err != nil {
		log.Error(err)
		return 0
	}

	// create map of public tracking ID's => site ID
	trackingIDMap := make(map[string]int64, len(sites)+1)
	for _, s := range sites {
		trackingIDMap[s.TrackingID] = s.ID
	}

	// if no explicit site ID was given in the tracking request, default to site with ID 1
	trackingIDMap[""] = 1

	// add each pageview to the various statistics we gather
	for _, p := range pageviews {

		// discard pageview if site tracking ID is unknown
		siteID, ok := trackingIDMap[p.SiteTrackingID]
		if !ok {
			log.Debugf("discarding pageview because of unrecognized site tracking ID %s", p.SiteTrackingID)
			continue
		}

		// get existing site stats so we can add this pageview to it
		site, err := agg.getSiteStats(results, siteID, p.Timestamp)
		if err != nil {
			log.Error(err)
			continue
		}
		site.HandlePageview(p)

		pageStats, err := agg.getPageStats(results, siteID, p.Timestamp, p.Hostname, p.Pathname)
		if err != nil {
			log.Error(err)
			continue
		}
		pageStats.HandlePageview(p)

		// referrer stats
		if p.Referrer != "" {
			hostname, pathname, err := parseUrlParts(p.Referrer)
			if err != nil {
				log.Error(err)
				continue
			}

			referrerStats, err := agg.getReferrerStats(results, siteID, p.Timestamp, hostname, pathname)
			if err != nil {
				log.Error(err)
				continue
			}
			referrerStats.HandlePageview(p)
		}

	}

	// update stats
	for _, site := range results.Sites {
		if err := agg.database.SaveSiteStats(site); err != nil {
			log.Error(err)
		}
	}

	for _, pageStats := range results.Pages {
		if err := agg.database.SavePageStats(pageStats); err != nil {
			log.Error(err)
		}
	}

	for _, referrerStats := range results.Referrers {
		if err := agg.database.SaveReferrerStats(referrerStats); err != nil {
			log.Error(err)
		}
	}

	// finally, remove pageviews that we just processed
	if err := agg.database.DeletePageviews(pageviews); err != nil {
		log.Error(err)
	}

	return n
}

func parseUrlParts(s string) (string, string, error) {
	u, err := url.Parse(s)
	if err != nil {
		return "", "", err
	}

	return u.Scheme + "://" + u.Host, u.Path, nil
}
