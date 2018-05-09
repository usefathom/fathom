package datastore

import (
	"database/sql"
	"time"

	"github.com/usefathom/fathom/pkg/models"
)

func GetReferrerStats(date time.Time, url string) (*models.ReferrerStats, error) {
	stats := &models.ReferrerStats{}
	query := dbx.Rebind(`SELECT * FROM daily_referrer_stats WHERE url = ? AND date = ? LIMIT 1`)
	err := dbx.Get(stats, query, url, date.Format("2006-01-02"))
	if err != nil && err == sql.ErrNoRows {
		return nil, ErrNoResults
	}
	return stats, err
}

func InsertReferrerStats(s *models.ReferrerStats) error {
	query := dbx.Rebind(`INSERT INTO daily_referrer_stats(visitors, pageviews, bounce_rate, avg_duration, url, date) VALUES(?, ?, ?, ?, ?, ?)`)
	_, err := dbx.Exec(query, s.Visitors, s.Pageviews, s.BounceRate, s.AvgDuration, s.URL, s.Date.Format("2006-01-02"))
	return err
}

func UpdateReferrerStats(s *models.ReferrerStats) error {
	query := dbx.Rebind(`UPDATE daily_referrer_stats SET visitors = ?, pageviews = ?, bounce_rate = ROUND(?, 4), avg_duration = ROUND(?, 4) WHERE url = ? AND date = ?`)
	_, err := dbx.Exec(query, s.Visitors, s.Pageviews, s.BounceRate, s.AvgDuration, s.URL, s.Date.Format("2006-01-02"))
	return err
}

func GetAggregatedReferrerStats(startDate time.Time, endDate time.Time, limit int) ([]*models.ReferrerStats, error) {
	var result []*models.ReferrerStats
	query := dbx.Rebind(`SELECT url, SUM(visitors) AS visitors, SUM(pageviews) AS pageviews, ROUND(AVG(avg_duration), 4) AS avg_duration, ROUND(AVG(bounce_rate), 4) AS bounce_rate FROM daily_referrer_stats WHERE date >= ? AND date <= ? GROUP BY url ORDER BY pageviews DESC LIMIT ?`)
	err := dbx.Select(&result, query, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), limit)
	return result, err
}
