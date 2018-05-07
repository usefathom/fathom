package datastore

import (
	"database/sql"
	"github.com/usefathom/fathom/pkg/models"
	"time"
)

// SavePageview inserts a single pageview model into the connected database
func SavePageview(p *models.Pageview) error {
	query := dbx.Rebind(`INSERT INTO pageviews(session_id, pathname, is_new_visitor, is_unique, is_bounce, referrer, duration, timestamp) VALUES(?, ?, ?, ?, ?, ?, ?, ?)`)
	result, err := dbx.Exec(query, p.SessionID, p.Pathname, p.IsNewVisitor, p.IsUnique, p.IsBounce, p.Referrer, p.Duration, p.Timestamp)
	if err != nil {
		return err
	}

	p.ID, _ = result.LastInsertId()
	return nil
}

func UpdatePageview(p *models.Pageview) error {
	query := dbx.Rebind(`UPDATE pageviews SET is_bounce = ?, duration = ? WHERE id = ?`)
	_, err := dbx.Exec(query, p.IsBounce, p.Duration, p.ID)
	return err
}

func GetMostRecentPageviewBySessionID(sessionID string) (*models.Pageview, error) {
	result := &models.Pageview{}
	query := dbx.Rebind(`SELECT * FROM pageviews WHERE session_id = ? ORDER BY id DESC LIMIT 1`)
	err := dbx.Get(result, query, sessionID)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoResults
		}

		return nil, err
	}

	return result, nil
}

func GetProcessablePageviews() ([]*models.Pageview, error) {
	var results []*models.Pageview
	thirtyMinsAgo := time.Now().Add(-30 * time.Minute)
	fiveMinsAgo := time.Now().Add(-5 * time.Minute)
	query := dbx.Rebind(`SELECT * FROM pageviews WHERE ( duration > 0 OR timestamp < ? ) AND timestamp < ? LIMIT 500`)
	err := dbx.Select(&results, query, thirtyMinsAgo, fiveMinsAgo)
	return results, err
}
