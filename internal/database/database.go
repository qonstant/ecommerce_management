package database

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	"ecommerce_management/internal/config"
	"ecommerce_management/pkg/store"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var (
	DB             *sql.DB
	initOnce       sync.Once
	initDBInstance *dbSingleton
	instanceCreated bool
)

type dbSingleton struct{}

func InitDB() *dbSingleton {
	if !instanceCreated {
		initOnce.Do(func() {
			fmt.Println("Creating DB instance now.")
			var err error

			// Load configuration
			config, err := config.LoadConfig(".")
			if err != nil {
				log.Fatalf("Error loading config: %v", err)
			}

			// Print DB_SOURCE for debugging
			fmt.Println("DB_SOURCE:", config.DBSource)

			// Open a connection to the database using DB_SOURCE directly
			DB, err = sql.Open("postgres", config.DBSource)
			if err != nil {
				log.Fatalf("Error connecting to the database: %v", err)
			}

			// Check if the connection to the database is working
			err = DB.Ping()
			if err != nil {
				log.Fatalf("Error pinging the database: %v", err)
			}

			// Run database migrations
			if err := store.Migrate(config.DBSource); err != nil {
				log.Fatalf("Could not run database migrations: %v", err)
			}

			log.Println("Connected to the database successfully!")
			initDBInstance = &dbSingleton{}
			instanceCreated = true
		})
	} else {
		fmt.Println("DB instance already created.")
	}

	return initDBInstance
}
