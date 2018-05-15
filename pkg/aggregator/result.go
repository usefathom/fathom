package aggregator

import (
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
