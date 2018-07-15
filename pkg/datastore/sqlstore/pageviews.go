package sqlstore

import (
	"database/sql"
	"strings"
	"time"

	"github.com/usefathom/fathom/pkg/models"
)

// GetPageview selects a single pageview by its string ID
func (db *sqlstore) GetPageview(id string) (*models.Pageview, error) {
	result := &models.Pageview{}
	query := db.Rebind(`SELECT * FROM pageviews WHERE id = ? LIMIT 1`)
	err := db.Get(result, query, id)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoResults
		}

		return nil, err
	}

	return result, nil
}

// InsertPageviews inserts multiple pageviews using a single INSERT statement
func (db *sqlstore) InsertPageviews(pageviews []*models.Pageview) error {
	n := len(pageviews)
	if n == 0 {
		return nil
	}

	placeholders := make([]string, 0, n)
	values := make([]interface{}, 0, n*10)

	for i := 0; i < n; i++ {
		placeholders = append(placeholders, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		values = append(values, pageviews[i].ID, pageviews[i].Hostname, pageviews[i].Pathname, pageviews[i].IsNewVisitor, pageviews[i].IsNewSession, pageviews[i].IsUnique, pageviews[i].IsBounce, pageviews[i].Referrer, pageviews[i].Duration, pageviews[i].Timestamp)
	}

	query := `INSERT INTO pageviews(id, hostname, pathname, is_new_visitor, is_new_session, is_unique, is_bounce, referrer, duration, timestamp) VALUES `
	query = query + strings.Join(placeholders, ",")

	query = db.Rebind(query)
	_, err := db.Exec(query, values...)
	if err != nil {
		return err
	}

	return nil
}

// UpdatePageviews updates multiple pageviews using a single transaction
// Please note that this function only updates the IsBounce and Duration properties
func (db *sqlstore) UpdatePageviews(pageviews []*models.Pageview) error {
	if len(pageviews) == 0 {
		return nil
	}

	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	query := tx.Rebind(`UPDATE pageviews SET is_bounce = ?, duration = ? WHERE id = ?`)
	stmt, err := tx.Preparex(query)
	if err != nil {
		return err
	}

	for i := 0; i < len(pageviews); i++ {
		_, err = stmt.Exec(query, pageviews[i].IsBounce, pageviews[i].Duration)
	}

	err = tx.Commit()
	return err
}

// GetProcessablePageviews selects all pageviews which are "done" (ie not still waiting for bounce flag or duration)
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
		ids = append(ids, "'"+p.ID+"'")
	}
	query := db.Rebind(`DELETE FROM pageviews WHERE id IN(` + strings.Join(ids, ",") + `)`)
	_, err := db.Exec(query)
	return err
}
