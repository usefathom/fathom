package count

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
	"time"

	"github.com/usefathom/fathom/pkg/datastore"
)

// Total represents a daily aggregated total for a metric
type Total struct {
	ID          int64
	PageID      int64
	Value       string
	Count       int64
	CountUnique int64
	Date        string
}

// Point represents a data point, will always have a Label and Value
type Point struct {
	Label           string
	Value           int
	PercentageValue float64
}

func getLastArchivedDate() string {
	value, _ := datastore.GetOption("last_archived")
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

// Save the Total in the given database connection + table
func (t *Total) Save(Conn *sql.DB, table string) error {
	stmt, err := datastore.DB.Prepare(`INSERT INTO ` + table + `(
    value,
    count,
		count_unique,
    date
    ) VALUES( ?, ?, ?, ? ) ON DUPLICATE KEY UPDATE count = ?, count_unique = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(
		t.Value,
		t.Count,
		t.CountUnique,
		t.Date,
		t.Count,
		t.CountUnique,
	)
	t.ID, _ = result.LastInsertId()

	return err
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

func queryTotalRows(sql string, lastArchived string) *sql.Rows {
	stmt, err := datastore.DB.Prepare(sql)
	checkError(err)
	defer stmt.Close()

	rows, err := stmt.Query(lastArchived)
	checkError(err)
	return rows
}

func processTotalRows(rows *sql.Rows, table string) {
	datastore.DB.Exec("START TRANSACTION")
	for rows.Next() {
		var t Total
		err := rows.Scan(&t.Value, &t.Count, &t.CountUnique, &t.Date)
		checkError(err)
		t.Save(datastore.DB, table)
	}
	datastore.DB.Exec("COMMIT")

	rows.Close()
}
