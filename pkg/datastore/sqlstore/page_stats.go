package sqlstore

import (
	"time"

	"github.com/usefathom/fathom/pkg/models"
)

func (db *sqlstore) GetPageStats(siteID int64, date time.Time, hostname string, pathname string) (*models.PageStats, error) {
	stats := &models.PageStats{}
	query := db.Rebind(`SELECT * FROM daily_page_stats WHERE site_id = ? AND hostname = ? AND pathname = ? AND date = ? LIMIT 1`)
	err := db.Get(stats, query, siteID, hostname, pathname, date.Format("2006-01-02"))
	return stats, mapError(err)
}

func (db *sqlstore) SavePageStats(s *models.PageStats) error {
	if s.New {
		return db.insertPageStats(s)
	}

	return db.updatePageStats(s)
}

func (db *sqlstore) insertPageStats(s *models.PageStats) error {
	query := db.Rebind(`INSERT INTO daily_page_stats(pageviews, visitors, entries, bounce_rate, avg_duration, known_durations, site_id, hostname, pathname, date) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	_, err := db.Exec(query, s.Pageviews, s.Visitors, s.Entries, s.BounceRate, s.AvgDuration, s.KnownDurations, s.SiteID, s.Hostname, s.Pathname, s.Date.Format("2006-01-02"))
	return err
}

func (db *sqlstore) updatePageStats(s *models.PageStats) error {
	query := db.Rebind(`UPDATE daily_page_stats SET pageviews = ?, visitors = ?, entries = ?, bounce_rate = ?, avg_duration = ?, known_durations = ? WHERE site_id = ? AND hostname = ? AND pathname = ? AND date = ?`)
	_, err := db.Exec(query, s.Pageviews, s.Visitors, s.Entries, s.BounceRate, s.AvgDuration, s.KnownDurations, s.SiteID, s.Hostname, s.Pathname, s.Date.Format("2006-01-02"))
	return err
}

func (db *sqlstore) GetAggregatedPageStats(siteID int64, startDate time.Time, endDate time.Time, limit int64) ([]*models.PageStats, error) {
	var result []*models.PageStats
	query := db.Rebind(`SELECT 
		hostname, 
		pathname, 
		SUM(pageviews) AS pageviews, 
		SUM(visitors) AS visitors, 
		SUM(entries) AS entries, 
		COALESCE(SUM(entries*bounce_rate) / NULLIF(SUM(entries), 0), 0.00) AS bounce_rate, 
		COALESCE(SUM(avg_duration*pageviews) / SUM(pageviews), 0.00) AS avg_duration 
		FROM daily_page_stats WHERE site_id = ? AND date >= ? AND date <= ? 
		GROUP BY hostname, pathname 
		ORDER BY pageviews DESC LIMIT ?`)
	err := db.Select(&result, query, siteID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), limit)
	return result, err
}

func (db *sqlstore) GetAggregatedPageStatsPageviews(siteID int64, startDate time.Time, endDate time.Time) (int64, error) {
	var result int64
	query := db.Rebind(`SELECT COALESCE(SUM(pageviews), 0) FROM daily_page_stats WHERE site_id = ? AND date >= ? AND date <= ?`)
	err := db.Get(&result, query, siteID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	return result, err
}
