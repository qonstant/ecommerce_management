package store

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrate(dataSourceName string) error {
	// Get the absolute path for the migrations directory
	dir, err := filepath.Abs("./db/migrations")
	if err != nil {
		return fmt.Errorf("could not get absolute path for migrations: %w", err)
	}

	// Check if the migrations directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Printf("Migrations directory '%s' does not exist, skipping migrations", dir)
		return nil // No migrations to run
	}

	// Print the contents of the migrations directory
	files, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("could not read migrations directory: %w", err)
	}

	log.Println("Migration files:")
	for _, file := range files {
		log.Println(file.Name())
	}

	// Perform database migrations using golang-migrate
	if !strings.Contains(dataSourceName, "://") {
		return errors.New("undefined data source name " + dataSourceName)
	}

	m, err := migrate.New(fmt.Sprintf("file://%s", dir), dataSourceName)
	if err != nil {
		return fmt.Errorf("error creating migration instance: %w", err)
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("error running migrations: %w", err)
	}

	log.Println("Database migrations ran successfully!")
	return nil
}
