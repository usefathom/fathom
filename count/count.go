package count

import (
	"database/sql"
	"log"
	"time"

	"github.com/dannyvankooten/ana/db"
	"github.com/dannyvankooten/ana/options"
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
	value := options.Get("last_archived")
	return value
}

// Archive aggregates data into daily totals
func Archive() {
	lastArchived := getLastArchivedDate()

	CreateVisitorTotals(lastArchived)
	CreatePageviewTotals(lastArchived)
	CreateScreenTotals(lastArchived)
	CreateLanguageTotals(lastArchived)
	CreateBrowserTotals(lastArchived)
	CreateReferrerTotals(lastArchived)

	err := options.Set("last_archived", time.Now().Format("2006-01-02"))
	checkError(err)
}

// Save the Total in the given database connection + table
func (t *Total) Save(Conn *sql.DB, table string) error {
	stmt, err := db.Conn.Prepare(`INSERT INTO ` + table + `(
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
	total := 0

	// append point slices
	for rows.Next() {
		var d Point
		err := rows.Scan(&d.Label, &d.Value)
		checkError(err)
		results = append(results, d)

		// sum total
		total += d.Value
	}

	// calculate percentage values for each point
	// TODO: This doesn't work when a limit on the # of points is given
	for i, d := range results {
		results[i].PercentageValue = float64(d.Value) / float64(total) * 100
	}

	return results
}

func fill(start int64, end int64, points []Point) []Point {
	// be smart about received timestamps
	if start > end {
		tmp := end
		end = start
		start = tmp
	}

	startTime := time.Unix(start, 0)
	endTime := time.Unix(end, 0)
	var newPoints []Point
	step := time.Hour * 24

	for startTime.Before(endTime) || startTime.Equal(endTime) {
		point := Point{
			Value: 0,
			Label: startTime.Format("2006-01-02"),
		}

		for j, p := range points {
			if p.Label == point.Label || p.Label == startTime.Format("2006-01") {
				point.Value = p.Value
				points[j] = points[len(points)-1]
				break
			}
		}

		newPoints = append(newPoints, point)
		startTime = startTime.Add(step)
	}

	return newPoints
}

func queryTotalRows(sql string, lastArchived string) *sql.Rows {
	stmt, err := db.Conn.Prepare(sql)
	checkError(err)
	defer stmt.Close()

	rows, err := stmt.Query(lastArchived)
	checkError(err)
	return rows
}

func processTotalRows(rows *sql.Rows, table string) {
	db.Conn.Exec("START TRANSACTION")
	for rows.Next() {
		var t Total
		err := rows.Scan(&t.Value, &t.Count, &t.CountUnique, &t.Date)
		checkError(err)
		t.Save(db.Conn, table)
	}
	db.Conn.Exec("COMMIT")

	rows.Close()
}
