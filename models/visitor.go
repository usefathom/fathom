package models

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
)

type Visitor struct {
	ID               int64
	Key              string
	BrowserName      string
	BrowserVersion   string
	BrowserLanguage  string
	Country          string
	DeviceOS         string
	IpAddress        string
	ScreenResolution string
}

func (v *Visitor) Save(conn *sql.DB) error {
	// prepare statement for inserting data
	stmt, err := conn.Prepare(`INSERT INTO visitors (
      visitor_key,
      ip_address,
      device_os,
      browser_name,
      browser_version,
      browser_language,
      screen_resolution,
      country
    ) VALUES( ?, ?, ?, ?, ?, ?, ?, ? )`)
	if err != nil {
		return err
	}
	defer stmt.Close()
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

// GenerateKey generates the "unique" visitor key from date, user agent + screen resolution
func (v *Visitor) GenerateKey(date string, ipAddress string, userAgent string) string {
	byteKey := md5.Sum([]byte(date + ipAddress + userAgent))
	return hex.EncodeToString(byteKey[:])
}
