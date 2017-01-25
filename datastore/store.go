package datastore

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql" // mysql driver
	_ "github.com/lib/pq"              // postgresql driver
)

// DB ...
var DB *sql.DB

// Init creates a database connection pool
func Init() *sql.DB {
	driver := os.Getenv("ANA_DATABASE_DRIVER")
	if driver == "" {
		driver = "mysql"
	}

	DB = New(driver, getDSN(driver))
	return DB
}

// New creates a new database pool
func New(driver string, config string) *sql.DB {
	db, err := sql.Open(driver, config)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	return db
}

func getDSN(driver string) string {
	var dsn = fmt.Sprintf(
		"%s:%s@%s/%s",
		os.Getenv("ANA_DATABASE_USER"),
		os.Getenv("ANA_DATABASE_PASSWORD"),
		os.Getenv("ANA_DATABASE_HOST"),
		os.Getenv("ANA_DATABASE_NAME"),
	)

	if driver == "postgres" {
		dsn = "postgres://" + dsn
	}

	return dsn
}
