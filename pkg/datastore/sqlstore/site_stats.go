package sqlstore

import (
	"database/sql"
	"time"

	"github.com/usefathom/fathom/pkg/models"
)

func (db *sqlstore) GetSiteStats(siteID int64, date time.Time) (*models.SiteStats, error) {
	stats := &models.SiteStats{}
	query := db.Rebind(`SELECT * FROM daily_site_stats WHERE site_id = ? AND date = ? LIMIT 1`)
	err := db.Get(stats, query, siteID, date.Format("2006-01-02"))
	if err != nil && err == sql.ErrNoRows {
		return nil, ErrNoResults
	}
	return stats, err
}

func (db *sqlstore) InsertSiteStats(s *models.SiteStats) error {
	query := db.Rebind(`INSERT INTO daily_site_stats(site_id, visitors, sessions, pageviews, bounce_rate, avg_duration, known_durations, date) VALUES(?, ?, ?, ?, ?, ?, ?)`)
	_, err := db.Exec(query, s.SiteID, s.Visitors, s.Sessions, s.Pageviews, s.BounceRate, s.AvgDuration, s.KnownDurations, s.Date.Format("2006-01-02"))
	return err
}

func (db *sqlstore) UpdateSiteStats(s *models.SiteStats) error {
	query := db.Rebind(`UPDATE daily_site_stats SET visitors = ?, sessions = ?, pageviews = ?, bounce_rate = ROUND(?, 4), avg_duration = ROUND(?, 4), known_durations = ? WHERE site_id = ? AND date = ?`)
	_, err := db.Exec(query, s.Visitors, s.Sessions, s.Pageviews, s.BounceRate, s.AvgDuration, s.KnownDurations, s.SiteID, s.Date.Format("2006-01-02"))
	return err
}

func (db *sqlstore) GetSiteStatsPerDay(siteID int64, startDate time.Time, endDate time.Time) ([]*models.SiteStats, error) {
	results := []*models.SiteStats{}
	sql := `SELECT * FROM daily_site_stats WHERE site_id = ? AND date >= ? AND date <= ?`
	query := db.Rebind(sql)
	err := db.Select(&results, query, siteID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	return results, err
}

func (db *sqlstore) GetAggregatedSiteStats(siteID int64, startDate time.Time, endDate time.Time) (*models.SiteStats, error) {
	stats := &models.SiteStats{}
	query := db.Rebind(`SELECT 
		COALESCE(SUM(pageviews), 0) AS pageviews,
		COALESCE(SUM(visitors), 0) AS visitors,
		COALESCE(SUM(sessions), 0) AS sessions,
		COALESCE(ROUND(SUM(pageviews*avg_duration) / NULLIF(SUM(pageviews), 0), 4), 0.00) AS avg_duration,
		COALESCE(ROUND(SUM(sessions*bounce_rate) / NULLIF(SUM(sessions), 0), 4), 0.00) AS bounce_rate
	 FROM daily_site_stats WHERE site_id = ? AND date >= ? AND date <= ? LIMIT 1`)
	err := db.Get(stats, query, siteID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	if err != nil && err == sql.ErrNoRows {
		return nil, ErrNoResults
	}

	return stats, err
}

func (db *sqlstore) GetTotalSiteViews(siteID int64, startDate time.Time, endDate time.Time) (int64, error) {
	sql := `SELECT COALESCE(SUM(pageviews), 0) FROM daily_site_stats WHERE site_id = ? AND date >= ? AND date <= ?`
	query := db.Rebind(sql)
	var total int64
	err := db.Get(&total, query, siteID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	return total, err
}

func (db *sqlstore) GetTotalSiteVisitors(siteID int64, startDate time.Time, endDate time.Time) (int64, error) {
	sql := `SELECT COALESCE(SUM(visitors), 0) FROM daily_site_stats WHERE site_id = ? AND date >= ? AND date <= ?`
	query := db.Rebind(sql)
	var total int64
	err := db.Get(&total, query, siteID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	return total, err
}

func (db *sqlstore) GetTotalSiteSessions(siteID int64, startDate time.Time, endDate time.Time) (int64, error) {
	sql := `SELECT COALESCE(SUM(sessions), 0) FROM daily_site_stats WHERE site_id = ? AND date >= ? AND date <= ?`
	query := db.Rebind(sql)
	var total int64
	err := db.Get(&total, query, siteID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	return total, err
}

func (db *sqlstore) GetAverageSiteDuration(siteID int64, startDate time.Time, endDate time.Time) (float64, error) {
	sql := `SELECT COALESCE(ROUND(SUM(pageviews*avg_duration)/SUM(pageviews), 4), 0.00) FROM daily_site_stats WHERE site_id = ? AND date >= ? AND date <= ?`
	query := db.Rebind(sql)
	var total float64
	err := db.Get(&total, query, siteID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	return total, err
}

func (db *sqlstore) GetAverageSiteBounceRate(siteID int64, startDate time.Time, endDate time.Time) (float64, error) {
	sql := `SELECT COALESCE(ROUND(SUM(sessions*bounce_rate)/SUM(sessions), 4), 0.00) FROM daily_site_stats WHERE site_id = ? AND date >= ? AND date <= ?`
	query := db.Rebind(sql)
	var total float64
	err := db.Get(&total, query, siteID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	return total, err
}

func (db *sqlstore) GetRealtimeVisitorCount(siteID int64) (int64, error) {
	var siteTrackingID string
	var total int64
	if err := db.Get(&siteTrackingID, db.Rebind(`SELECT tracking_id FROM sites WHERE id = ?`), siteID); err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	sql := `SELECT COUNT(*) FROM pageviews p WHERE site_tracking_id = ? AND ( duration = 0 OR is_bounce = TRUE) AND  timestamp > ?`
	query := db.Rebind(sql)
	if err := db.Get(&total, query, siteTrackingID, time.Now().Add(-5*time.Minute)); err != nil {
		return 0, err
	}

	return total, nil
}
