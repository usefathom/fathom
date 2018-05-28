package sqlstore

import (
	"database/sql"
	"github.com/usefathom/fathom/pkg/models"
	"time"
)

func (db *sqlstore) GetPageStats(date time.Time, hostname string, pathname string) (*models.PageStats, error) {
	stats := &models.PageStats{}
	query := db.Rebind(`SELECT * FROM daily_page_stats WHERE hostname = ? AND pathname = ? AND date = ? LIMIT 1`)
	err := db.Get(stats, query, hostname, pathname, date.Format("2006-01-02"))
	if err != nil && err == sql.ErrNoRows {
		return nil, ErrNoResults
	}
	return stats, err
}

func (db *sqlstore) InsertPageStats(s *models.PageStats) error {
	query := db.Rebind(`INSERT INTO daily_page_stats(pageviews, visitors, entries, bounce_rate, avg_duration, known_durations, hostname, pathname, date) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	_, err := db.Exec(query, s.Pageviews, s.Visitors, s.Entries, s.BounceRate, s.AvgDuration, s.KnownDurations, s.Hostname, s.Pathname, s.Date.Format("2006-01-02"))
	return err
}

func (db *sqlstore) UpdatePageStats(s *models.PageStats) error {
	query := db.Rebind(`UPDATE daily_page_stats SET pageviews = ?, visitors = ?, entries = ?, bounce_rate = ROUND(?, 4), avg_duration = ROUND(?, 4), known_durations = ? WHERE hostname = ? AND pathname = ? AND date = ?`)
	_, err := db.Exec(query, s.Pageviews, s.Visitors, s.Entries, s.BounceRate, s.AvgDuration, s.KnownDurations, s.Hostname, s.Pathname, s.Date.Format("2006-01-02"))
	return err
}

func (db *sqlstore) GetAggregatedPageStats(startDate time.Time, endDate time.Time, limit int) ([]*models.PageStats, error) {
	var result []*models.PageStats
	query := db.Rebind(`SELECT hostname, pathname, SUM(pageviews) AS pageviews, SUM(visitors) AS visitors, SUM(entries) AS entries, COALESCE(ROUND(SUM(entries*bounce_rate)/SUM(entries), 4), 0.00) AS bounce_rate, COALESCE(ROUND(SUM(avg_duration*pageviews)/SUM(pageviews), 4), 0.00) AS avg_duration FROM daily_page_stats WHERE date >= ? AND date <= ? GROUP BY hostname, pathname ORDER BY pageviews DESC LIMIT ?`)
	err := db.Select(&result, query, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), limit)
	return result, err
}

func (db *sqlstore) GetAggregatedPageStatsPageviews(startDate time.Time, endDate time.Time) (int, error) {
	var result int
	query := db.Rebind(`SELECT COALESCE(SUM(pageviews), 0) FROM daily_page_stats WHERE date >= ? AND date <= ?`)
	err := db.Get(&result, query, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	return result, err
}
