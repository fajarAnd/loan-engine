package main

import (
	"flag"
	"github.com/fajar-andriansyah/loan-engine/config"
	database2 "github.com/fajar-andriansyah/loan-engine/infrastructure/database"
	"log"
	_ "os"
	_ "strconv"
)

func main() {
	// Load configuration
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Parse command line flags
	var (
		action  = flag.String("action", "up", "Migration action: up, down, force, version")
		version = flag.Int("version", -1, "Migration version (for force action)")
	)
	flag.Parse()

	// Connect to database
	dbConfig := database2.LoadDatabaseConfig()
	db, err := database2.Connect(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create migration instance with schema
	migration, err := database2.NewMigration(db, dbConfig.Schema)
	if err != nil {
		log.Fatalf("Failed to create migration: %v", err)
	}
	defer migration.Close()

	// Execute migration action
	switch *action {
	case "up":
		if err := migration.Up(); err != nil {
			log.Fatalf("Failed to run migrations up: %v", err)
		}
		log.Printf("Migrations up completed successfully for schema: %s", dbConfig.Schema)

	case "down":
		if err := migration.Down(); err != nil {
			log.Fatalf("Failed to run migrations down: %v", err)
		}
		log.Printf("Migrations down completed successfully for schema: %s", dbConfig.Schema)

	case "force":
		if *version == -1 {
			log.Fatal("Version flag is required for force action")
		}
		if err := migration.Force(*version); err != nil {
			log.Fatalf("Failed to force migration version: %v", err)
		}
		log.Printf("Migration version forced to %d successfully for schema: %s", *version, dbConfig.Schema)

	case "version":
		v, dirty, err := migration.Version()
		if err != nil {
			log.Fatalf("Failed to get migration version: %v", err)
		}
		log.Printf("Current migration version: %d (dirty: %t) for schema: %s", v, dirty, dbConfig.Schema)

	default:
		log.Fatalf("Unknown action: %s. Available actions: up, down, force, version", *action)
	}
}
