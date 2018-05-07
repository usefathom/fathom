package datastore

import (
	"database/sql"
	"github.com/usefathom/fathom/pkg/models"
	"time"
)

func GetSiteStats(date time.Time) (*models.SiteStats, error) {
	stats := &models.SiteStats{}
	query := dbx.Rebind(`SELECT * FROM daily_site_stats WHERE date = ? LIMIT 1`)
	err := dbx.Get(stats, query, date.Format("2006-01-02"))
	if err != nil && err == sql.ErrNoRows {
		return nil, ErrNoResults
	}
	return stats, err
}

func InsertSiteStats(s *models.SiteStats) error {
	query := dbx.Rebind(`INSERT INTO daily_site_stats(visitors, sessions, pageviews, bounces, avg_duration, date) VALUES(?, ?, ?, ?, ?, ?)`)
	_, err := dbx.Exec(query, s.Visitors, s.Sessions, s.Pageviews, s.Bounces, s.AvgDuration, s.Date.Format("2006-01-02"))
	return err
}

func UpdateSiteStats(s *models.SiteStats) error {
	query := dbx.Rebind(`UPDATE daily_site_stats SET visitors = ?, sessions = ?, pageviews = ?, bounces = ?, avg_duration = ? WHERE date = ?`)
	_, err := dbx.Exec(query, s.Visitors, s.Sessions, s.Pageviews, s.Bounces, s.AvgDuration, s.Date.Format("2006-01-02"))
	return err
}

func GetTotalSiteViews(startDate time.Time, endDate time.Time) (int64, error) {
	sql := `SELECT COALESCE(SUM(pageviews), 0) FROM daily_site_stats WHERE date >= ? AND date <= ?`
	query := dbx.Rebind(sql)
	var total int64
	err := dbx.Get(&total, query, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	return total, err
}

func GetTotalSiteVisitors(startDate time.Time, endDate time.Time) (int64, error) {
	sql := `SELECT COALESCE(SUM(visitors), 0) FROM daily_site_stats WHERE date >= ? AND date <= ?`
	query := dbx.Rebind(sql)
	var total int64
	err := dbx.Get(&total, query, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	return total, err
}

func GetTotalSiteSessions(startDate time.Time, endDate time.Time) (int64, error) {
	sql := `SELECT COALESCE(SUM(sessions), 0) FROM daily_site_stats WHERE date >= ? AND date <= ?`
	query := dbx.Rebind(sql)
	var total int64
	err := dbx.Get(&total, query, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	return total, err
}

func GetAverageSiteDuration(startDate time.Time, endDate time.Time) (int64, error) {
	sql := `SELECT COALESCE(ROUND(AVG(avg_duration), 0), 0) FROM daily_site_stats WHERE date >= ? AND date <= ?`
	query := dbx.Rebind(sql)
	var total int64
	err := dbx.Get(&total, query, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	return total, err
}

func GetAverageSiteBounceRate(startDate time.Time, endDate time.Time) (int64, error) {
	sql := `SELECT COALESCE(ROUND(AVG(bounces), 0), 0) FROM daily_site_stats WHERE date >= ? AND date <= ?`
	query := dbx.Rebind(sql)
	var total int64
	err := dbx.Get(&total, query, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	return total, err
}

func GetRealtimeVisitorCount() (int64, error) {
	sql := `SELECT COUNT(DISTINCT(session_id)) FROM pageviews WHERE timestamp > ?`
	query := dbx.Rebind(sql)
	var total int64
	err := dbx.Get(&total, query, time.Now().Add(-5*time.Minute))
	return total, err
}
