package sqlstore

import (
	"database/sql"
	"time"

	"github.com/usefathom/fathom/pkg/models"
)

func (db *sqlstore) GetReferrerStats(siteID int64, date time.Time, hostnameID int64, pathnameID int64) (*models.ReferrerStats, error) {
	stats := &models.ReferrerStats{New: false}
	query := db.Rebind(`SELECT * FROM referrer_stats WHERE site_id = ? AND ts = ? AND hostname_id = ? AND pathname_id = ? LIMIT 1`)
	err := db.Get(stats, query, siteID, date.Format(DATE_FORMAT), hostnameID, pathnameID)
	if err == sql.ErrNoRows {
		return nil, ErrNoResults
	}

	return stats, mapError(err)
}

func (db *sqlstore) SaveReferrerStats(s *models.ReferrerStats) error {
	if s.New {
		return db.insertReferrerStats(s)
	}

	return db.updateReferrerStats(s)
}

func (db *sqlstore) insertReferrerStats(s *models.ReferrerStats) error {
	query := db.Rebind(`INSERT INTO referrer_stats(visitors, pageviews, bounce_rate, avg_duration, known_durations, groupname, site_id, hostname_id, pathname_id, ts) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	_, err := db.Exec(query, s.Visitors, s.Pageviews, s.BounceRate, s.AvgDuration, s.KnownDurations, s.Group, s.SiteID, s.HostnameID, s.PathnameID, s.Date.Format(DATE_FORMAT))
	return err
}

func (db *sqlstore) updateReferrerStats(s *models.ReferrerStats) error {
	query := db.Rebind(`UPDATE referrer_stats SET visitors = ?, pageviews = ?, bounce_rate = ?, avg_duration = ?, known_durations = ?, groupname = ? WHERE site_id = ? AND hostname_id = ? AND pathname_id = ? AND ts = ?`)
	_, err := db.Exec(query, s.Visitors, s.Pageviews, s.BounceRate, s.AvgDuration, s.KnownDurations, s.Group, s.SiteID, s.HostnameID, s.PathnameID, s.Date.Format(DATE_FORMAT))
	return err
}

func (db *sqlstore) SelectAggregatedReferrerStats(siteID int64, startDate time.Time, endDate time.Time, offset int64, limit int64) ([]*models.ReferrerStats, error) {
	var result []*models.ReferrerStats

	sql := `SELECT 
		MIN(h.name) AS hostname,
		MIN(p.name) AS pathname,
		COALESCE(MIN(groupname), '') AS groupname,  
		SUM(visitors) AS visitors, 
		SUM(pageviews) AS pageviews, 
		SUM(pageviews*bounce_rate) / SUM(pageviews) AS bounce_rate, 
		SUM(pageviews*avg_duration) / SUM(pageviews) AS avg_duration 
	FROM referrer_stats s
		LEFT JOIN hostnames h ON h.id = s.hostname_id 
		LEFT JOIN pathnames p ON p.id = s.pathname_id 
	WHERE site_id = ? AND ts >= ? AND ts <= ? `

	if db.Config.Driver == "sqlite3" {
		sql = sql + `GROUP BY COALESCE(NULLIF(groupname, ''), hostname_id || pathname_id ) `
	} else {
		sql = sql + `GROUP BY COALESCE(NULLIF(groupname, ''), CONCAT(hostname_id, pathname_id) ) `
	}
	sql = sql + ` ORDER BY pageviews DESC LIMIT ?, ?`

	query := db.Rebind(sql)

	err := db.Select(&result, query, siteID, startDate.Format(DATE_FORMAT), endDate.Format(DATE_FORMAT), offset, limit)
	return result, mapError(err)
}

func (db *sqlstore) GetAggregatedReferrerStatsPageviews(siteID int64, startDate time.Time, endDate time.Time) (int64, error) {
	var result int64
	query := db.Rebind(`SELECT COALESCE(SUM(pageviews), 0) FROM referrer_stats WHERE site_id = ? AND ts >= ? AND ts <= ?`)
	err := db.Get(&result, query, siteID, startDate.Format(DATE_FORMAT), endDate.Format(DATE_FORMAT))
	return result, mapError(err)
}
