package aggregator

import (
	"github.com/usefathom/fathom/pkg/models"
)

type results struct {
	Sites     map[string]*models.SiteStats
	Pages     map[string]*models.PageStats
	Referrers map[string]*models.ReferrerStats
}

func newResults() *results {
	return &results{
		Sites:     map[string]*models.SiteStats{},
		Pages:     map[string]*models.PageStats{},
		Referrers: map[string]*models.ReferrerStats{},
	}
}
