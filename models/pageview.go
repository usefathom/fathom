package models

import (
	"database/sql"
)

type Pageview struct {
	ID              int64
	PageID          int64
	VisitorID       int64
	ReferrerKeyword string
	ReferrerUrl     string
	Timestamp       string
}

func (pv *Pageview) Save(conn *sql.DB) error {
	// prepare statement for inserting data
	stmt, err := conn.Prepare(`INSERT INTO pageviews (
     page_id,
     visitor_id,
     referrer_url,
     referrer_keyword,
     timestamp
   ) VALUES( ?, ?, ?, ?, ? )`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(
		pv.PageID,
		pv.VisitorID,
		pv.ReferrerUrl,
		pv.ReferrerKeyword,
		pv.Timestamp,
	)
	if err != nil {
		return err
	}

	pv.ID, err = result.LastInsertId()
	return err
}
