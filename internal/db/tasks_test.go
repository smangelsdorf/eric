package db

import (
	"testing"
)

func TestCreateAndGetTask(t *testing.T) {
	dbPath := t.TempDir() + "/test.db"
	database, err := Open(dbPath)
	if err != nil {
		t.Fatalf("opening database: %v", err)
	}
	defer database.Close()

	RegisterProject(database, "origin-proj", "Origin")
	RegisterProject(database, "dest-proj", "Destination")

	task, err := CreateTask(database, "Test task", "origin-proj", "dest-proj", "/tmp/test.md")
	if err != nil {
		t.Fatalf("creating task: %v", err)
	}
	if task.ID != "ERIC-1" {
		t.Fatalf("expected ERIC-1, got %s", task.ID)
	}
	if task.Status != "open" {
		t.Fatalf("expected open status, got %s", task.Status)
	}

	got, err := GetTask(database, "ERIC-1")
	if err != nil {
		t.Fatalf("getting task: %v", err)
	}
	if got == nil {
		t.Fatal("expected task, got nil")
	}
	if got.Summary != "Test task" {
		t.Fatalf("expected summary 'Test task', got %q", got.Summary)
	}
}

func TestTaskSequence(t *testing.T) {
	dbPath := t.TempDir() + "/test.db"
	database, err := Open(dbPath)
	if err != nil {
		t.Fatalf("opening database: %v", err)
	}
	defer database.Close()

	RegisterProject(database, "a", "")
	RegisterProject(database, "b", "")

	t1, _ := CreateTask(database, "First", "a", "b", "/tmp/1.md")
	t2, _ := CreateTask(database, "Second", "a", "b", "/tmp/2.md")

	if t1.ID != "ERIC-1" {
		t.Fatalf("expected ERIC-1, got %s", t1.ID)
	}
	if t2.ID != "ERIC-2" {
		t.Fatalf("expected ERIC-2, got %s", t2.ID)
	}
}

func TestListTasksWithFilter(t *testing.T) {
	dbPath := t.TempDir() + "/test.db"
	database, err := Open(dbPath)
	if err != nil {
		t.Fatalf("opening database: %v", err)
	}
	defer database.Close()

	RegisterProject(database, "a", "")
	RegisterProject(database, "b", "")
	RegisterProject(database, "c", "")

	CreateTask(database, "Task 1", "a", "b", "/tmp/1.md")
	CreateTask(database, "Task 2", "a", "c", "/tmp/2.md")
	CreateTask(database, "Task 3", "b", "c", "/tmp/3.md")

	// Filter by origin
	tasks, err := ListTasks(database, TaskFilter{Origin: "a"})
	if err != nil {
		t.Fatalf("listing tasks: %v", err)
	}
	if len(tasks) != 2 {
		t.Fatalf("expected 2 tasks from origin a, got %d", len(tasks))
	}

	// Filter by destination
	tasks, err = ListTasks(database, TaskFilter{Destination: "c"})
	if err != nil {
		t.Fatalf("listing tasks: %v", err)
	}
	if len(tasks) != 2 {
		t.Fatalf("expected 2 tasks to destination c, got %d", len(tasks))
	}

	// No filter
	tasks, err = ListTasks(database, TaskFilter{})
	if err != nil {
		t.Fatalf("listing tasks: %v", err)
	}
	if len(tasks) != 3 {
		t.Fatalf("expected 3 tasks, got %d", len(tasks))
	}
}

func TestUpdateTaskStatus(t *testing.T) {
	dbPath := t.TempDir() + "/test.db"
	database, err := Open(dbPath)
	if err != nil {
		t.Fatalf("opening database: %v", err)
	}
	defer database.Close()

	RegisterProject(database, "a", "")
	RegisterProject(database, "b", "")

	CreateTask(database, "Task", "a", "b", "/tmp/1.md")

	if err := UpdateTaskStatus(database, "ERIC-1", "closed"); err != nil {
		t.Fatalf("updating status: %v", err)
	}

	task, _ := GetTask(database, "ERIC-1")
	if task.Status != "closed" {
		t.Fatalf("expected closed, got %s", task.Status)
	}
}

func TestGetTaskNotFound(t *testing.T) {
	dbPath := t.TempDir() + "/test.db"
	database, err := Open(dbPath)
	if err != nil {
		t.Fatalf("opening database: %v", err)
	}
	defer database.Close()

	task, err := GetTask(database, "ERIC-999")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if task != nil {
		t.Fatal("expected nil for nonexistent task")
	}
}
