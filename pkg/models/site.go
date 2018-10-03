package models

// Site represents a group for tracking data
type Site struct {
	ID         int64  `db:"id"`
	TrackingID string `db:"tracking_id"`
	Name       string `db:"name"`
}
