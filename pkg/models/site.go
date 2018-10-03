package models

// Site represents a group for tracking data
type Site struct {
	ID         int64  `db:"id" json:"id"`
	TrackingID string `db:"tracking_id" json:"tracking_id"`
	Name       string `db:"name" json:"name"`
}
