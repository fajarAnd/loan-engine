package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
	Host         string
	Port         int
	User         string
	Password     string
	DBName       string
	Schema       string // Added schema field
	SSLMode      string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  time.Duration
}

func LoadDatabaseConfig() *Config {
	return &Config{
		Host:         viper.GetString("database.host"),
		Port:         viper.GetInt("database.port"),
		User:         viper.GetString("database.user"),
		Password:     viper.GetString("database.password"),
		DBName:       viper.GetString("database.name"),
		Schema:       viper.GetString("database.schema"),
		SSLMode:      viper.GetString("database.sslmode"),
		MaxOpenConns: viper.GetInt("database.max_open_conns"),
		MaxIdleConns: viper.GetInt("database.max_idle_conns"),
		MaxLifetime:  viper.GetDuration("database.max_lifetime"),
	}
}

func (c *Config) PostgresDSN() string {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)

	// Add search_path if schema is specified
	if c.Schema != "" {
		dsn += fmt.Sprintf(" search_path=%s,public", c.Schema)
	}

	return dsn
}

func Connect(config *Config) (*sql.DB, error) {
	log.Info().
		Str("host", config.Host).
		Int("port", config.Port).
		Str("user", config.User).
		Str("dbname", config.DBName).
		Str("schema", config.Schema).
		Str("sslmode", config.SSLMode).
		Msg("Connecting to PostgreSQL database")

	db, err := sql.Open("postgres", config.PostgresDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Configure connection pool
	if config.MaxOpenConns > 0 {
		db.SetMaxOpenConns(config.MaxOpenConns)
	}
	if config.MaxIdleConns > 0 {
		db.SetMaxIdleConns(config.MaxIdleConns)
	}
	if config.MaxLifetime > 0 {
		db.SetConnMaxLifetime(config.MaxLifetime)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set search path for the session if schema is specified
	if config.Schema != "" {
		_, err := db.Exec(fmt.Sprintf("SET search_path TO %s, public", config.Schema))
		if err != nil {
			return nil, fmt.Errorf("failed to set search path: %w", err)
		}
		log.Info().Str("schema", config.Schema).Msg("Database search path set")
	}

	log.Info().Msg("Successfully connected to PostgreSQL database")
	return db, nil
}
