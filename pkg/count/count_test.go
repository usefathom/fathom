package count

import (
	"testing"

	"github.com/usefathom/fathom/pkg/models"
)

func TestCalculatePointPercentages(t *testing.T) {
	totals := []*models.Total{
		&models.Total{
			Value: "Foo",
			Count: 5,
		},
	}

	totals = calculatePercentagesOfTotal(totals, 100)

	if totals[0].PercentageOfTotal != 5.00 {
		t.Errorf("Percentage value should be 5.00, is %.2f", totals[0].PercentageOfTotal)
	}
}
