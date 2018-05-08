package counter

import (
	"time"

	"github.com/usefathom/fathom/pkg/datastore"
	"github.com/usefathom/fathom/pkg/models"

	log "github.com/sirupsen/logrus"
)

func Aggregate() error {
	// Get unprocessed pageviews
	pageviews, err := datastore.GetProcessablePageviews()
	if err != nil && err != datastore.ErrNoResults {
		return err
	}

	//  Do we have anything to process?
	if len(pageviews) == 0 {
		return nil
	}

	sites := map[string]*models.SiteStats{}
	pages := map[string]*models.PageStats{}
	referrers := map[string]*models.ReferrerStats{}

	for _, p := range pageviews {
		date := p.Timestamp.Format("2006-01-02")

		var site *models.SiteStats
		if site, ok := sites[date]; !ok {
			site, _ = getSiteStats(p.Timestamp)
			sites[date] = site
		}

		// site stats
		site.Pageviews += 1
		site.AvgDuration = ((site.AvgDuration * (site.Pageviews - 1)) + p.Duration) / site.Pageviews

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

		// page stats
		var pageStats *models.PageStats
		var ok bool
		if pageStats, ok = pages[date+p.Pathname]; !ok {
			pageStats, err = getPageStats(p.Timestamp, p.Pathname)
			if err != nil {
				log.Error(err)
				continue
			}
			pages[date+p.Pathname] = pageStats
		}

		pageStats.Pageviews += 1
		if p.IsUnique {
			pageStats.Visitors += 1
		}

		pageStats.AvgDuration = (pageStats.AvgDuration*(pageStats.Pageviews-1) + p.Duration) / pageStats.Pageviews

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
			var referrerStats *models.ReferrerStats
			var ok bool
			if referrerStats, ok = referrers[date+p.Referrer]; !ok {
				referrerStats, err = getReferrerStats(p.Timestamp, p.Referrer)
				if err != nil {
					log.Error(err)
					continue
				}
				referrers[date+p.Referrer] = referrerStats
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
		}

	}

	// update stats
	for _, site := range sites {
		err = datastore.UpdateSiteStats(site)
		if err != nil {
			log.Error(err)
		}
	}

	for _, pageStats := range pages {
		err = datastore.UpdatePageStats(pageStats)
		if err != nil {
			log.Error(err)
			return err
		}
	}

	for _, referrerStats := range referrers {
		err = datastore.UpdateReferrerStats(referrerStats)
		if err != nil {
			log.Error(err)
			return err
		}
	}

	// finally, remove pageviews that we just processed
	err = datastore.DeletePageviews(pageviews)
	if err != nil {
		log.Error(err)
	}

	return nil
}

func getSiteStats(date time.Time) (*models.SiteStats, error) {
	stats, err := datastore.GetSiteStats(date)
	if err != nil && err != datastore.ErrNoResults {
		return nil, err
	}

	if stats != nil {
		return stats, nil
	}

	stats = &models.SiteStats{
		Date: date,
	}
	err = datastore.InsertSiteStats(stats)
	return stats, err
}

func getPageStats(date time.Time, pathname string) (*models.PageStats, error) {
	stats, err := datastore.GetPageStats(date, pathname)
	if err != nil && err != datastore.ErrNoResults {
		return nil, err
	}

	if stats != nil {
		return stats, nil
	}

	stats = &models.PageStats{
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
