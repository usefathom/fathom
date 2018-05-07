package datastore

import (
	"time"

	"github.com/usefathom/fathom/pkg/models"
)

func GetAggregatedReferrerStats(startDate time.Time, endDate time.Time, limit int) ([]*models.ReferrerStats, error) {
	var result []*models.ReferrerStats
	query := dbx.Rebind(`SELECT url, SUM(visitors) AS visitors, SUM(pageviews) AS pageviews FROM daily_referrer_stats WHERE date >= ? AND date <= ? GROUP BY url, visitors, pageviews ORDER BY pageviews DESC LIMIT ?`)
	err := dbx.Select(&result, query, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), limit)
	return result, err
}
