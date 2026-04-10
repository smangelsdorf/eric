package db

import (
	"database/sql"
	"fmt"
	"time"
)

type Reply struct {
	TaskID      string
	Seq         int
	Origin      string
	Destination string
	FilePath    string
	CreatedAt   time.Time
}

func CreateReply(database *sql.DB, taskID, origin, destination string, filePathFn func(taskID string, seq int) string) (*Reply, error) {
	tx, err := database.Begin()
	if err != nil {
		return nil, fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback()

	var seq int
	var createdAt time.Time
	err = tx.QueryRow(`
		INSERT INTO replies (task_id, seq, origin, destination, file_path)
		SELECT ?, COALESCE(MAX(seq), 0) + 1, ?, ?, ''
		FROM replies WHERE task_id = ?
		RETURNING seq, created_at`,
		taskID, origin, destination, taskID,
	).Scan(&seq, &createdAt)
	if err != nil {
		return nil, fmt.Errorf("inserting reply: %w", err)
	}

	filePath := filePathFn(taskID, seq)
	if _, err := tx.Exec(`UPDATE replies SET file_path = ? WHERE task_id = ? AND seq = ?`, filePath, taskID, seq); err != nil {
		return nil, fmt.Errorf("updating reply file path: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("committing transaction: %w", err)
	}

	return &Reply{
		TaskID:      taskID,
		Seq:         seq,
		Origin:      origin,
		Destination: destination,
		FilePath:    filePath,
		CreatedAt:   createdAt,
	}, nil
}

func ListReplies(database *sql.DB, taskID string) ([]Reply, error) {
	rows, err := database.Query(
		`SELECT task_id, seq, origin, destination, file_path, created_at
		 FROM replies WHERE task_id = ? ORDER BY seq`, taskID,
	)
	if err != nil {
		return nil, fmt.Errorf("listing replies: %w", err)
	}
	defer rows.Close()

	var replies []Reply
	for rows.Next() {
		var r Reply
		if err := rows.Scan(&r.TaskID, &r.Seq, &r.Origin, &r.Destination, &r.FilePath, &r.CreatedAt); err != nil {
			return nil, fmt.Errorf("scanning reply: %w", err)
		}
		replies = append(replies, r)
	}
	return replies, rows.Err()
}
