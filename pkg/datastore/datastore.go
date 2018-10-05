package datastore

import (
	"time"

	"github.com/usefathom/fathom/pkg/datastore/sqlstore"
	"github.com/usefathom/fathom/pkg/models"
)

// ErrNoResults is returned whenever a single-item query returns 0 results
var ErrNoResults = sqlstore.ErrNoResults // ???

// Datastore represents a database implementations
type Datastore interface {
	// users
	GetUser(int64) (*models.User, error)
	GetUserByEmail(string) (*models.User, error)
	SaveUser(*models.User) error
	DeleteUser(*models.User) error
	CountUsers() (int64, error)

	// sites
	GetSites() ([]*models.Site, error)
	SaveSite(s *models.Site) error
	DeleteSite(s *models.Site) error

	// site stats
	GetSiteStats(int64, time.Time) (*models.SiteStats, error)
	GetSiteStatsPerDay(int64, time.Time, time.Time) ([]*models.SiteStats, error)
	InsertSiteStats(*models.SiteStats) error
	UpdateSiteStats(*models.SiteStats) error
	GetAggregatedSiteStats(int64, time.Time, time.Time) (*models.SiteStats, error)
	GetTotalSiteViews(int64, time.Time, time.Time) (int64, error)
	GetTotalSiteVisitors(int64, time.Time, time.Time) (int64, error)
	GetTotalSiteSessions(int64, time.Time, time.Time) (int64, error)
	GetAverageSiteDuration(int64, time.Time, time.Time) (float64, error)
	GetAverageSiteBounceRate(int64, time.Time, time.Time) (float64, error)
	GetRealtimeVisitorCount(int64) (int64, error)

	// pageviews
	InsertPageviews([]*models.Pageview) error
	UpdatePageviews([]*models.Pageview) error
	GetPageview(string) (*models.Pageview, error)
	GetProcessablePageviews() ([]*models.Pageview, error)
	DeletePageviews([]*models.Pageview) error

	// page stats
	GetPageStats(int64, time.Time, string, string) (*models.PageStats, error)
	InsertPageStats(*models.PageStats) error
	UpdatePageStats(*models.PageStats) error
	GetAggregatedPageStats(int64, time.Time, time.Time, int64) ([]*models.PageStats, error)
	GetAggregatedPageStatsPageviews(int64, time.Time, time.Time) (int64, error)

	// referrer stats
	GetReferrerStats(int64, time.Time, string, string) (*models.ReferrerStats, error)
	InsertReferrerStats(*models.ReferrerStats) error
	UpdateReferrerStats(*models.ReferrerStats) error
	GetAggregatedReferrerStats(int64, time.Time, time.Time, int64) ([]*models.ReferrerStats, error)
	GetAggregatedReferrerStatsPageviews(int64, time.Time, time.Time) (int64, error)

	// misc
	Health() error
	Close() error
}

// New instantiates a new datastore from the given configuration struct
func New(c *sqlstore.Config) Datastore {
	return sqlstore.New(c)
}
