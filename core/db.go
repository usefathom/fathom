package core

import (
  "database/sql"
  _"github.com/go-sql-driver/mysql"
  "log"
  "os"
  "fmt"
)

var DB *sql.DB

func SetupDatabaseConnection() *sql.DB {
  var err error
  var dataSourceName = fmt.Sprintf("%s:%s@%s/%s", os.Getenv("ANA_DATABASE_USER"), os.Getenv("ANA_DATABASE_PASSWORD"), os.Getenv("ANA_DATABASE_HOST"), os.Getenv("ANA_DATABASE_NAME"))

  DB, err = sql.Open("mysql", dataSourceName)
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
