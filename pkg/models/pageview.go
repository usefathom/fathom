package models

type Pageview struct {
	ID              int64
	PageID          int64
	VisitorID       int64
	ReferrerKeyword string
	ReferrerUrl     string
	Timestamp       string
}

type PageviewCount struct {
	Hostname    string `json:"hostname"`
	Path        string `json:"path"`
	Count       int    `json:"count"`
	CountUnique int    `json:"count_unique"`
}
