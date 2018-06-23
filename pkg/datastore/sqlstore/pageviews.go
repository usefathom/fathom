package sqlstore

import (
	"database/sql"
	"strconv"
	"strings"
	"time"

	"github.com/usefathom/fathom/pkg/models"
)

// SavePageview inserts a single pageview model into the connected database
func (db *sqlstore) SavePageview(p *models.Pageview) error {
	query := db.Rebind(`INSERT INTO pageviews(hostname, pathname, session_id, is_new_visitor, is_new_session, is_unique, is_bounce, referrer, duration, timestamp) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	result, err := db.Exec(query, p.Hostname, p.Pathname, p.SessionID, p.IsNewVisitor, p.IsNewSession, p.IsUnique, p.IsBounce, p.Referrer, p.Duration, p.Timestamp)
	if err != nil {
		return err
	}

	p.ID, _ = result.LastInsertId()
	return nil
}

func (db *sqlstore) UpdatePageview(p *models.Pageview) error {
	query := db.Rebind(`UPDATE pageviews SET is_bounce = ?, duration = ? WHERE id = ?`)
	_, err := db.Exec(query, p.IsBounce, p.Duration, p.ID)
	return err
}

func (db *sqlstore) GetMostRecentPageviewBySessionID(sessionID string) (*models.Pageview, error) {
	result := &models.Pageview{}
	query := db.Rebind(`SELECT * FROM pageviews WHERE session_id = ? ORDER BY id DESC LIMIT 1`)
	err := db.Get(result, query, sessionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoResults
		}

		return nil, err
	}

	return result, nil
}

func (db *sqlstore) GetProcessablePageviews() ([]*models.Pageview, error) {
	var results []*models.Pageview
	thirtyMinsAgo := time.Now().Add(-30 * time.Minute)
	// We use FALSE here, even though SQLite has no BOOLEAN value. If it fails, maybe we can roll our own Rebind?
	query := db.Rebind(`SELECT * FROM pageviews WHERE ( duration > 0 AND is_bounce = FALSE ) OR timestamp < ? LIMIT 500`)
	err := db.Select(&results, query, thirtyMinsAgo)
	return results, err
}

func (db *sqlstore) DeletePageviews(pageviews []*models.Pageview) error {
	ids := []string{}
	for _, p := range pageviews {
		ids = append(ids, strconv.FormatInt(p.ID, 10))
	}
	query := db.Rebind(`DELETE FROM pageviews WHERE id IN(` + strings.Join(ids, ",") + `)`)
	_, err := db.Exec(query)
	return err
}
