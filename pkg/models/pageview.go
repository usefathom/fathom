package models

type Pageview struct {
	ID              int64  `db:"id"`
	PageID          int64  `db:"page_id"`
	VisitorID       int64  `db:"visitor_id"`
	Bounced         bool   `db:"bounced"`
	ReferrerKeyword string `db:"referrer_keyword"`
	ReferrerUrl     string `db:"referrer_url"`
	Timestamp       string `db:"timestamp"`
}
