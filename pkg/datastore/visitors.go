package datastore

import (
	"github.com/usefathom/fathom/pkg/models"
)

var v models.Visitor

// GetVisitorByKey ...
func GetVisitorByKey(key string) (*models.Visitor, error) {
	// query by unique visitor key
	stmt, err := DB.Prepare("SELECT v.id FROM visitors v WHERE v.visitor_key = ? LIMIT 1")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	err = stmt.QueryRow(key).Scan(&v.ID)
	return &v, err
}

// SaveVisitor ...
func SaveVisitor(v *models.Visitor) error {
	// prepare statement for inserting data
	stmt, err := DB.Prepare(`INSERT INTO visitors (
      visitor_key,
      ip_address,
      device_os,
      browser_name,
      browser_version,
      browser_language,
      screen_resolution,
      country
    ) VALUES( ?, ?, ?, ?, ?, ?, ?, ? )`)
	defer stmt.Close()
	if err != nil {
		return err
	}

	result, err := stmt.Exec(
		v.Key,
		v.IpAddress,
		v.DeviceOS,
		v.BrowserName,
		v.BrowserVersion,
		v.BrowserLanguage,
		v.ScreenResolution,
		v.Country,
	)
	if err != nil {
		return err
	}

	v.ID, err = result.LastInsertId()
	return err
}
