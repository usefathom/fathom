package sqlstore

import (
	"database/sql"
	"time"

	"github.com/usefathom/fathom/pkg/models"
)

func (db *sqlstore) GetPageStats(siteID int64, date time.Time, hostnameID int64, pathnameID int64) (*models.PageStats, error) {
	stats := &models.PageStats{New: false}
	query := db.Rebind(`SELECT * FROM page_stats WHERE site_id = ? AND hostname_id = ? AND pathname_id = ? AND ts = ? LIMIT 1`)
	err := db.Get(stats, query, siteID, hostnameID, pathnameID, date.Format(DATE_FORMAT))
	if err == sql.ErrNoRows {
		return nil, ErrNoResults
	}

	return stats, mapError(err)
}

func (db *sqlstore) SavePageStats(s *models.PageStats) error {
	if s.New {
		return db.insertPageStats(s)
	}

	return db.updatePageStats(s)
}

func (db *sqlstore) insertPageStats(s *models.PageStats) error {
	query := db.Rebind(`INSERT INTO page_stats(pageviews, visitors, entries, bounce_rate, avg_duration, known_durations, site_id, hostname_id, pathname_id, ts) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	_, err := db.Exec(query, s.Pageviews, s.Visitors, s.Entries, s.BounceRate, s.AvgDuration, s.KnownDurations, s.SiteID, s.HostnameID, s.PathnameID, s.Date.Format(DATE_FORMAT))
	return err
}

func (db *sqlstore) updatePageStats(s *models.PageStats) error {
	query := db.Rebind(`UPDATE page_stats SET pageviews = ?, visitors = ?, entries = ?, bounce_rate = ?, avg_duration = ?, known_durations = ? WHERE site_id = ? AND hostname_id = ? AND pathname_id = ? AND ts = ?`)
	_, err := db.Exec(query, s.Pageviews, s.Visitors, s.Entries, s.BounceRate, s.AvgDuration, s.KnownDurations, s.SiteID, s.HostnameID, s.PathnameID, s.Date.Format(DATE_FORMAT))
	return err
}

func (db *sqlstore) GetAggregatedPageStats(siteID int64, startDate time.Time, endDate time.Time, limit int64) ([]*models.PageStats, error) {
	var result []*models.PageStats
	query := db.Rebind(`SELECT 
		h.name AS hostname,
		p.name AS pathname,
		MAX(SUM(pageviews), 1) AS pageviews, 
		MAX(SUM(visitors), 1) AS visitors, 
		SUM(entries) AS entries, 
		COALESCE(SUM(entries*bounce_rate) / SUM(entries), 0.00) AS bounce_rate, 
		SUM(pageviews*avg_duration) / SUM(pageviews) AS avg_duration 
		FROM page_stats s 
			LEFT JOIN hostnames h ON h.id = s.hostname_id 
			LEFT JOIN pathnames p ON p.id = s.pathname_id 
		WHERE site_id = ? AND ts >= ? AND ts <= ? 
		GROUP BY hostname, pathname 
		ORDER BY pageviews DESC LIMIT ?`)
	err := db.Select(&result, query, siteID, startDate.Format(DATE_FORMAT), endDate.Format(DATE_FORMAT), limit)
	return result, err
}

func (db *sqlstore) GetAggregatedPageStatsPageviews(siteID int64, startDate time.Time, endDate time.Time) (int64, error) {
	var result int64
	query := db.Rebind(`SELECT COALESCE(SUM(pageviews), 0) FROM page_stats WHERE site_id = ? AND ts >= ? AND ts <= ?`)
	err := db.Get(&result, query, siteID, startDate.Format(DATE_FORMAT), endDate.Format(DATE_FORMAT))
	return result, err
}
