package db

import (
	"database/sql"
	"fmt"
	"time"
)

type Task struct {
	ID          string
	Summary     string
	Origin      string
	Destination string
	Status      string
	FilePath    string
	CreatedAt   time.Time
}

func nextTaskID(tx *sql.Tx) (string, error) {
	var seq int
	if err := tx.QueryRow(`SELECT next_id FROM task_seq WHERE rowid = 1`).Scan(&seq); err != nil {
		return "", fmt.Errorf("reading sequence: %w", err)
	}
	if _, err := tx.Exec(`UPDATE task_seq SET next_id = ? WHERE rowid = 1`, seq+1); err != nil {
		return "", fmt.Errorf("updating sequence: %w", err)
	}
	return fmt.Sprintf("ERIC-%d", seq), nil
}

func CreateTask(db *sql.DB, summary, origin, destination, filePath string) (*Task, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback()

	id, err := nextTaskID(tx)
	if err != nil {
		return nil, err
	}

	var createdAt time.Time
	err = tx.QueryRow(
		`INSERT INTO tasks (id, summary, origin, destination, file_path)
		 VALUES (?, ?, ?, ?, ?)
		 RETURNING created_at`,
		id, summary, origin, destination, filePath,
	).Scan(&createdAt)
	if err != nil {
		return nil, fmt.Errorf("inserting task: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("committing transaction: %w", err)
	}

	return &Task{
		ID:          id,
		Summary:     summary,
		Origin:      origin,
		Destination: destination,
		Status:      "open",
		FilePath:    filePath,
		CreatedAt:   createdAt,
	}, nil
}

type TaskFilter struct {
	Origin      string
	Destination string
	Status      string
}

func ListTasks(db *sql.DB, filter TaskFilter) ([]Task, error) {
	query := `SELECT id, summary, origin, destination, status, file_path, created_at FROM tasks WHERE 1=1`
	var args []any

	if filter.Origin != "" {
		query += ` AND origin = ?`
		args = append(args, filter.Origin)
	}
	if filter.Destination != "" {
		query += ` AND destination = ?`
		args = append(args, filter.Destination)
	}
	if filter.Status != "" {
		query += ` AND status = ?`
		args = append(args, filter.Status)
	}
	query += ` ORDER BY created_at DESC`

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("listing tasks: %w", err)
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Summary, &t.Origin, &t.Destination, &t.Status, &t.FilePath, &t.CreatedAt); err != nil {
			return nil, fmt.Errorf("scanning task: %w", err)
		}
		tasks = append(tasks, t)
	}
	return tasks, rows.Err()
}

func GetTask(db *sql.DB, id string) (*Task, error) {
	var t Task
	err := db.QueryRow(
		`SELECT id, summary, origin, destination, status, file_path, created_at
		 FROM tasks WHERE id = ?`, id,
	).Scan(&t.ID, &t.Summary, &t.Origin, &t.Destination, &t.Status, &t.FilePath, &t.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("getting task: %w", err)
	}
	return &t, nil
}

func UpdateTaskStatus(db *sql.DB, id, status string) error {
	result, err := db.Exec(`UPDATE tasks SET status = ? WHERE id = ?`, status, id)
	if err != nil {
		return fmt.Errorf("updating task status: %w", err)
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return fmt.Errorf("task %s not found", id)
	}
	return nil
}
