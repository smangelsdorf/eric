package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func DefaultDBPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("getting home directory: %w", err)
	}
	dir := filepath.Join(home, ".eric")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("creating data directory: %w", err)
	}
	return filepath.Join(dir, "eric.db"), nil
}

func Open(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}
	if _, err := db.Exec(`PRAGMA journal_mode=WAL`); err != nil {
		db.Close()
		return nil, fmt.Errorf("setting journal mode: %w", err)
	}
	if _, err := db.Exec(`PRAGMA foreign_keys = ON`); err != nil {
		db.Close()
		return nil, fmt.Errorf("enabling foreign keys: %w", err)
	}
	if err := migrate(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("running migrations: %w", err)
	}
	return db, nil
}

func migrate(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS projects (
			name TEXT PRIMARY KEY,
			description TEXT NOT NULL DEFAULT ''
		);

		CREATE TABLE IF NOT EXISTS tasks (
			id TEXT PRIMARY KEY,
			summary TEXT NOT NULL,
			origin TEXT NOT NULL REFERENCES projects(name),
			destination TEXT NOT NULL REFERENCES projects(name),
			status TEXT NOT NULL DEFAULT 'open',
			file_path TEXT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS task_seq (
			next_id INTEGER NOT NULL
		);

		INSERT OR IGNORE INTO task_seq (rowid, next_id) VALUES (1, 1);
	`)
	return err
}
