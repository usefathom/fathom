package sqlstore

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/go-sql-driver/mysql" // mysql driver
	"github.com/gobuffalo/packr/v2"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"           // postgresql driver
	_ "github.com/mattn/go-sqlite3" //sqlite3 driver
	migrate "github.com/rubenv/sql-migrate"
	log "github.com/sirupsen/logrus"
)

const (
	MYSQL    = "mysql"
	POSTGRES = "postgres"
	SQLITE   = "sqlite3"

	DATE_FORMAT = "2006-01-02 15:00:00"
)

type sqlstore struct {
	*sqlx.DB

	Driver string
	Config *Config
}

// ErrNoResults is returned when a query yielded 0 results
var ErrNoResults = errors.New("datastore: query returned 0 results")

// New creates a new database pool
func New(c *Config) *sqlstore {
	dsn := c.DSN()
	dbx, err := sqlx.Connect(c.Driver, dsn)
	if err != nil {
		log.Fatalf("Error connecting to database: %s", err)
	}
	db := &sqlstore{dbx, c.Driver, c}

	if c.Host == "" || c.Driver == SQLITE {
		log.Printf("Connected to %s database: %s", c.Driver, c.Dbname())
	} else {
		log.Printf("Connected to %s database: %s on %s", c.Driver, c.Dbname(), c.Host)
	}

	// apply database migrations (if any)
	db.Migrate()

	return db
}

func (db *sqlstore) Migrate() {
	migrationSource := &migrate.PackrMigrationSource{
		Box: packr.NewBox("./migrations"),
		Dir: db.Config.Driver,
	}
	migrate.SetTable("migrations")

	migrations, err := migrationSource.FindMigrations()
	if err != nil {
		log.Errorf("Error loading database migrations: %s", err)
	}

	if len(migrations) == 0 {
		log.Fatalf("Missing database migrations")
	}

	n, err := migrate.Exec(db.DB.DB, db.Config.Driver, migrationSource, migrate.Up)
	if err != nil {
		log.Errorf("Error applying database migrations: %s", err)
	}

	if n > 0 {
		log.Infof("Applied %d database migrations!", n)
	}
}

// Health check health of database
func (db *sqlstore) Health() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	return db.PingContext(ctx)
}

// Closes the db pool
func (db *sqlstore) Close() error {
	return db.DB.Close()
}

func mapError(err error) error {
	if err == sql.ErrNoRows {
		return ErrNoResults
	}

	return nil
}
