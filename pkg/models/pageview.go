package models

import (
	"time"
)

type Pageview struct {
	ID           string    `db:"id"`
	Hostname     string    `db:"hostname"`
	Pathname     string    `db:"pathname"`
	IsNewVisitor bool      `db:"is_new_visitor"`
	IsNewSession bool      `db:"is_new_session"`
	IsUnique     bool      `db:"is_unique"`
	IsBounce     bool      `db:"is_bounce"`
	Referrer     string    `db:"referrer"`
	Duration     int64     `db:"duration"`
	Timestamp    time.Time `db:"timestamp"`
}
