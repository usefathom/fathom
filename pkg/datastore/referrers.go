package datastore

import "github.com/usefathom/fathom/pkg/models"

func ReferrerCountPerDay(before string, after string) ([]*models.Total, error) {
	var results []*models.Total

	query := dbx.Rebind(`
		SELECT
	      pv.referrer_url AS value,
		  COUNT(*) AS count,
	      COUNT(DISTINCT(pv.visitor_id)) AS count_unique,
	      DATE_FORMAT(pv.timestamp, '%Y-%m-%d') AS date_group
	    FROM pageviews pv
	    WHERE pv.referrer_url IS NOT NULL
	    AND pv.referrer_url != ''
	    AND pv.timestamp < ? AND pv.timestamp > ? 
	    GROUP BY date_group, pv.referrer_url`)

	err := dbx.Select(&results, query, before, after)
	return results, err
}
