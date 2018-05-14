package aggregator

import (
	"time"

	"github.com/usefathom/fathom/pkg/models"
)

type Results struct {
	Sites     map[string]*models.SiteStats
	Pages     map[string]*models.PageStats
	Referrers map[string]*models.ReferrerStats
}

func NewResults() *Results {
	return &Results{
		Sites:     map[string]*models.SiteStats{},
		Pages:     map[string]*models.PageStats{},
		Referrers: map[string]*models.ReferrerStats{},
	}
}

func (r *Results) GetSiteStats(t time.Time) (*models.SiteStats, error) {
	var stats *models.SiteStats
	var ok bool
	var err error

	date := t.Format("2006-01-02")
	if stats, ok = r.Sites[date]; !ok {
		stats, err = getSiteStats(t)
		if err != nil {
			return nil, err
		}
		r.Sites[date] = stats
	}
	return stats, nil
}

func (r *Results) GetPageStats(t time.Time, hostname string, pathname string) (*models.PageStats, error) {
	var stats *models.PageStats
	var ok bool
	var err error

	date := t.Format("2006-01-02")
	if stats, ok = r.Pages[date+hostname+pathname]; !ok {
		stats, err = getPageStats(t, hostname, pathname)
		if err != nil {
			return nil, err
		}
		r.Pages[date+hostname+pathname] = stats
	}
	return stats, nil
}

func (r *Results) GetReferrerStats(t time.Time, referrer string) (*models.ReferrerStats, error) {
	var stats *models.ReferrerStats
	var ok bool
	var err error

	date := t.Format("2006-01-02")
	if stats, ok = r.Referrers[date+referrer]; !ok {
		stats, err = getReferrerStats(t, referrer)
		if err != nil {
			return nil, err
		}
		r.Referrers[date+referrer] = stats
	}

	return stats, nil
}
