package aggregator

import (
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/usefathom/fathom/pkg/datastore"
	"github.com/usefathom/fathom/pkg/models"

	log "github.com/sirupsen/logrus"
)

type Aggregator struct {
	database datastore.Datastore
}

type Report struct {
	Processed int
	PoolEmpty bool
	Duration  time.Duration
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
func (agg *Aggregator) Run() Report {
	startTime := time.Now()

	// Get unprocessed pageviews
	limit := 10000
	pageviews, err := agg.database.GetProcessablePageviews(limit)
	emptyReport := Report{
		Processed: 0,
	}
	if err != nil && err != datastore.ErrNoResults {
		log.Error(err)
		return emptyReport
	}

	//  Do we have anything to process?
	n := len(pageviews)
	if n == 0 {
		return emptyReport
	}

	results := &results{
		Sites:     map[string]*models.SiteStats{},
		Pages:     map[string]*models.PageStats{},
		Referrers: map[string]*models.ReferrerStats{},
	}

	sites, err := agg.database.GetSites()
	if err != nil {
		log.Error(err)
		return emptyReport
	}

	// create map of public tracking ID's => site ID
	trackingIDMap := make(map[string]int64, len(sites)+1)
	for _, s := range sites {
		trackingIDMap[s.TrackingID] = s.ID
	}

	// if no explicit site ID was given in the tracking request, default to site with ID 1
	trackingIDMap[""] = 1

	// setup referrer spam blacklist
	blacklist, err := newBlacklist()
	if err != nil {
		log.Error(err)
		return emptyReport
	}

	// add each pageview to the various statistics we gather
	for _, p := range pageviews {
		// discard pageview if site tracking ID is unknown
		siteID, ok := trackingIDMap[p.SiteTrackingID]
		if !ok {
			log.Debugf("Skipping pageview because of unrecognized site tracking ID %s", p.SiteTrackingID)
			continue
		}

		// start with referrer because we may want to skip this pageview altogether if it is referrer spam
		if p.Referrer != "" {
			ref, err := parseReferrer(p.Referrer)
			if err != nil {
				log.Debugf("Skipping pageview from referrer %s because of malformed referrer URL", p.Referrer)
				continue
			}

			// ignore out pageviews from blacklisted referrers
			// we use Hostname() here to discard port numbers
			if blacklist.Has(ref.Hostname()) {
				log.Debugf("Skipping pageview from referrer %s because of blacklist", p.Referrer)
				continue
			}

			hostname := ref.Scheme + "://" + ref.Host
			referrerStats, err := agg.getReferrerStats(results, siteID, p.Timestamp, hostname, ref.Path)
			if err != nil {
				log.Error(err)
				continue
			}
			referrerStats.HandlePageview(p)
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

	endTime := time.Now()
	dur := endTime.Sub(startTime)

	report := Report{
		Processed: n,
		PoolEmpty: n < limit,
		Duration:  dur,
	}
	log.Debugf("processed %d pageviews. took: %s, pool empty: %v", report.Processed, report.Duration, report.PoolEmpty)
	return report
}

// parseReferrer parses the referrer string & normalizes it
func parseReferrer(r string) (*url.URL, error) {
	u, err := url.Parse(r)
	if err != nil {
		return nil, err
	}

	// always require a hostname
	if u.Host == "" {
		return nil, errors.New("malformed URL, empty host")
	}

	// remove AMP & UTM vars
	if u.RawQuery != "" {
		q := u.Query()
		keys := []string{"amp", "utm_campaign", "utm_medium", "utm_source"}
		for _, k := range keys {
			q.Del(k)
		}
		u.RawQuery = q.Encode()
	}

	// remove amp/ suffix (but keep trailing slash)
	if strings.HasSuffix(u.Path, "/amp/") {
		u.Path = u.Path[0:(len(u.Path) - 4)]
	}

	// re-parse our normalized string into a new URL struct
	return url.Parse(u.String())
}
