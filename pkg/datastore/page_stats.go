package datastore

import (
	"database/sql"
	"github.com/usefathom/fathom/pkg/models"
	"time"
)

func GetPageStats(date time.Time, hostname string, pathname string) (*models.PageStats, error) {
	stats := &models.PageStats{}
	query := dbx.Rebind(`SELECT * FROM daily_page_stats WHERE hostname = ? AND pathname = ? AND date = ? LIMIT 1`)
	err := dbx.Get(stats, query, hostname, pathname, date.Format("2006-01-02"))
	if err != nil && err == sql.ErrNoRows {
		return nil, ErrNoResults
	}
	return stats, err
}

func InsertPageStats(s *models.PageStats) error {
	query := dbx.Rebind(`INSERT INTO daily_page_stats(pageviews, visitors, entries, bounce_rate, avg_duration, hostname, pathname, date) VALUES(?, ?, ?, ?, ?, ?, ?, ?)`)
	_, err := dbx.Exec(query, s.Pageviews, s.Visitors, s.Entries, s.BounceRate, s.AvgDuration, s.Hostname, s.Pathname, s.Date.Format("2006-01-02"))
	return err
}

func UpdatePageStats(s *models.PageStats) error {
	query := dbx.Rebind(`UPDATE daily_page_stats SET pageviews = ?, visitors = ?, entries = ?, bounce_rate = ?, avg_duration = ? WHERE hostname = ? AND pathname = ? AND date = ?`)
	_, err := dbx.Exec(query, s.Pageviews, s.Visitors, s.Entries, s.BounceRate, s.AvgDuration, s.Hostname, s.Pathname, s.Date.Format("2006-01-02"))
	return err
}

func GetAggregatedPageStats(startDate time.Time, endDate time.Time, limit int) ([]*models.PageStats, error) {
	var result []*models.PageStats
	query := dbx.Rebind(`SELECT hostname, pathname, SUM(pageviews) AS pageviews, SUM(visitors) AS visitors, SUM(entries) AS entries, ROUND(AVG(bounce_rate), 0) AS bounce_rate FROM daily_page_stats WHERE date >= ? AND date <= ? GROUP BY hostname, pathname, pageviews, visitors, entries, bounce_rate ORDER BY pageviews DESC LIMIT ?`)
	err := dbx.Select(&result, query, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), limit)
	return result, err
}
