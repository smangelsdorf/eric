package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func DefaultStoragePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("getting home directory: %w", err)
	}
	dir := filepath.Join(home, ".eric", "tasks")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("creating storage directory: %w", err)
	}
	return dir, nil
}

func TaskFilePath(storageDir, taskID string) string {
	return filepath.Join(storageDir, taskID+".md")
}

func WriteTaskFile(storageDir, taskID, summary, origin, destination, content string) (string, error) {
	filePath := TaskFilePath(storageDir, taskID)

	now := time.Now().UTC().Format(time.RFC3339)
	md := fmt.Sprintf("# %s\n\n"+
		"- **ID**: %s\n"+
		"- **Origin**: %s\n"+
		"- **Destination**: %s\n"+
		"- **Created**: %s\n\n"+
		"---\n\n"+
		"%s\n",
		summary, taskID, origin, destination, now, content,
	)

	if err := os.WriteFile(filePath, []byte(md), 0644); err != nil {
		return "", fmt.Errorf("writing task file: %w", err)
	}
	return filePath, nil
}

func AppendToTaskFile(filePath, content string) error {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("opening task file: %w", err)
	}
	defer f.Close()

	now := time.Now().UTC().Format(time.RFC3339)
	entry := fmt.Sprintf("\n---\n\n**Update** (%s):\n\n%s\n", now, content)

	if _, err := f.WriteString(entry); err != nil {
		return fmt.Errorf("appending to task file: %w", err)
	}
	return nil
}

func ReadTaskFile(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("reading task file: %w", err)
	}
	return string(data), nil
}

// FileContains checks whether a task file contains the given query (case-insensitive).
func FileContains(filePath, query string) (bool, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return false, fmt.Errorf("reading task file: %w", err)
	}
	return strings.Contains(strings.ToLower(string(data)), strings.ToLower(query)), nil
}
