package db

import (
	"testing"
)

func TestRegisterAndListProjects(t *testing.T) {
	dbPath := t.TempDir() + "/test.db"
	database, err := Open(dbPath)
	if err != nil {
		t.Fatalf("opening database: %v", err)
	}
	defer database.Close()

	if err := RegisterProject(database, "project-a", "First project"); err != nil {
		t.Fatalf("registering project: %v", err)
	}
	if err := RegisterProject(database, "project-b", "Second project"); err != nil {
		t.Fatalf("registering project: %v", err)
	}

	projects, err := ListProjects(database)
	if err != nil {
		t.Fatalf("listing projects: %v", err)
	}
	if len(projects) != 2 {
		t.Fatalf("expected 2 projects, got %d", len(projects))
	}
	if projects[0].Name != "project-a" || projects[1].Name != "project-b" {
		t.Fatalf("unexpected project order: %v", projects)
	}
}

func TestRegisterProjectUpsert(t *testing.T) {
	dbPath := t.TempDir() + "/test.db"
	database, err := Open(dbPath)
	if err != nil {
		t.Fatalf("opening database: %v", err)
	}
	defer database.Close()

	if err := RegisterProject(database, "proj", "Original"); err != nil {
		t.Fatalf("registering project: %v", err)
	}
	if err := RegisterProject(database, "proj", "Updated"); err != nil {
		t.Fatalf("re-registering project: %v", err)
	}

	projects, err := ListProjects(database)
	if err != nil {
		t.Fatalf("listing projects: %v", err)
	}
	if len(projects) != 1 {
		t.Fatalf("expected 1 project, got %d", len(projects))
	}
	if projects[0].Description != "Updated" {
		t.Fatalf("expected updated description, got %q", projects[0].Description)
	}
}

func TestProjectExists(t *testing.T) {
	dbPath := t.TempDir() + "/test.db"
	database, err := Open(dbPath)
	if err != nil {
		t.Fatalf("opening database: %v", err)
	}
	defer database.Close()

	exists, err := ProjectExists(database, "nope")
	if err != nil {
		t.Fatalf("checking existence: %v", err)
	}
	if exists {
		t.Fatal("expected project not to exist")
	}

	RegisterProject(database, "yep", "")
	exists, err = ProjectExists(database, "yep")
	if err != nil {
		t.Fatalf("checking existence: %v", err)
	}
	if !exists {
		t.Fatal("expected project to exist")
	}
}
