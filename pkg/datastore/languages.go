package datastore

import "github.com/usefathom/fathom/pkg/models"

func LanguageCountPerDay(before string, after string) ([]*models.Total, error) {
	var results []*models.Total

	query := dbx.Rebind(`
		SELECT
		   v.browser_language AS value,
		   COUNT(*) AS count,
		   COUNT(DISTINCT(pv.visitor_id)) AS count_unique,
		   DATE_FORMAT(pv.timestamp, "%Y-%m-%d") AS date_group
		FROM pageviews pv
		LEFT JOIN visitors v ON v.id = pv.visitor_id
		WHERE pv.timestamp < ? AND pv.timestamp > ?
		GROUP BY date_group, v.browser_language`)

	err := dbx.Select(&results, query, before, after)
	return results, err
}
