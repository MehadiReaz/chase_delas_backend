package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "modernc.org/sqlite"
)

var (
	// Connection pool settings
	MaxOpenConns    = 25
	MaxIdleConns    = 25
	ConnMaxLifetime = 5 * time.Minute
)

type DB struct {
	SQL *sql.DB
}

// Connect initializes and returns a database connection pool
func Connect() (*DB, error) {
	dbPath := getDBPath()

	sqlDB, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("database connection failed: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxOpenConns(MaxOpenConns)
	sqlDB.SetMaxIdleConns(MaxIdleConns)
	sqlDB.SetConnMaxLifetime(ConnMaxLifetime)

	// Verify connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	log.Println("Database connection established")
	return &DB{SQL: sqlDB}, nil
}

func (db *DB) Close() error {
	if err := db.SQL.Close(); err != nil {
		return fmt.Errorf("database closure failed: %w", err)
	}
	log.Println("Database connection closed")
	return nil
}

func getDBPath() string {
	if path := os.Getenv("DB_PATH"); path != "" {
		return path
	}
	return "./data.db"
}
