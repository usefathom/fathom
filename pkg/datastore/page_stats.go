package datastore

import (
	"database/sql"
	"github.com/usefathom/fathom/pkg/models"
	"time"
)

func GetPageStats(date time.Time, pathname string) (*models.PageStats, error) {
	stats := &models.PageStats{}
	query := dbx.Rebind(`SELECT * FROM daily_page_stats WHERE date = ? AND pathname = ? LIMIT 1`)
	err := dbx.Get(stats, query, date.Format("2006-01-02"), pathname)
	if err != nil && err == sql.ErrNoRows {
		return nil, ErrNoResults
	}
	return stats, err
}

func InsertPageStats(s *models.PageStats) error {
	query := dbx.Rebind(`INSERT INTO daily_page_stats(views, entries, unique_views, bounces, avg_duration, pathname, date) VALUES(?, ?, ?, ?, ?, ?, ?)`)
	_, err := dbx.Exec(query, s.Views, s.Entries, s.UniqueViews, s.Bounces, s.AvgDuration, s.Pathname, s.Date.Format("2006-01-02"))
	return err
}

func UpdatePageStats(s *models.PageStats) error {
	query := dbx.Rebind(`UPDATE daily_page_stats SET views = ?, entries = ?, unique_views = ?, bounces = ?, avg_duration = ? WHERE pathname = ? AND date = ?`)
	_, err := dbx.Exec(query, s.Views, s.Entries, s.UniqueViews, s.Bounces, s.AvgDuration, s.Pathname, s.Date.Format("2006-01-02"))
	return err
}

func GetAggregatedPageStats(startDate time.Time, endDate time.Time, limit int) ([]*models.PageStats, error) {
	var result []*models.PageStats
	query := dbx.Rebind(`SELECT pathname, SUM(views) AS views, SUM(unique_views) AS unique_views, SUM(entries) AS entries, ROUND(AVG(bounces), 0) AS bounces FROM daily_page_stats WHERE date >= ? AND date <= ? GROUP BY pathname, views, unique_views, entries, bounces ORDER BY views DESC LIMIT ?`)
	err := dbx.Select(&result, query, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), limit)
	return result, err
}
