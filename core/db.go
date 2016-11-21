package core

import (
  "database/sql"
  _"github.com/go-sql-driver/mysql"
  "log"
)

var DB *sql.DB

func SetupDatabaseConnection() *sql.DB {
  // setup db connection
  var err error
  DB, err = sql.Open("mysql", "root:root@/ana")
  if err != nil {
      log.Fatal(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
  }

  // Open doesn't open a connection. Validate DSN data:
  err = DB.Ping()
  if err != nil {
      log.Fatal(err.Error()) // proper error handling instead of panic in your app
  }

  return DB
}
