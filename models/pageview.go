package models

type Pageview struct {
	ID              int64
	PageID          int64
	VisitorID       int64
	ReferrerKeyword string
	ReferrerUrl     string
	Timestamp       string
}
