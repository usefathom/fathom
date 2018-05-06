package datastore

import (
	"database/sql"
	"github.com/usefathom/fathom/pkg/models"
	"time"
)

const sqlSelectSiteStat = `SELECT * FROM daily_site_stats WHERE date = ? LIMIT 1`
const sqlInsertSiteStats = `INSERT INTO daily_site_stats(visitors, pageviews, bounced, bounced_n, avg_duration, avg_duration_n, date) VALUES(?, ?, ?, ?, ?, ?, ?)`

/*
visitors INT NOT NULL,
pageviews INT NOT NULL,
bounced INT NOT NULL,
bounced_n INT NOT NULL,
avg_duration INT NOT NULL,
avg_duration_n INT NOT NULL,
date DATE NOT NULL
*/
func GetSiteStats(date time.Time) (*models.SiteStats, error) {
	stats := &models.SiteStats{}
	query := dbx.Rebind(sqlSelectSiteStat)
	err := dbx.Get(stats, query, date)
	if err != nil && err == sql.ErrNoRows {
		return nil, ErrNoResults
	}
	return stats, err
}

func SaveSiteStats(s *models.SiteStats) error {
	query := dbx.Rebind(sqlInsertSiteStats)
	_, err := dbx.Exec(query, s.Visitors, s.Pageviews, s.Bounced, s.BouncedN, s.AvgDuration, s.AvgDurationN, s.Date)
	return err
}
