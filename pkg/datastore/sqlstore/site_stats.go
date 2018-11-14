package sqlstore

import (
	"database/sql"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/usefathom/fathom/pkg/models"
)

func (db *sqlstore) GetSiteStats(siteID int64, date time.Time) (*models.SiteStats, error) {
	stats := &models.SiteStats{New: false}
	query := db.Rebind(`SELECT * FROM site_stats WHERE site_id = ? AND ts = ? LIMIT 1`)

	err := db.Get(stats, query, siteID, date.Format(DATE_FORMAT))
	if err == sql.ErrNoRows {
		return nil, ErrNoResults
	}

	return stats, mapError(err)
}

func (db *sqlstore) SaveSiteStats(s *models.SiteStats) error {
	if s.New {
		return db.insertSiteStats(s)
	}

	return db.updateSiteStats(s)
}

func (db *sqlstore) insertSiteStats(s *models.SiteStats) error {
	query := db.Rebind(`INSERT INTO site_stats(site_id, visitors, sessions, pageviews, bounce_rate, avg_duration, known_durations, ts) VALUES(?, ?, ?, ?, ?, ?, ?, ?)`)
	_, err := db.Exec(query, s.SiteID, s.Visitors, s.Sessions, s.Pageviews, s.BounceRate, s.AvgDuration, s.KnownDurations, s.Date.Format(DATE_FORMAT))
	return err
}

func (db *sqlstore) updateSiteStats(s *models.SiteStats) error {
	query := db.Rebind(`UPDATE site_stats SET visitors = ?, sessions = ?, pageviews = ?, bounce_rate = ?, avg_duration = ?, known_durations = ? WHERE site_id = ? AND ts = ?`)
	_, err := db.Exec(query, s.Visitors, s.Sessions, s.Pageviews, s.BounceRate, s.AvgDuration, s.KnownDurations, s.SiteID, s.Date.Format(DATE_FORMAT))
	return err
}

func (db *sqlstore) SelectSiteStats(siteID int64, startDate time.Time, endDate time.Time) ([]*models.SiteStats, error) {
	results := []*models.SiteStats{}
	query := db.Rebind(`SELECT *
		FROM site_stats 
		WHERE site_id = ? AND ts >= ? AND ts <= ? 
		ORDER BY ts DESC`)
	err := db.Select(&results, query, siteID, startDate.Format(DATE_FORMAT), endDate.Format(DATE_FORMAT))
	return results, err
}

func (db *sqlstore) GetAggregatedSiteStats(siteID int64, startDate time.Time, endDate time.Time) (*models.SiteStats, error) {
	stats := &models.SiteStats{}
	query := db.Rebind(`SELECT 
		SUM(pageviews) AS pageviews,
		SUM(visitors) AS visitors,
		SUM(sessions) AS sessions,
		SUM(pageviews*avg_duration) / SUM(pageviews) AS avg_duration,
		COALESCE(SUM(sessions*bounce_rate) / SUM(sessions), 0.00) AS bounce_rate
	 FROM site_stats 
	 WHERE site_id = ? AND ts >= ? AND ts <= ? LIMIT 1`)
	err := db.Get(stats, query, siteID, startDate.Format(DATE_FORMAT), endDate.Format(DATE_FORMAT))
	return stats, mapError(err)
}

func (db *sqlstore) GetRealtimeVisitorCount(siteID int64) (int64, error) {
	var siteTrackingID string
	if err := db.Get(&siteTrackingID, db.Rebind(`SELECT tracking_id FROM sites WHERE id = ? LIMIT 1`), siteID); err != nil && err != sql.ErrNoRows {
		log.Error(err)
		return 0, mapError(err)
	}

	var sql string
	var total int64

	// for backwards compatibility with tracking snippets without an explicit site tracking ID (< 1.1.0)
	if siteID == 1 {
		sql = `SELECT COUNT(*) FROM pageviews p WHERE ( site_tracking_id = ? OR site_tracking_id = '' ) AND is_finished = FALSE AND timestamp > ?`
	} else {
		sql = `SELECT COUNT(*) FROM pageviews p WHERE site_tracking_id = ? AND is_finished = FALSE AND timestamp > ?`
	}

	query := db.Rebind(sql)
	if err := db.Get(&total, query, siteTrackingID, time.Now().Add(-5*time.Minute)); err != nil {
		return 0, mapError(err)
	}

	return total, nil
}
