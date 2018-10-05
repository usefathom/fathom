package aggregator

import (
	"fmt"
	"strings"
	"time"

	"github.com/usefathom/fathom/pkg/datastore"
	"github.com/usefathom/fathom/pkg/models"
)

func (agg *aggregator) getSiteStats(r *results, siteID int64, t time.Time) (*models.SiteStats, error) {
	// get from map
	cacheKey := fmt.Sprintf("%d-%s", siteID, t.Format("2006-01-02"))
	if stats, ok := r.Sites[cacheKey]; ok {
		return stats, nil

	}

	// get from db
	stats, err := agg.database.GetSiteStats(siteID, t)
	if err != nil && err != datastore.ErrNoResults {
		return nil, err
	}

	// create in db
	if stats == nil {
		stats = &models.SiteStats{
			SiteID: siteID,
			Date:   t,
			New:    true,
		}
	}

	r.Sites[cacheKey] = stats
	return stats, nil
}

func (agg *aggregator) getPageStats(r *results, siteID int64, t time.Time, hostname string, pathname string) (*models.PageStats, error) {
	cacheKey := fmt.Sprintf("%d-%s-%s-%s", siteID, t.Format("2006-01-02"), hostname, pathname)
	if stats, ok := r.Pages[cacheKey]; ok {
		return stats, nil
	}

	stats, err := agg.database.GetPageStats(siteID, t, hostname, pathname)
	if err != nil && err != datastore.ErrNoResults {
		return nil, err
	}

	if stats == nil {
		stats = &models.PageStats{
			SiteID:   siteID,
			New:      true,
			Hostname: hostname,
			Pathname: pathname,
			Date:     t,
		}

	}

	r.Pages[cacheKey] = stats
	return stats, nil
}

func (agg *aggregator) getReferrerStats(r *results, siteID int64, t time.Time, hostname string, pathname string) (*models.ReferrerStats, error) {
	cacheKey := fmt.Sprintf("%d-%s-%s-%s", siteID, t.Format("2006-01-02"), hostname, pathname)
	if stats, ok := r.Referrers[cacheKey]; ok {
		return stats, nil
	}

	// get from db
	stats, err := agg.database.GetReferrerStats(siteID, t, hostname, pathname)
	if err != nil && err != datastore.ErrNoResults {
		return nil, err
	}

	// create in db
	if stats == nil {
		stats = &models.ReferrerStats{
			SiteID:   siteID,
			New:      true,
			Hostname: hostname,
			Pathname: pathname,
			Date:     t,
			Group:    "",
		}

		// TODO: Abstract this
		if strings.Contains(stats.Hostname, "www.google.") {
			stats.Group = "Google"
		}
	}

	r.Referrers[cacheKey] = stats
	return stats, nil
}
