package db

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	//"Import this package, but I'm not using it directly" blank identifier
)

func RunMigrations(dbURL string) error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	// Look in src/db/migrations
	migrationsPath := filepath.Join(wd, "db", "migrations")
	migrationsURL := fmt.Sprintf("file://%s", migrationsPath)

	m, err := migrate.New(
		migrationsURL,
		dbURL,
	)

	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("failed to get migration version: %w", err)
	}

	if err == migrate.ErrNilVersion {
		log.Println("No migrations have been applied yet")
	} else {
		log.Printf("Database is at migration version: %d (dirty: %v)", version, dirty)
	}

	log.Println("Migrations completed successfully")
	return nil
}
