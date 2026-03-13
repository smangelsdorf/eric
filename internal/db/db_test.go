package db

import (
	"testing"
)

func TestMigrate(t *testing.T) {
	dbPath := t.TempDir() + "/test.db"
	database, err := Open(dbPath)
	if err != nil {
		t.Fatalf("opening database: %v", err)
	}
	defer database.Close()

	// Running Open again should be idempotent
	database2, err := Open(dbPath)
	if err != nil {
		t.Fatalf("opening database again: %v", err)
	}
	database2.Close()
}
