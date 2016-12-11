package db

import(
  "database/sql"
  "log"
)


type Archive struct {
  ID int64
  Metric string
  Value string
  Count int64
  Date string
}

func (a *Archive) Save(conn *sql.DB) error {
  stmt, err := conn.Prepare(`INSERT INTO archive(
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
  stmt, err := Conn.Prepare(`
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

  Conn.Exec("START TRANSACTION")
  for rows.Next() {
    a := Archive{
      Metric: "pageviews",
      Value: "",
    }
    err = rows.Scan(&a.Count, &a.Date);
    checkError(err)
    a.Save(Conn)
  }
  Conn.Exec("COMMIT")
}

func CreateVisitorArchives() {

  /*
  SELECT
    COUNT(DISTINCT(visitor_id)) AS count, DATE_FORMAT(timestamp, ?) AS date_group
    FROM pageviews pv
    WHERE UNIX_TIMESTAMP(pv.timestamp) <= ? AND UNIX_TIMESTAMP(pv.timestamp) >= ?
    GROUP BY date_group
    */
  stmt, err := Conn.Prepare(`
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

  Conn.Exec("START TRANSACTION")
  for rows.Next() {
    a := Archive{
      Metric: "visitors",
      Value: "",
    }
    err = rows.Scan(&a.Count, &a.Date);
    checkError(err)
    a.Save(Conn)
  }
  Conn.Exec("COMMIT")
}

func checkError(err error) {
  if err != nil {
    log.Fatal(err)
  }
}
