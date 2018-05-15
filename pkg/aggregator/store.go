package aggregator

import (
	"time"

	"github.com/usefathom/fathom/pkg/datastore"
	"github.com/usefathom/fathom/pkg/models"
)

func (agg *aggregator) getSiteStats(r *Results, t time.Time) (*models.SiteStats, error) {
	// get from map
	date := t.Format("2006-01-02")
	if stats, ok := r.Sites[date]; ok {
		return stats, nil

	}

	// get from db
	stats, err := agg.database.GetSiteStats(t)
	if err != nil && err != datastore.ErrNoResults {
		return nil, err
	}

	// create in db
	if stats == nil {
		stats = &models.SiteStats{
			Date: t,
		}

		err = agg.database.InsertSiteStats(stats)
		if err != nil {
			return nil, err
		}
	}

	r.Sites[date] = stats
	return stats, nil
}

func (agg *aggregator) getPageStats(r *Results, t time.Time, hostname string, pathname string) (*models.PageStats, error) {
	date := t.Format("2006-01-02")
	if stats, ok := r.Pages[date+hostname+pathname]; ok {
		return stats, nil
	}

	stats, err := agg.database.GetPageStats(t, hostname, pathname)
	if err != nil && err != datastore.ErrNoResults {
		return nil, err
	}

	if stats == nil {
		stats = &models.PageStats{
			Hostname: hostname,
			Pathname: pathname,
			Date:     t,
		}
		err = agg.database.InsertPageStats(stats)
		if err != nil {
			return nil, err
		}
	}

	r.Pages[date+hostname+pathname] = stats
	return stats, nil
}

func (agg *aggregator) getReferrerStats(r *Results, t time.Time, url string) (*models.ReferrerStats, error) {
	date := t.Format("2006-01-02")
	if stats, ok := r.Referrers[date+url]; ok {
		return stats, nil
	}

	// get from db
	stats, err := agg.database.GetReferrerStats(t, url)
	if err != nil && err != datastore.ErrNoResults {
		return nil, err
	}

	// create in db
	if stats == nil {
		stats = &models.ReferrerStats{
			URL:  url,
			Date: t,
		}
		err = agg.database.InsertReferrerStats(stats)
		if err != nil {
			return nil, err
		}
	}

	r.Referrers[date+url] = stats
	return stats, nil
}
