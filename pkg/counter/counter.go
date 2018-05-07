package counter

import (
	"time"

	"github.com/usefathom/fathom/pkg/datastore"
	"github.com/usefathom/fathom/pkg/models"

	log "github.com/sirupsen/logrus"
)

func Aggregate() error {
	now := time.Now()

	// Get unprocessed pageviews
	pageviews, err := datastore.GetProcessablePageviews()
	if err != nil && err != datastore.ErrNoResults {
		return err
	}

	//  Do we have anything to process?
	if len(pageviews) == 0 {
		return nil
	}

	// site stats
	siteStats, err := getSiteStats(now)
	if err != nil {
		return err
	}

	pages := map[string]*models.PageStats{}
	referrers := map[string]*models.ReferrerStats{}

	for _, p := range pageviews {
		// site stats
		siteStats.Pageviews += 1

		if p.IsNewVisitor {
			siteStats.Visitors += 1

			// TODO: Only new sessions can bounce, not only new visitors. So this is inaccurate right now.
			if p.IsBounce {
				siteStats.Bounces = ((siteStats.Sessions * siteStats.Bounces) + 1) / (siteStats.Sessions + 1)
			} else {
				siteStats.Bounces = ((siteStats.Sessions * siteStats.Bounces) + 0) / (siteStats.Sessions + 1)
			}
			siteStats.Sessions += 1
		}

		siteStats.AvgDuration = ((siteStats.AvgDuration * (siteStats.Pageviews - 1)) + p.Duration) / siteStats.Pageviews

		// page stats
		var pageStats *models.PageStats
		var ok bool
		if pageStats, ok = pages[p.Pathname]; !ok {
			pageStats, err = getPageStats(now, p.Pathname)
			if err != nil {
				log.Error(err)
				continue
			}
			pages[p.Pathname] = pageStats
		}

		pageStats.Views += 1
		if p.IsUnique {
			pageStats.UniqueViews += 1
		}

		pageStats.AvgDuration = (pageStats.AvgDuration*(pageStats.Views-1) + p.Duration) / pageStats.Views

		if p.IsNewVisitor {
			if p.IsBounce {
				pageStats.Bounces = ((pageStats.Entries * pageStats.Bounces) + 1) / (pageStats.Entries + 1)
			} else {
				pageStats.Bounces = ((pageStats.Entries * pageStats.Bounces) + 0) / (pageStats.Entries + 1)
			}
			pageStats.Entries += 1
		}

		if p.Referrer != "" {
			var referrerStats *models.ReferrerStats
			var ok bool
			if referrerStats, ok = referrers[p.Referrer]; !ok {
				referrerStats, err = getReferrerStats(now, p.Referrer)
				if err != nil {
					log.Error(err)
					continue
				}
				referrers[p.Referrer] = referrerStats
			}

			referrerStats.Pageviews += 1

			if p.IsNewVisitor {
				referrerStats.Visitors += 1
			}
		}

	}

	// update stats
	err = datastore.UpdateSiteStats(siteStats)
	if err != nil {
		log.Error(err)
		return err
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
