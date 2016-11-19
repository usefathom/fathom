package ana

import (
  "database/sql"
  "log"
  _"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func connectToDatabase() *sql.DB {
  // setup db connection
  var err error
  db, err = sql.Open("mysql", "root:root@/ana")
  if err != nil {
      log.Fatal(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
  }
  defer db.Close()

  // Open doesn't open a connection. Validate DSN data:
  err = db.Ping()
  if err != nil {
      log.Fatal(err.Error()) // proper error handling instead of panic in your app
  }
  
  return db
}
