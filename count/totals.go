package count

import(
  "github.com/dannyvankooten/ana/db"
)

func TotalVisitors(before int64, after int64) float64 {
  // get total
  stmt, err := db.Conn.Prepare(`
    SELECT
    SUM(a.count)
    FROM archive a
    WHERE a.metric = 'visitors' AND UNIX_TIMESTAMP(a.date) <= ? AND UNIX_TIMESTAMP(a.date) >= ?`)
  checkError(err)
  defer stmt.Close()
  var total float64
  stmt.QueryRow(before, after).Scan(&total)
  return total
}
