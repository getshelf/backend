package postgres

import (
	"database/sql"
	"embed"

	"github.com/golang-migrate/migrate/v4"
	databasePostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func RunMigrations(db *sql.DB) error {
	driver, err := iofs.New(migrationFiles, "migrations")
	if err != nil {
		return err
	}
	databaseDriver, err := databasePostgres.WithInstance(db, &databasePostgres.Config{})
	if err != nil {
		return err
	}
	migrator, err := migrate.NewWithInstance("iofs", driver, "postgres", databaseDriver)
	if err != nil {
		return err
	}
	defer migrator.Close()
	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
