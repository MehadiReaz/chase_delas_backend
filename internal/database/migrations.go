package database

import (
	"fmt"
	"log"
	"os"
)

// RunMigrations executes database schema migrations
func (db *DB) RunMigrations() error {
	// Table migrations
	tableMigrations := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			role TEXT NOT NULL CHECK(role IN ('superadmin', 'admin', 'vendor', 'user')) DEFAULT 'user'
		);`,
		// ... other table migrations
	}

	// Data migrations
	dataMigrations := []string{
		fmt.Sprintf(`INSERT INTO users (id, email, password, role) 
		SELECT 'superadmin1', '%s', '%s', 'superadmin'
		WHERE NOT EXISTS (
			SELECT 1 FROM users WHERE email = '%s'
		);`,
			os.Getenv("SUPERADMIN_EMAIL"),
			os.Getenv("SUPERADMIN_PASSWORD_HASH"),
			os.Getenv("SUPERADMIN_EMAIL")),
	}

	// Run table migrations first
	for i, query := range tableMigrations {
		if _, err := db.SQL.Exec(query); err != nil {
			return fmt.Errorf("table migration %d failed: %w\nQuery: %s", i+1, err, query)
		}
	}

	// Then run data migrations
	for i, query := range dataMigrations {
		if _, err := db.SQL.Exec(query); err != nil {
			return fmt.Errorf("data migration %d failed: %w\nQuery: %s", i+1, err, query)
		}
	}

	log.Println("Database migrations completed successfully")
	return nil
}
