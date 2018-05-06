package models

import (
	"time"
)

type RawPageview struct {
	ID           int64     `db:"id"`
	SessionID    string    `db:"session_id"`
	Pathname     string    `db:"pathname"`
	IsNewVisitor bool      `db:"is_new_visitor"`
	IsUnique     bool      `db:"is_unique"`
	IsBounce     bool      `db:"is_bounce"`
	Referrer     string    `db:"referrer"`
	Duration     int64     `db:"duration"`
	Timestamp    time.Time `db:"timestamp"`
}
