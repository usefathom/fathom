package sqlstore

import (
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/usefathom/fathom/pkg/models"
)

// GetPageview selects a single pageview by its string ID
func (db *sqlstore) GetPageview(id string) (*models.Pageview, error) {
	result := &models.Pageview{}
	query := db.Rebind(`SELECT * FROM pageviews WHERE id = ? LIMIT 1`)
	err := db.Get(result, query, id)

	if err != nil {
		return nil, mapError(err)
	}

	return result, nil
}

// InsertPageviews bulks-insert multiple pageviews using a single INSERT statement
// IMPORTANT: This does not insert the actual IsBounce, Duration and IsFinished values
func (db *sqlstore) InsertPageviews(pageviews []*models.Pageview) error {
	n := len(pageviews)
	if n == 0 {
		return nil
	}

	// generate placeholders string
	placeholderTemplate := "(?, ?, ?, ?, ?, ?, ?, ?, ?, TRUE, FALSE, 0),"
	placeholders := strings.Repeat(placeholderTemplate, n)
	placeholders = placeholders[:len(placeholders)-1]
	nPlaceholders := strings.Count(placeholderTemplate, "?")

	// init values slice with correct length
	nValues := n * nPlaceholders
	values := make([]interface{}, nValues)

	// overwrite nil values in slice
	j := 0
	for i := range pageviews {

		// test for columns with ignored values
		if pageviews[i].IsBounce != true || pageviews[i].Duration > 0 || pageviews[i].IsFinished != false {
			log.Warnf("inserting pageview with invalid column values for bulk-insert")
		}

		j = i * nPlaceholders
		values[j] = pageviews[i].ID
		values[j+1] = pageviews[i].SiteTrackingID
		values[j+2] = pageviews[i].Hostname
		values[j+3] = pageviews[i].Pathname
		values[j+4] = pageviews[i].IsNewVisitor
		values[j+5] = pageviews[i].IsNewSession
		values[j+6] = pageviews[i].IsUnique
		values[j+7] = pageviews[i].Referrer
		values[j+8] = pageviews[i].Timestamp
	}

	// string together query & execute with values
	query := `INSERT INTO pageviews(id, site_tracking_id, hostname, pathname, is_new_visitor, is_new_session, is_unique, referrer, timestamp, is_bounce, is_finished, duration) VALUES ` + placeholders
	query = db.Rebind(query)
	_, err := db.Exec(query, values...)
	if err != nil {
		return err
	}

	return nil
}

// UpdatePageviews updates multiple pageviews using a single transaction
// IMPORTANT: this function only updates the IsFinished, IsBounce and Duration values
func (db *sqlstore) UpdatePageviews(pageviews []*models.Pageview) error {
	if len(pageviews) == 0 {
		return nil
	}

	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	query := tx.Rebind(`UPDATE pageviews SET is_bounce = ?, duration = ?, is_finished = ? WHERE id = ?`)
	stmt, err := tx.Preparex(query)
	if err != nil {
		return err
	}

	for i := range pageviews {
		_, err := stmt.Exec(pageviews[i].IsBounce, pageviews[i].Duration, pageviews[i].IsFinished, pageviews[i].ID)

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	return err
}

// GetProcessablePageviews selects all pageviews which are "done" (ie not still waiting for bounce flag or duration)
func (db *sqlstore) GetProcessablePageviews() ([]*models.Pageview, error) {
	var results []*models.Pageview
	thirtyMinsAgo := time.Now().Add(-30 * time.Minute)
	query := db.Rebind(`SELECT * FROM pageviews WHERE is_finished = TRUE OR timestamp < ? LIMIT 5000`)
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
