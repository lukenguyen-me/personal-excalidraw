package main

import (
	"database/sql"
	"embed"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

//go:embed *.sql
var migrationsFS embed.FS

func main() {
	// Define flags
	up := flag.Bool("up", false, "Run all up migrations")
	down := flag.Bool("down", false, "Run all down migrations")
	steps := flag.Int("steps", 0, "Run N migration steps (positive for up, negative for down)")
	version := flag.Uint("version", 0, "Migrate to specific version")
	force := flag.Int("force", -1, "Force set migration version (use with caution)")

	flag.Parse()

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Build database connection string
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "personal_excalidraw")
	dbSSLMode := getEnv("DB_SSL_MODE", "disable")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode,
	)

	// Open database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	// Create postgres driver
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal("Failed to create postgres driver:", err)
	}

	// Create source from embedded files
	source, err := iofs.New(migrationsFS, ".")
	if err != nil {
		log.Fatal("Failed to create migration source:", err)
	}

	// Create migrate instance
	m, err := migrate.NewWithInstance("iofs", source, "postgres", driver)
	if err != nil {
		log.Fatal("Failed to create migrate instance:", err)
	}

	// Get current version
	currentVersion, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		log.Fatal("Failed to get current version:", err)
	}

	if err == migrate.ErrNilVersion {
		log.Println("Current version: None (no migrations applied)")
	} else {
		log.Printf("Current version: %d (dirty: %v)\n", currentVersion, dirty)
	}

	// Execute migration command
	switch {
	case *force >= 0:
		log.Printf("Forcing version to %d...\n", *force)
		if err := m.Force(*force); err != nil {
			log.Fatal("Force failed:", err)
		}
		log.Println("Version forced successfully")

	case *up:
		log.Println("Running all up migrations...")
		if err := m.Up(); err != nil {
			if err == migrate.ErrNoChange {
				log.Println("No migrations to apply")
			} else {
				log.Fatal("Up migration failed:", err)
			}
		} else {
			log.Println("Migrations applied successfully")
		}

	case *down:
		log.Println("Running all down migrations...")
		if err := m.Down(); err != nil {
			if err == migrate.ErrNoChange {
				log.Println("No migrations to revert")
			} else {
				log.Fatal("Down migration failed:", err)
			}
		} else {
			log.Println("Migrations reverted successfully")
		}

	case *steps != 0:
		log.Printf("Running %d migration steps...\n", *steps)
		if err := m.Steps(*steps); err != nil {
			if err == migrate.ErrNoChange {
				log.Println("No migrations to apply")
			} else {
				log.Fatal("Steps migration failed:", err)
			}
		} else {
			log.Println("Migration steps completed successfully")
		}

	case *version != 0:
		log.Printf("Migrating to version %d...\n", *version)
		if err := m.Migrate(*version); err != nil {
			if err == migrate.ErrNoChange {
				log.Println("Already at target version")
			} else {
				log.Fatal("Migrate to version failed:", err)
			}
		} else {
			log.Println("Migrated to version successfully")
		}

	default:
		fmt.Println("Migration tool for personal-excalidraw")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  go run migrations/migrate.go [flags]")
		fmt.Println()
		fmt.Println("Flags:")
		flag.PrintDefaults()
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  go run migrations/migrate.go -up          # Run all pending migrations")
		fmt.Println("  go run migrations/migrate.go -down        # Rollback all migrations")
		fmt.Println("  go run migrations/migrate.go -steps 1     # Apply 1 migration")
		fmt.Println("  go run migrations/migrate.go -steps -1    # Rollback 1 migration")
		fmt.Println("  go run migrations/migrate.go -version 2   # Migrate to version 2")
		fmt.Println("  go run migrations/migrate.go -force 1     # Force set version to 1")
		os.Exit(0)
	}

	// Get new version
	newVersion, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		log.Fatal("Failed to get new version:", err)
	}

	if err == migrate.ErrNilVersion {
		log.Println("New version: None")
	} else {
		log.Printf("New version: %d (dirty: %v)\n", newVersion, dirty)
	}
}

// getEnv retrieves an environment variable with a default fallback
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
