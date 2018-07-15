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

	// generate placeholders string
	placeholders := strings.Repeat("(?, ?, ?, ?, ?, ?, ?, ?, ?, ?),", n)
	placeholders = placeholders[:len(placeholders)-1]

	// init values slice with correct length
	nValues := n * 10
	values := make([]interface{}, nValues)

	// overwrite nil values in slice
	j := 0
	for i := range pageviews {
		j = i * 10
		values[j] = pageviews[i].ID
		values[j+1] = pageviews[i].Hostname
		values[j+2] = pageviews[i].Pathname
		values[j+3] = pageviews[i].IsNewVisitor
		values[j+4] = pageviews[i].IsNewSession
		values[j+5] = pageviews[i].IsUnique
		values[j+6] = pageviews[i].IsBounce
		values[j+7] = pageviews[i].Referrer
		values[j+8] = pageviews[i].Duration
		values[j+9] = pageviews[i].Timestamp
	}

	// string together query & execute with values
	query := `INSERT INTO pageviews(id, hostname, pathname, is_new_visitor, is_new_session, is_unique, is_bounce, referrer, duration, timestamp) VALUES ` + placeholders
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

	for i := range pageviews {
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
