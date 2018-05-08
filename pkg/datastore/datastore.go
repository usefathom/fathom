package datastore

import (
	"errors"

	_ "github.com/go-sql-driver/mysql" // mysql driver
	"github.com/gobuffalo/packr"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"           // postgresql driver
	_ "github.com/mattn/go-sqlite3" //sqlite3 driver
	migrate "github.com/rubenv/sql-migrate"
	log "github.com/sirupsen/logrus"
)

var dbx *sqlx.DB

// ErrNoResults is returned when a query yielded 0 results
var ErrNoResults = errors.New("datastore: query returned 0 results")

// Init creates a database connection pool (using sqlx)
func Init(c *Config) *sqlx.DB {
	dbx = New(c)

	// run migrations
	runMigrations(c.Driver)

	return dbx
}

// New creates a new database pool
func New(c *Config) *sqlx.DB {
	dbx := sqlx.MustConnect(c.Driver, c.DSN())
	return dbx
}

// TODO: Move to command (but still auto-run on boot).
func runMigrations(driver string) {
	migrations := &migrate.PackrMigrationSource{
		Box: packr.NewBox("./migrations"),
		Dir: "./" + driver,
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
