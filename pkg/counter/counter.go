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

	for _, p := range pageviews {
		// site stats
		siteStats.Pageviews += 1

		if p.IsNewVisitor {
			siteStats.Visitors += 1

			// TODO: Only new sessions can bounce, not only new visitors. So this is inaccurate right now.
			if p.IsBounce {
				siteStats.Bounced = ((siteStats.BouncedN * siteStats.Bounced) + 10) / (siteStats.BouncedN + 1)
			} else {
				siteStats.Bounced = ((siteStats.BouncedN * siteStats.Bounced) + 0) / (siteStats.BouncedN + 1)
			}
			siteStats.BouncedN += 1
		}

		siteStats.AvgDuration = ((siteStats.AvgDuration * siteStats.AvgDurationN) + p.Duration) / (siteStats.AvgDurationN + 1)
		siteStats.AvgDurationN += 1

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

		pageStats.AvgDuration = ((pageStats.AvgDuration * pageStats.AvgDurationN) + p.Duration) / (pageStats.AvgDurationN + 1)
		pageStats.AvgDurationN += 1

		if p.IsNewVisitor {
			if p.IsBounce {
				pageStats.Bounced = ((pageStats.BouncedN * pageStats.Bounced) + 1) / (pageStats.BouncedN + 1)
			} else {
				pageStats.Bounced = ((pageStats.BouncedN * pageStats.Bounced) + 0) / (pageStats.BouncedN + 1)
			}
			pageStats.BouncedN += 1
		}

		// TODO: referrer stats
	}

	for _, pageStats := range pages {
		err = datastore.UpdatePageStats(pageStats)
		if err != nil {
			log.Error(err)
			return err
		}
	}

	err = datastore.UpdateSiteStats(siteStats)
	if err != nil {
		log.Error(err)
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
