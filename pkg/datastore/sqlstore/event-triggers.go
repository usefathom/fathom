package sqlstore

import (
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/usefathom/fathom/pkg/models"
)

// GetEventTrigger selects a single envent-trigger by its string ID
func (db *sqlstore) GetEventTrigger(id string) (*models.EventTrigger, error) {
	result := &models.EventTrigger{}
	query := db.Rebind(`SELECT * FROM event_triggers WHERE id = ? LIMIT 1`)
	err := db.Get(result, query, id)

	if err != nil {
		return nil, mapError(err)
	}

	return result, nil
}

// InsertEventTriggers bulks-insert multiple event-triggers using a single INSERT statement
func (db *sqlstore) InsertEventTriggers(events []*models.EventTrigger) error {
	n := len(events)
	if n == 0 {
		return nil
	}

	// generate placeholders string
	placeholderTemplate := "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?),"
	placeholders := strings.Repeat(placeholderTemplate, n)
	placeholders = placeholders[:len(placeholders)-1]
	nPlaceholders := strings.Count(placeholderTemplate, "?")

	// init values slice with correct length
	nValues := n * nPlaceholders
	values := make([]interface{}, nValues)

	// overwrite nil values in slice
	j := 0
	for i := range events {
		// test for columns with ignored values
		if events[i].IsBounce != true {
			log.Warnf("inserting pageview with invalid column values for bulk-insert")
		}

		j = i * nPlaceholders
		values[j] = events[i].ID
		values[j+1] = events[i].SiteTrackingID
		values[j+2] = events[i].Hostname
		values[j+3] = events[i].Pathname
		values[j+4] = events[i].IsNewVisitor
		values[j+5] = events[i].IsNewSession
		values[j+6] = events[i].IsUnique
		values[j+7] = events[i].Referrer
		values[j+8] = events[i].Timestamp
		values[j+9] = events[i].EventName
		values[j+10] = events[i].EventContent
	}

	// string together query & execute with values
	query := `INSERT INTO event_triggers(id, site_tracking_id, hostname, pathname, is_new_visitor, is_new_session, is_unique, referrer, timestamp, event_name, event_content) VALUES ` + placeholders
	query = db.Rebind(query)
	_, err := db.Exec(query, values...)
	if err != nil {
		return err
	}

	return nil
}

func (db *sqlstore) DeleteEventTriggers(events []*models.EventTrigger) error {
	ids := []string{}
	for _, e := range events {
		ids = append(ids, "'"+e.ID+"'")
	}
	query := db.Rebind(`DELETE FROM event_triggers WHERE id IN(` + strings.Join(ids, ",") + `)`)
	_, err := db.Exec(query)
	return err
}
