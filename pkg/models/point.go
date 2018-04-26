package models

// Point represents a data point, will always have a Label and Value
type Point struct {
	Label           string
	Value           int
	PercentageValue float64
}
