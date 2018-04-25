package datastore

import (
	"database/sql"
	"github.com/usefathom/fathom/pkg/models"
)

var v models.Visitor

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
	query := dbx.Rebind(`INSERT INTO visitors(visitor_key, ip_address, device_os, browser_name, browser_version, browser_language, screen_resolution, country) VALUES( ?, ?, ?, ?, ?, ?, ?, ? )`)
	result, err := dbx.Exec(query, v.Key, v.IpAddress, v.DeviceOS, v.BrowserName, v.BrowserVersion, v.BrowserLanguage, v.ScreenResolution, v.Country)
	if err != nil {
		return err
	}

	v.ID, _ = result.LastInsertId()
	return nil
}
