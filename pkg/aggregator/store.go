package aggregator

import (
	"fmt"
	"strings"
	"time"

	"github.com/usefathom/fathom/pkg/datastore"
	"github.com/usefathom/fathom/pkg/models"
)

func (agg *Aggregator) getSiteStats(r *results, siteID int64, t time.Time) (*models.SiteStats, error) {
	cacheKey := fmt.Sprintf("%d-%s", siteID, t.Format("2006-01-02T15"))
	if stats, ok := r.Sites[cacheKey]; ok {
		return stats, nil

	}

	// get from db
	stats, err := agg.database.GetSiteStats(siteID, t)
	if err != nil && err != datastore.ErrNoResults {
		return nil, err
	}

	if stats == nil {
		stats = &models.SiteStats{
			SiteID: siteID,
			New:    true,
			Date:   t,
		}
	}

	r.Sites[cacheKey] = stats
	return stats, nil
}

func (agg *Aggregator) getPageStats(r *results, siteID int64, t time.Time, hostname string, pathname string) (*models.PageStats, error) {
	cacheKey := fmt.Sprintf("%d-%s-%s-%s", siteID, t.Format("2006-01-02T15"), hostname, pathname)
	if stats, ok := r.Pages[cacheKey]; ok {
		return stats, nil
	}

	hostnameID, err := agg.database.HostnameID(hostname)
	if err != nil {
		return nil, err
	}

	pathnameID, err := agg.database.PathnameID(pathname)
	if err != nil {
		return nil, err
	}

	stats, err := agg.database.GetPageStats(siteID, t, hostnameID, pathnameID)
	if err != nil && err != datastore.ErrNoResults {
		return nil, err
	}

	if stats == nil {
		stats = &models.PageStats{
			SiteID:     siteID,
			New:        true,
			HostnameID: hostnameID,
			PathnameID: pathnameID,
			Date:       t,
		}

	}

	r.Pages[cacheKey] = stats
	return stats, nil
}

func (agg *Aggregator) getReferrerStats(r *results, siteID int64, t time.Time, hostname string, pathname string) (*models.ReferrerStats, error) {
	cacheKey := fmt.Sprintf("%d-%s-%s-%s", siteID, t.Format("2006-01-02T15"), hostname, pathname)
	if stats, ok := r.Referrers[cacheKey]; ok {
		return stats, nil
	}

	hostnameID, err := agg.database.HostnameID(hostname)
	if err != nil {
		return nil, err
	}

	pathnameID, err := agg.database.PathnameID(pathname)
	if err != nil {
		return nil, err
	}

	// get from db
	stats, err := agg.database.GetReferrerStats(siteID, t, hostnameID, pathnameID)
	if err != nil && err != datastore.ErrNoResults {
		return nil, err
	}

	if stats == nil {
		stats = &models.ReferrerStats{
			SiteID:     siteID,
			New:        true,
			HostnameID: hostnameID,
			PathnameID: pathnameID,
			Date:       t,
			Group:      "",
		}

		if strings.Contains(hostname, "www.google.") {
			stats.Group = "Google"
		} else if strings.Contains(stats.Hostname, "www.bing.") {
			stats.Group = "Bing"
		} else if strings.Contains(stats.Hostname, "www.baidu.") {
			stats.Group = "Baidu"
		} else if strings.Contains(stats.Hostname, "www.yandex.") {
			stats.Group = "Yandex"
		} else if strings.Contains(stats.Hostname, "search.yahoo.") {
			stats.Group = "Yahoo!"
		} else if strings.Contains(stats.Hostname, "www.findx.") {
			stats.Group = "Findx"
		}
	}

	r.Referrers[cacheKey] = stats
	return stats, nil
}
