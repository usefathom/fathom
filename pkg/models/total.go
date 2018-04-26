package models

// Total represents a daily aggregated total for a metric
type Total struct {
	ID          int64
	PageID      int64 `db:"page_id"`
	Value       string
	Count       int64
	CountUnique int64  `db:"count_unique"`
	Date        string `db:"date_group"`
}
