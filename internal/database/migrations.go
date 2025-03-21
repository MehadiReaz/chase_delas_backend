package database

import (
	"fmt"
	"log"
)

// RunMigrations executes database schema migrations
func (db *DB) RunMigrations() error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		`CREATE TABLE IF NOT EXISTS books (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			author TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	for i, query := range migrations {
		if _, err := db.SQL.Exec(query); err != nil {
			return fmt.Errorf("migration %d failed: %w\nQuery: %s", i+1, err, query)
		}
	}

	log.Println("Database migrations completed successfully")
	return nil
}
