package datastore

import (
	"database/sql"

	"github.com/usefathom/fathom/pkg/models"
)

// GetVisitorByKey ...
func GetVisitorByKey(key string) (*models.Visitor, error) {
	v := &models.Visitor{}
	query := dbx.Rebind(`SELECT v.id FROM visitors v WHERE v.visitor_key = ? LIMIT 1`)
	err := dbx.Get(v, query, key)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoResults
		}

		return nil, err
	}

	return v, nil
}

// SaveVisitor inserts a single visitor model into the connected database
func SaveVisitor(v *models.Visitor) error {
	query := dbx.Rebind(`INSERT INTO visitors(visitor_key, device_os, browser_name, browser_version, browser_language, screen_resolution, country) VALUES( ?, ?, ?, ?, ?, ?, ? )`)
	result, err := dbx.Exec(query, v.Key, v.DeviceOS, v.BrowserName, v.BrowserVersion, v.BrowserLanguage, v.ScreenResolution, v.Country)
	if err != nil {
		return err
	}

	v.ID, _ = result.LastInsertId()
	return nil
}

// RealtimeVisitors returns the total number of visitors in the last 3 minutes
// TODO: Query visitors table instead, using a last_seen column.
func RealtimeVisitors() (int, error) {
	var result int
	query := dbx.Rebind(`
		SELECT COUNT(DISTINCT(pv.visitor_id))
		FROM pageviews pv
		WHERE pv.timestamp >= DATE_SUB(CURRENT_TIMESTAMP, INTERVAL 3 HOUR_MINUTE) AND pv.timestamp <= CURRENT_TIMESTAMP`)
	err := dbx.Get(&result, query)
	return result, err
}

func VisitorCountPerDay(before string, after string) ([]*models.Total, error) {
	var results []*models.Total

	query := dbx.Rebind(`
		SELECT
		  COUNT(DISTINCT(pv.visitor_id)) AS count,
		  DATE_FORMAT(pv.timestamp, '%Y-%m-%d') AS date_group
		FROM pageviews pv
		WHERE pv.timestamp < ? AND pv.timestamp > ?
		GROUP BY date_group`)

	err := dbx.Select(&results, query, before, after)
	return results, err
}
