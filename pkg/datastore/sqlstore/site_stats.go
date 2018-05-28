package sqlstore

import (
	"database/sql"
	"github.com/usefathom/fathom/pkg/models"
	"time"
)

func (db *sqlstore) GetSiteStats(date time.Time) (*models.SiteStats, error) {
	stats := &models.SiteStats{}
	query := db.Rebind(`SELECT * FROM daily_site_stats WHERE date = ? LIMIT 1`)
	err := db.Get(stats, query, date.Format("2006-01-02"))
	if err != nil && err == sql.ErrNoRows {
		return nil, ErrNoResults
	}
	return stats, err
}

func (db *sqlstore) InsertSiteStats(s *models.SiteStats) error {
	query := db.Rebind(`INSERT INTO daily_site_stats(visitors, sessions, pageviews, bounce_rate, avg_duration, known_durations, date) VALUES(?, ?, ?, ?, ?, ?, ?)`)
	_, err := db.Exec(query, s.Visitors, s.Sessions, s.Pageviews, s.BounceRate, s.AvgDuration, s.KnownDurations, s.Date.Format("2006-01-02"))
	return err
}

func (db *sqlstore) UpdateSiteStats(s *models.SiteStats) error {
	query := db.Rebind(`UPDATE daily_site_stats SET visitors = ?, sessions = ?, pageviews = ?, bounce_rate = ROUND(?, 4), avg_duration = ROUND(?, 4), known_durations = ? WHERE date = ?`)
	_, err := db.Exec(query, s.Visitors, s.Sessions, s.Pageviews, s.BounceRate, s.AvgDuration, s.KnownDurations, s.Date.Format("2006-01-02"))
	return err
}

func (db *sqlstore) GetSiteStatsPerDay(startDate time.Time, endDate time.Time) ([]*models.SiteStats, error) {
	results := []*models.SiteStats{}
	sql := `SELECT * FROM daily_site_stats WHERE date >= ? AND date <= ?`
	query := db.Rebind(sql)
	err := db.Select(&results, query, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	return results, err
}

func (db *sqlstore) GetTotalSiteViews(startDate time.Time, endDate time.Time) (int, error) {
	sql := `SELECT COALESCE(SUM(pageviews), 0) FROM daily_site_stats WHERE date >= ? AND date <= ?`
	query := db.Rebind(sql)
	var total int
	err := db.Get(&total, query, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	return total, err
}

func (db *sqlstore) GetTotalSiteVisitors(startDate time.Time, endDate time.Time) (int, error) {
	sql := `SELECT COALESCE(SUM(visitors), 0) FROM daily_site_stats WHERE date >= ? AND date <= ?`
	query := db.Rebind(sql)
	var total int
	err := db.Get(&total, query, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	return total, err
}

func (db *sqlstore) GetTotalSiteSessions(startDate time.Time, endDate time.Time) (int, error) {
	sql := `SELECT COALESCE(SUM(sessions), 0) FROM daily_site_stats WHERE date >= ? AND date <= ?`
	query := db.Rebind(sql)
	var total int
	err := db.Get(&total, query, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	return total, err
}

func (db *sqlstore) GetAverageSiteDuration(startDate time.Time, endDate time.Time) (float64, error) {
	sql := `SELECT COALESCE(ROUND(SUM(pageviews*avg_duration)/SUM(pageviews), 4), 0.00) FROM daily_site_stats WHERE date >= ? AND date <= ?`
	query := db.Rebind(sql)
	var total float64
	err := db.Get(&total, query, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	return total, err
}

func (db *sqlstore) GetAverageSiteBounceRate(startDate time.Time, endDate time.Time) (float64, error) {
	sql := `SELECT COALESCE(ROUND(SUM(sessions*bounce_rate)/SUM(sessions), 4), 0.00) FROM daily_site_stats WHERE date >= ? AND date <= ?`
	query := db.Rebind(sql)
	var total float64
	err := db.Get(&total, query, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	return total, err
}

func (db *sqlstore) GetRealtimeVisitorCount() (int, error) {
	sql := `SELECT COUNT(DISTINCT(session_id)) FROM pageviews WHERE timestamp > ?`
	query := db.Rebind(sql)
	var total int
	err := db.Get(&total, query, time.Now().Add(-5*time.Minute))
	return total, err
}
