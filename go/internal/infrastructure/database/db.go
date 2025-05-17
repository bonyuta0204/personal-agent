package database

import (
	"fmt"

	"github.com/bonyuta0204/personal-agent/go/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// NewDBConnection creates a new database connection
func NewDBConnection(cfg *config.DatabaseConfig) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Name,
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// CloseDB closes the database connection
func CloseDB(db *sqlx.DB) error {
	if db != nil {
		return db.Close()
	}
	return nil
}
