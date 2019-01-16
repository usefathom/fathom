package models

import (
	"time"
)

type EventTrigger struct {
	ID             string    `db:"id"`
	SiteTrackingID string    `db:"site_tracking_id"`
	Hostname       string    `db:"hostname"`
	Pathname       string    `db:"pathname"`
	EventName      string    `db:"event_name"`
	EventContent   string    `db:"event_content"`
	IsNewVisitor   bool      `db:"is_new_visitor"`
	IsNewSession   bool      `db:"is_new_session"`
	IsUnique       bool      `db:"is_unique"`
	IsBounce       bool      `db:"is_bounce"`
	IsFinished     bool      `db:"is_finished"`
	Referrer       string    `db:"referrer"`
	Timestamp      time.Time `db:"timestamp"`
}
