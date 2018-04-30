package datastore

import (
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // mysql driver
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // postgresql driver
	"github.com/rubenv/sql-migrate"
	"log"
)

var dbx *sqlx.DB

// ErrNoResults is returned when a query yielded 0 results
var ErrNoResults = errors.New("datastore: query returned 0 results")

// Init creates a database connection pool (using sqlx)
func Init(driver string, host string, name string, user string, password string) *sqlx.DB {
	dbx = New(driver, getDSN(driver, host, name, user, password))

	// run migrations
	runMigrations(driver)

	return dbx
}

// New creates a new database pool
func New(driver string, dsn string) *sqlx.DB {
	dbx := sqlx.MustConnect(driver, dsn)
	return dbx
}

func getDSN(driver string, host string, name string, user string, password string) string {
	var dsn string

	switch driver {
	case "postgres":
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s", host, user, password, name)
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@%s/%s?parseTime=true", user, password, host, name)
	}

	return dsn
}

func runMigrations(driver string) {
	migrations := migrate.FileMigrationSource{
		Dir: "pkg/datastore/migrations", // TODO: Move to bindata
	}

	migrate.SetTable("migrations")

	n, err := migrate.Exec(dbx.DB, driver, migrations, migrate.Up)

	if err != nil {
		log.Fatal("Database migrations failed: ", err)
	}

	if n > 0 {
		log.Printf("Applied %d database migrations!\n", n)
	}
}
