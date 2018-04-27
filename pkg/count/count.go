package count

import (
	"database/sql"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/usefathom/fathom/pkg/datastore"
)

// Point represents a data point, will always have a Label and Value
type Point struct {
	Label           string
	Value           int
	PercentageValue float64
}

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

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func newPointSlice(rows *sql.Rows) []Point {
	results := make([]Point, 0)

	// append point slices
	for rows.Next() {
		var d Point
		err := rows.Scan(&d.Label, &d.Value)
		checkError(err)
		results = append(results, d)
	}

	return results
}

func calculatePointPercentages(points []Point, total int) []Point {
	// calculate percentage values for each point
	for i, d := range points {
		points[i].PercentageValue = float64(d.Value) / float64(total) * 100
	}

	return points
}
