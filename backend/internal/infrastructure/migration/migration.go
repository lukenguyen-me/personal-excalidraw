package migration

import (
	"database/sql"
	"embed"
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed sql/*.sql
var migrationFiles embed.FS

// Runner handles database migrations
type Runner struct {
	logger *slog.Logger
}

// NewRunner creates a new migration runner
func NewRunner(logger *slog.Logger) *Runner {
	return &Runner{
		logger: logger,
	}
}

// Run executes all pending migrations
func (r *Runner) Run(db *sql.DB, dbName string) error {
	r.logger.Info("Starting database migrations")

	// Create postgres driver
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	// Create source from embedded files
	sourceDriver, err := iofs.New(migrationFiles, "sql")
	if err != nil {
		return fmt.Errorf("failed to create migration source: %w", err)
	}

	// Create migrate instance
	m, err := migrate.NewWithInstance("iofs", sourceDriver, dbName, driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	// Get current version
	currentVersion, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("failed to get current version: %w", err)
	}

	if err == migrate.ErrNilVersion {
		r.logger.Info("Current migration version: None (no migrations applied)")
	} else {
		r.logger.Info("Current migration version",
			"version", currentVersion,
			"dirty", dirty,
		)
	}

	// Run migrations
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			r.logger.Info("No pending migrations to apply")
			return nil
		}
		return fmt.Errorf("migration failed: %w", err)
	}

	// Get new version
	newVersion, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("failed to get new version: %w", err)
	}

	r.logger.Info("Migrations applied successfully",
		"version", newVersion,
		"dirty", dirty,
	)

	return nil
}
