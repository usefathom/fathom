package count

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/usefathom/fathom/pkg/datastore"
	"github.com/usefathom/fathom/pkg/models"
)

func getLastArchivedDate() string {
	value, _ := datastore.GetOption("last_archived")
	if value == "" {
		return time.Now().AddDate(-1, 0, 0).Format("2006-01-02")
	}

	return value
}

// Archive aggregates data into daily totals
func Archive() {
	start := time.Now()

	lastArchived := getLastArchivedDate()
	CreatePageviewTotals(lastArchived)
	CreateVisitorTotals(lastArchived)
	CreateScreenTotals(lastArchived)
	CreateLanguageTotals(lastArchived)
	CreateBrowserTotals(lastArchived)
	CreateReferrerTotals(lastArchived)
	datastore.SetOption("last_archived", time.Now().Format("2006-01-02"))

	end := time.Now()
	log.Infof("finished aggregating metrics. ran for %dms.", (end.UnixNano()-start.UnixNano())/1000000)
}

func calculatePointPercentages(points []*models.Point, total int) []*models.Point {
	// calculate percentage values for each point
	for _, p := range points {
		p.PercentageValue = float64(p.Value) / float64(total) * 100.00
	}

	return points
}
