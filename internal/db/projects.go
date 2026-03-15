package db

import (
	"database/sql"
	"fmt"
)

type Project struct {
	Name        string
	Description string
}

// RegisterProject creates or updates a project. Returns true if the project
// already existed (i.e. was updated rather than created).
func RegisterProject(db *sql.DB, name, description string) (updated bool, err error) {
	var exists bool
	if err := db.QueryRow(`SELECT COUNT(*) > 0 FROM projects WHERE name = ?`, name).Scan(&exists); err != nil {
		return false, fmt.Errorf("checking project: %w", err)
	}

	_, err = db.Exec(
		`INSERT INTO projects (name, description) VALUES (?, ?)
		 ON CONFLICT(name) DO UPDATE SET description = excluded.description`,
		name, description,
	)
	if err != nil {
		return false, fmt.Errorf("registering project: %w", err)
	}
	return exists, nil
}

func ListProjects(db *sql.DB) ([]Project, error) {
	rows, err := db.Query(`SELECT name, description FROM projects ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("listing projects: %w", err)
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var p Project
		if err := rows.Scan(&p.Name, &p.Description); err != nil {
			return nil, fmt.Errorf("scanning project: %w", err)
		}
		projects = append(projects, p)
	}
	return projects, rows.Err()
}

func ProjectExists(db *sql.DB, name string) (bool, error) {
	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM projects WHERE name = ?`, name).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("checking project: %w", err)
	}
	return count > 0, nil
}
