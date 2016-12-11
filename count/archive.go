package count

import(
  "database/sql"
  "log"
  "github.com/dannyvankooten/ana/db"
)

type Archive struct {
  ID int64
  Metric string
  Value string
  Count int64
  Date string
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
