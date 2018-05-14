package aggregator

import (
	"time"

	"github.com/usefathom/fathom/pkg/datastore"
	"github.com/usefathom/fathom/pkg/models"

	log "github.com/sirupsen/logrus"
)

func Run() {
	// Get unprocessed pageviews
	pageviews, err := datastore.GetProcessablePageviews()
	if err != nil && err != datastore.ErrNoResults {
		log.Error(err)
		return
	}

	//  Do we have anything to process?
	if len(pageviews) == 0 {
		return
	}

	results := Process(pageviews)

	// update stats
	for _, site := range results.Sites {
		err = datastore.UpdateSiteStats(site)
		if err != nil {
			log.Error(err)
		}
	}

	for _, pageStats := range results.Pages {
		err = datastore.UpdatePageStats(pageStats)
		if err != nil {
			log.Error(err)
		}
	}

	for _, referrerStats := range results.Referrers {
		err = datastore.UpdateReferrerStats(referrerStats)
		if err != nil {
			log.Error(err)
		}
	}

	// finally, remove pageviews that we just processed
	err = datastore.DeletePageviews(pageviews)
	if err != nil {
		log.Error(err)
	}
}

func Process(pageviews []*models.Pageview) *Results {
	log.Debugf("processing %d pageviews", len(pageviews))
	results := NewResults()

	for _, p := range pageviews {
		site, err := results.GetSiteStats(p.Timestamp)
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

		pageStats, err := results.GetPageStats(p.Timestamp, p.Hostname, p.Pathname)
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
			referrerStats, err := results.GetReferrerStats(p.Timestamp, p.Referrer)
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

func getSiteStats(t time.Time) (*models.SiteStats, error) {
	stats, err := datastore.GetSiteStats(t)
	if err != nil && err != datastore.ErrNoResults {
		return nil, err
	}

	if stats != nil {
		return stats, nil
	}

	stats = &models.SiteStats{
		Date: t,
	}
	err = datastore.InsertSiteStats(stats)
	return stats, err
}

func getPageStats(date time.Time, hostname string, pathname string) (*models.PageStats, error) {
	stats, err := datastore.GetPageStats(date, hostname, pathname)
	if err != nil && err != datastore.ErrNoResults {
		return nil, err
	}

	if stats != nil {
		return stats, nil
	}

	stats = &models.PageStats{
		Hostname: hostname,
		Pathname: pathname,
		Date:     date,
	}
	err = datastore.InsertPageStats(stats)
	return stats, err
}

func getReferrerStats(date time.Time, url string) (*models.ReferrerStats, error) {
	stats, err := datastore.GetReferrerStats(date, url)
	if err != nil && err != datastore.ErrNoResults {
		return nil, err
	}

	if stats != nil {
		return stats, nil
	}

	stats = &models.ReferrerStats{
		URL:  url,
		Date: date,
	}
	err = datastore.InsertReferrerStats(stats)
	return stats, err
}
