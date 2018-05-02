package models

type Count struct {
	URL            string  `json:"url"`
	Views          int64   `json:"views"`
	Uniques        int64   `json:"uniques"`
	PercentOfTotal float64 `json:"percent_of_total"`
}
