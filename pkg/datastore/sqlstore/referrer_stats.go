package sqlstore

import (
	"time"

	"github.com/usefathom/fathom/pkg/models"
)

func (db *sqlstore) GetReferrerStats(siteID int64, date time.Time, hostname string, pathname string) (*models.ReferrerStats, error) {
	stats := &models.ReferrerStats{}
	query := db.Rebind(`SELECT * FROM daily_referrer_stats WHERE site_id = ? AND date = ? AND hostname = ? AND pathname = ? LIMIT 1`)
	err := db.Get(stats, query, siteID, date.Format("2006-01-02"), hostname, pathname)
	return stats, mapError(err)
}

func (db *sqlstore) SaveReferrerStats(s *models.ReferrerStats) error {
	if s.New {
		return db.insertReferrerStats(s)
	}

	return db.updateReferrerStats(s)
}

func (db *sqlstore) insertReferrerStats(s *models.ReferrerStats) error {
	query := db.Rebind(`INSERT INTO daily_referrer_stats(visitors, pageviews, bounce_rate, avg_duration, known_durations, groupname, site_id, hostname, pathname, date) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	_, err := db.Exec(query, s.Visitors, s.Pageviews, s.BounceRate, s.AvgDuration, s.KnownDurations, s.Group, s.SiteID, s.Hostname, s.Pathname, s.Date.Format("2006-01-02"))
	return err
}

func (db *sqlstore) updateReferrerStats(s *models.ReferrerStats) error {
	query := db.Rebind(`UPDATE daily_referrer_stats SET visitors = ?, pageviews = ?, bounce_rate = ?, avg_duration = ?, known_durations = ?, groupname = ? WHERE site_id = ? AND hostname = ? AND pathname = ? AND date = ?`)
	_, err := db.Exec(query, s.Visitors, s.Pageviews, s.BounceRate, s.AvgDuration, s.KnownDurations, s.Group, s.SiteID, s.Hostname, s.Pathname, s.Date.Format("2006-01-02"))
	return err
}

func (db *sqlstore) GetAggregatedReferrerStats(siteID int64, startDate time.Time, endDate time.Time, limit int64) ([]*models.ReferrerStats, error) {
	var result []*models.ReferrerStats

	sql := `SELECT 
		MIN(hostname) AS hostname,
		MIN(pathname) AS pathname,
		COALESCE(MIN(groupname), '') AS groupname,  
		SUM(visitors) AS visitors, 
		SUM(pageviews) AS pageviews, 
		COALESCE(SUM(pageviews*NULLIF(bounce_rate, 0)) / SUM(pageviews), 0.00) AS bounce_rate, 
		COALESCE(SUM(avg_duration*pageviews) / SUM(pageviews), 0.00) AS avg_duration 
	FROM daily_referrer_stats 
	WHERE site_id = ? AND date >= ? AND date <= ? `

	if db.Config.Driver == "sqlite3" {
		sql = sql + `GROUP BY COALESCE(NULLIF(groupname, ''), hostname || pathname ) `
	} else {
		sql = sql + `GROUP BY COALESCE(NULLIF(groupname, ''), CONCAT(hostname, pathname) ) `
	}
	sql = sql + ` ORDER BY pageviews DESC LIMIT ?`

	query := db.Rebind(sql)

	err := db.Select(&result, query, siteID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), limit)
	return result, mapError(err)
}

func (db *sqlstore) GetAggregatedReferrerStatsPageviews(siteID int64, startDate time.Time, endDate time.Time) (int64, error) {
	var result int64
	query := db.Rebind(`SELECT COALESCE(SUM(pageviews), 0) FROM daily_referrer_stats WHERE site_id = ? AND date >= ? AND date <= ?`)
	err := db.Get(&result, query, siteID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	return result, mapError(err)
}
