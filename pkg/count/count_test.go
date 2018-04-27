package count

import (
	"testing"
	"time"

	"github.com/usefathom/fathom/pkg/models"
)

func TestCalculatePointPercentages(t *testing.T) {
	points := []*models.Point{
		&models.Point{
			Label: "Foo",
			Value: 5,
		},
	}

	points = calculatePointPercentages(points, 100)

	if points[0].PercentageValue != 5.00 {
		t.Errorf("Percentage value should be 5.00, is %.2f", points[0].PercentageValue)
	}
}
