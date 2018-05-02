package models

// Total represents a daily aggregated total for a metric
type Total struct {
	ID                int64   `json:"-"`
	PageID            int64   `db:"page_id" json:"-"`
	Value             string  `db:"value" json:"value"`
	Count             int64   `db:"count" json:"count"`
	CountUnique       int64   `db:"count_unique" json:"count_unique"`
	PercentageOfTotal float64 `db:"-" json:"percentage_of_total"`
	Date              string  `db:"date_group" json:"date,omitempty"`
}
