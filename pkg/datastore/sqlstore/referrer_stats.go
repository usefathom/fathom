package sqlstore

import (
	"database/sql"
	"time"

	"github.com/usefathom/fathom/pkg/models"
)

func (db *sqlstore) GetReferrerStats(date time.Time, url string) (*models.ReferrerStats, error) {
	stats := &models.ReferrerStats{}
	query := db.Rebind(`SELECT * FROM daily_referrer_stats WHERE url = ? AND date = ? LIMIT 1`)
	err := db.Get(stats, query, url, date.Format("2006-01-02"))
	if err != nil && err == sql.ErrNoRows {
		return nil, ErrNoResults
	}
	return stats, err
}

func (db *sqlstore) InsertReferrerStats(s *models.ReferrerStats) error {
	query := db.Rebind(`INSERT INTO daily_referrer_stats(visitors, pageviews, bounce_rate, avg_duration, url, date) VALUES(?, ?, ?, ?, ?, ?)`)
	_, err := db.Exec(query, s.Visitors, s.Pageviews, s.BounceRate, s.AvgDuration, s.URL, s.Date.Format("2006-01-02"))
	return err
}

func (db *sqlstore) UpdateReferrerStats(s *models.ReferrerStats) error {
	query := db.Rebind(`UPDATE daily_referrer_stats SET visitors = ?, pageviews = ?, bounce_rate = ROUND(?, 4), avg_duration = ROUND(?, 4) WHERE url = ? AND date = ?`)
	_, err := db.Exec(query, s.Visitors, s.Pageviews, s.BounceRate, s.AvgDuration, s.URL, s.Date.Format("2006-01-02"))
	return err
}

func (db *sqlstore) GetAggregatedReferrerStats(startDate time.Time, endDate time.Time, limit int) ([]*models.ReferrerStats, error) {
	var result []*models.ReferrerStats
	query := db.Rebind(`SELECT url, SUM(visitors) AS visitors, SUM(pageviews) AS pageviews, COALESCE(ROUND(SUM(pageviews*bounce_rate)/SUM(pageviews), 4), 0.00) AS bounce_rate, COALESCE(ROUND(SUM(avg_duration*pageviews)/SUM(pageviews), 4), 0.00) AS avg_duration FROM daily_referrer_stats WHERE date >= ? AND date <= ? GROUP BY url ORDER BY pageviews DESC LIMIT ?`)
	err := db.Select(&result, query, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), limit)
	return result, err
}

func (db *sqlstore) GetAggregatedReferrerStatsPageviews(startDate time.Time, endDate time.Time) (int, error) {
	var result int
	query := db.Rebind(`SELECT COALESCE(SUM(pageviews), 0) FROM daily_referrer_stats WHERE date >= ? AND date <= ?`)
	err := db.Get(&result, query, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	return result, err
}
