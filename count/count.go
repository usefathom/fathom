package count

import(
  "database/sql"
  "log"
  "github.com/dannyvankooten/ana/db"
  "time"
)

type Archive struct {
  ID int64
  Metric string
  Value string
  Count int64
  Date string
}

type Point struct {
  Label string
  Value int
  PercentageValue float64
}

func (a *Archive) Save(Conn *sql.DB) error {
  stmt, err := db.Conn.Prepare(`INSERT INTO archive(
    metric,
    value,
    count,
    date
    ) VALUES( ?, ?, ?, ? )`)
  if err != nil {
      return err
  }
  defer stmt.Close()

  result, err := stmt.Exec(
    a.Metric,
    a.Value,
    a.Count,
    a.Date,
  )
  a.ID, _ = result.LastInsertId()

  return err
}


func CreateArchives() {
  CreatePageviewArchives()
  CreateVisitorArchives()
}

func CreatePageviewArchives() {
  stmt, err := db.Conn.Prepare(`
    SELECT
      COUNT(*) AS count,
      DATE_FORMAT(pv.timestamp, "%Y-%m-%d") AS date_group
    FROM pageviews pv
    WHERE NOT EXISTS(
      SELECT a.id
      FROM archive a
      WHERE a.metric = 'pageviews' AND a.date = DATE_FORMAT(pv.timestamp, "%Y-%m-%d")
    )
    GROUP BY date_group`)
  checkError(err)
  defer stmt.Close()

  rows, err := stmt.Query()
  checkError(err)
  defer rows.Close()

  db.Conn.Exec("START TRANSACTION")
  for rows.Next() {
    a := Archive{
      Metric: "pageviews",
      Value: "",
    }
    err = rows.Scan(&a.Count, &a.Date);
    checkError(err)
    a.Save(db.Conn)
  }
  db.Conn.Exec("COMMIT")
}

func CreateVisitorArchives() {
  stmt, err := db.Conn.Prepare(`
    SELECT
      COUNT(DISTINCT(pv.visitor_id)) AS count,
      DATE_FORMAT(pv.timestamp, "%Y-%m-%d") AS date_group
    FROM pageviews pv
    WHERE NOT EXISTS(
      SELECT a.id
      FROM archive a
      WHERE a.metric = 'visitors' AND a.date = DATE_FORMAT(pv.timestamp, "%Y-%m-%d")
    )
    GROUP BY date_group`)
  checkError(err)
  defer stmt.Close()

  rows, err := stmt.Query()
  checkError(err)
  defer rows.Close()

  db.Conn.Exec("START TRANSACTION")
  for rows.Next() {
    a := Archive{
      Metric: "visitors",
      Value: "",
    }
    err = rows.Scan(&a.Count, &a.Date);
    checkError(err)
    a.Save(db.Conn)
  }
  db.Conn.Exec("COMMIT")
}

func checkError(err error) {
  if err != nil {
    log.Fatal(err)
  }
}

func Custom(sql string, before int64, after int64, limit int, total float64) []Point {
  stmt, err := db.Conn.Prepare(sql)
  checkError(err)
  defer stmt.Close()

  rows, err := stmt.Query(before, after, limit)
  checkError(err)
  defer rows.Close()

  results := newPointSlice(rows, total)
  return results
}

func newPointSlice(rows *sql.Rows, total float64) []Point {
  results := make([]Point, 0)
  for rows.Next() {
    var d Point
    err := rows.Scan(&d.Label, &d.Value);
    checkError(err)

    d.PercentageValue = float64(d.Value) / total * 100
    results = append(results, d)
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
  newPoints := make([]Point, 0)
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
