package datastore

import (
	"github.com/usefathom/fathom/pkg/models"
)

const sqlInsertRawPageview = `INSERT INTO raw_pageviews(session_id, pathname, is_new_visitor, is_unique, is_bounce, referrer, duration, timestamp) VALUES(?, ?, ?, ?, ?, ?, ?, ?)`
const sqlSelectRawPageviews = `SELECT * FROM raw_pageviews`

// SaveRawPageview inserts a single pageview model into the connected database
func SaveRawPageview(p *models.RawPageview) error {
	query := dbx.Rebind(sqlInsertRawPageview)
	result, err := dbx.Exec(query, p.SessionID, p.Pathname, p.IsNewVisitor, p.IsUnique, p.IsBounce, p.Referrer, p.Duration, p.Timestamp)
	if err != nil {
		return err
	}

	p.ID, _ = result.LastInsertId()
	return nil
}

// SaveRawPageviews inserts multiple pageviews
func SaveRawPageviews(p []*models.RawPageview) error {
	return nil // TODO: Implement this method
}

func GetRawPageviews() ([]*models.RawPageview, error) {
	var results []*models.RawPageview
	query := dbx.Rebind(sqlSelectRawPageviews)
	err := dbx.Select(&results, query)
	return results, err
}
