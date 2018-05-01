package models

// Point represents a data point, will always have a Label and Value
type Point struct {
	Label           string  `json:"label"`
	Value           int     `json:"value"`
	PercentageValue float64 `json:"perc_value,omitempty"`
	UniqueValue     int     `json:"unique_value,omitempty"`
}
