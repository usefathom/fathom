package datastore

import "github.com/usefathom/fathom/pkg/models"

func BouncesCountPerPageAndDay(before string, after string) ([]*models.Total, error) {
	query := dbx.Rebind(`SELECT
		pv.page_id,
		( COUNT(*) * 100 ) DIV ( SELECT ( COUNT(*) ) FROM pageviews WHERE page_id = pv.page_id AND bounced IS NOT NULL ) AS count,
		( COUNT(DISTINCT(pv.visitor_id)) * 100 ) DIV ( SELECT ( COUNT(*) ) FROM pageviews WHERE page_id = pv.page_id AND bounced IS NOT NULL ) AS count_unique,
		DATE_FORMAT(pv.timestamp, '%Y-%m-%d') AS date_group
		FROM pageviews pv
		WHERE pv.bounced = 1 AND pv.bounced IS NOT NULL AND pv.timestamp < ? AND pv.timestamp > ? 
		GROUP BY pv.page_id, date_group`)
	var results []*models.Total
	err := dbx.Select(&results, query, before, after)
	return results, err
}
