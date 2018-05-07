package datastore

import (
	"database/sql"
	"github.com/usefathom/fathom/pkg/models"
	"time"
)

const sqlSelectPageStat = `SELECT * FROM daily_page_stats WHERE date = ? AND pathname = ? LIMIT 1`
const sqlInsertPageStats = `INSERT INTO daily_page_stats(views, unique_views, bounced, bounced_n, avg_duration, avg_duration_n, pathname, date) VALUES(?, ?, ?, ?, ?, ?, ?, ?)`
const sqlUpdatePageStats = `UPDATE daily_page_stats SET views = ?, unique_views = ?, bounced = ?, bounced_n = ?, avg_duration = ?, avg_duration_n = ? WHERE pathname = ? AND date = ?`

func GetPageStats(date time.Time, pathname string) (*models.PageStats, error) {
	stats := &models.PageStats{}
	query := dbx.Rebind(sqlSelectPageStat)
	err := dbx.Get(stats, query, date.Format("2006-01-02"), pathname)
	if err != nil && err == sql.ErrNoRows {
		return nil, ErrNoResults
	}
	return stats, err
}

func InsertPageStats(s *models.PageStats) error {
	query := dbx.Rebind(sqlInsertPageStats)
	_, err := dbx.Exec(query, s.Views, s.UniqueViews, s.Bounced, s.BouncedN, s.AvgDuration, s.AvgDurationN, s.Pathname, s.Date.Format("2006-01-02"))
	return err
}

func UpdatePageStats(s *models.PageStats) error {
	query := dbx.Rebind(sqlUpdatePageStats)
	_, err := dbx.Exec(query, s.Views, s.UniqueViews, s.Bounced, s.BouncedN, s.AvgDuration, s.AvgDurationN, s.Pathname, s.Date.Format("2006-01-02"))
	return err
}
