package count

import (
	"testing"
	"time"
)

func TestFill(t *testing.T) {
	start := time.Date(2016, time.January, 1, 0, 0, 0, 0, time.Local)
	end := time.Date(2016, time.January, 5, 0, 0, 0, 0, time.Local)

	points := []Point{
		Point{
			Label: start.Format("2006-01-02"),
			Value: 1,
		},
		Point{
			Label: end.Format("2006-01-02"),
			Value: 1,
		},
	}

	filled := fill(start.Unix(), end.Unix(), points)
	if len(filled) != 5 {
		t.Error("Length of filled points slice does not match expected length")
	}
}

func TestCalculatePointPercentages(t *testing.T) {
	points := []Point{
		Point{
			Label: "Foo",
			Value: 5,
		},
	}

	points = calculatePointPercentages(points, 100)

	if points[0].PercentageValue != 5.00 {
		t.Errorf("Percentage value should be 5.00, is %.2f", points[0].PercentageValue)
	}
}
