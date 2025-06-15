package database

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

type Migration struct {
	db       *sql.DB
	migrator *migrate.Migrate
	schema   string
}

func NewMigration(db *sql.DB, schema string) (*Migration, error) {
	config := &postgres.Config{}

	// Set schema for migrations if specified
	if schema != "" {
		config.SchemaName = schema
	}

	driver, err := postgres.WithInstance(db, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrator: %w", err)
	}

	return &Migration{
		db:       db,
		migrator: m,
		schema:   schema,
	}, nil
}

func (m *Migration) Up() error {
	log.Info().
		Str("schema", m.schema).
		Msg("Running database migrations up...")

	err := m.migrator.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations up: %w", err)
	}

	if err == migrate.ErrNoChange {
		log.Info().Msg("No new migrations to apply")
	} else {
		log.Info().Msg("Database migrations completed successfully")
	}

	return nil
}

func (m *Migration) Down() error {
	log.Info().
		Str("schema", m.schema).
		Msg("Running database migrations down...")

	err := m.migrator.Down()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations down: %w", err)
	}

	log.Info().Msg("Database migrations down completed")
	return nil
}

func (m *Migration) Force(version int) error {
	log.Info().
		Int("version", version).
		Str("schema", m.schema).
		Msg("Forcing migration version...")

	err := m.migrator.Force(version)
	if err != nil {
		return fmt.Errorf("failed to force migration version: %w", err)
	}

	log.Info().Int("version", version).Msg("Migration version forced successfully")
	return nil
}

func (m *Migration) Version() (uint, bool, error) {
	version, dirty, err := m.migrator.Version()
	if err != nil {
		return 0, false, fmt.Errorf("failed to get migration version: %w", err)
	}

	log.Info().
		Uint("version", version).
		Bool("dirty", dirty).
		Str("schema", m.schema).
		Msg("Current migration version")
	return version, dirty, nil
}

func (m *Migration) Close() error {
	sourceErr, dbErr := m.migrator.Close()
	if sourceErr != nil {
		return sourceErr
	}
	return dbErr
}
