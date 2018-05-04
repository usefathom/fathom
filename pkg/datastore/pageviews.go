package datastore

import (
	"database/sql"

	"github.com/usefathom/fathom/pkg/models"
)

// SavePageview inserts a single pageview model into the connected database
func SavePageview(pv *models.Pageview) error {
	query := dbx.Rebind(`INSERT INTO pageviews(page_id, visitor_id, referrer_url, referrer_keyword, timestamp) VALUES( ?, ?, ?, ?, ?)`)
	result, err := dbx.Exec(query, pv.PageID, pv.VisitorID, pv.ReferrerUrl, pv.ReferrerKeyword, pv.Timestamp)
	if err != nil {
		return err
	}

	pv.ID, _ = result.LastInsertId()
	return nil
}

// SavePageviews inserts multiple pageview models into the connected database using a transaction
func SavePageviews(pvs []*models.Pageview) error {
	query := dbx.Rebind(`INSERT INTO pageviews(page_id, visitor_id, referrer_url, referrer_keyword, bounced, timestamp ) VALUES( ?, ?, ?, ?, ?, ? )`)
	tx, err := dbx.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, pv := range pvs {
		result, err := stmt.Exec(pv.PageID, pv.VisitorID, pv.ReferrerUrl, pv.ReferrerKeyword, pv.Bounced, pv.Timestamp)
		if err != nil {
			return err
		}

		pv.ID, err = result.LastInsertId()
	}

	err = tx.Commit()
	return err
}

// UpdatePageview updates an existing pageview
func UpdatePageview(p *models.Pageview) error {
	query := dbx.Rebind(`UPDATE pageviews SET bounced = ? WHERE id = ?`)
	_, err := dbx.Exec(query, p.Bounced, p.ID)
	return err
}

// GetPageview retrieves a pageview by its ID
func GetPageview(ID int64) (*models.Pageview, error) {
	p := &models.Pageview{}

	query := dbx.Rebind(`SELECT * FROM pageviews WHERE id = ? LIMIT 1`)
	err := dbx.Get(p, query, ID)

	if err != nil {
		return nil, ErrNoResults
	}

	return p, nil
}

func GetLastPageviewForVisitor(visitorID int64) (*models.Pageview, error) {
	p := &models.Pageview{}
	query := dbx.Rebind(`SELECT * FROM pageviews WHERE visitor_id = ? ORDER BY id DESC LIMIT 1`)
	err := dbx.Get(p, query, visitorID)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoResults
		}

		return nil, err
	}

	return p, nil
}

func PageviewCountPerPageAndDay(before string, after string) ([]*models.Total, error) {
	query := dbx.Rebind(`SELECT
	    	pv.page_id,
	    	COUNT(*) AS count,
			COUNT(DISTINCT(pv.visitor_id)) AS count_unique,
			DATE_FORMAT(pv.timestamp, '%Y-%m-%d') AS date_group
	    FROM pageviews pv
	    WHERE pv.timestamp < ? AND pv.timestamp > ?
	    GROUP BY pv.page_id, date_group`)
	var results []*models.Total
	err := dbx.Select(&results, query, before, after)
	return results, err
}
