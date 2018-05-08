package datastore

import (
	"database/sql"
	"strconv"
	"strings"
	"time"

	"github.com/usefathom/fathom/pkg/models"
)

// SavePageview inserts a single pageview model into the connected database
func SavePageview(p *models.Pageview) error {
	query := dbx.Rebind(`INSERT INTO pageviews(hostname, pathname, session_id, is_new_visitor, is_new_session, is_unique, is_bounce, referrer, duration, timestamp) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	result, err := dbx.Exec(query, p.Hostname, p.Pathname, p.SessionID, p.IsNewVisitor, p.IsNewSession, p.IsUnique, p.IsBounce, p.Referrer, p.Duration, p.Timestamp)
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

func DeletePageviews(pageviews []*models.Pageview) error {
	ids := []string{}
	for _, p := range pageviews {
		ids = append(ids, strconv.FormatInt(p.ID, 10))
	}
	query := dbx.Rebind(`DELETE FROM pageviews WHERE id IN(` + strings.Join(ids, ",") + `)`)
	_, err := dbx.Exec(query)
	return err
}
