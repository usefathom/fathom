package datastore

import (
	"database/sql"
	"github.com/usefathom/fathom/pkg/models"
	"time"
)

const sqlSelectSiteStat = `SELECT * FROM daily_site_stats WHERE date = ? LIMIT 1`
const sqlInsertSiteStats = `INSERT INTO daily_site_stats(visitors, pageviews, bounced, bounced_n, avg_duration, avg_duration_n, date) VALUES(?, ?, ?, ?, ?, ?, ?)`
const sqlUpdateSiteStats = `UPDATE daily_site_stats SET visitors = ?, pageviews = ?, bounced = ?, bounced_n = ?, avg_duration = ?, avg_duration_n = ? WHERE date = ?`

func GetSiteStats(date time.Time) (*models.SiteStats, error) {
	stats := &models.SiteStats{}
	query := dbx.Rebind(sqlSelectSiteStat)
	err := dbx.Get(stats, query, date.Format("2006-01-02"))
	if err != nil && err == sql.ErrNoRows {
		return nil, ErrNoResults
	}
	return stats, err
}

func InsertSiteStats(s *models.SiteStats) error {
	query := dbx.Rebind(sqlInsertSiteStats)
	_, err := dbx.Exec(query, s.Visitors, s.Pageviews, s.Bounced, s.BouncedN, s.AvgDuration, s.AvgDurationN, s.Date.Format("2006-01-02"))
	return err
}

func UpdateSiteStats(s *models.SiteStats) error {
	query := dbx.Rebind(sqlUpdateSiteStats)
	_, err := dbx.Exec(query, s.Visitors, s.Pageviews, s.Bounced, s.BouncedN, s.AvgDuration, s.AvgDurationN, s.Date.Format("2006-01-02"))
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

func GetAverageSiteDuration(startDate time.Time, endDate time.Time) (int64, error) {
	sql := `SELECT COALESCE(ROUND(AVG(avg_duration), 0), 0) FROM daily_site_stats WHERE date >= ? AND date <= ?`
	query := dbx.Rebind(sql)
	var total int64
	err := dbx.Get(&total, query, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	return total, err
}

func GetAverageSiteBounceRate(startDate time.Time, endDate time.Time) (int64, error) {
	sql := `SELECT COALESCE(ROUND(AVG(bounced), 0), 0) FROM daily_site_stats WHERE date >= ? AND date <= ?`
	query := dbx.Rebind(sql)
	var total int64
	err := dbx.Get(&total, query, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	return total, err
}

func GetRealtimeVisitorCount() (int64, error) {
	sql := `SELECT COUNT(DISTINCT(session_id)) FROM raw_pageviews WHERE timestamp > ?`
	query := dbx.Rebind(sql)
	var total int64
	err := dbx.Get(&total, query, time.Now().Add(-5*time.Minute))
	return total, err
}
