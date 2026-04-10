package db

import (
	"fmt"
	"testing"
)

func replyPathFn(taskID string, seq int) string {
	return fmt.Sprintf("/tmp/%s-%d.md", taskID, seq)
}

func TestCreateAndListReplies(t *testing.T) {
	dbPath := t.TempDir() + "/test.db"
	database, err := Open(dbPath)
	if err != nil {
		t.Fatalf("opening database: %v", err)
	}
	defer database.Close()

	RegisterProject(database, "a", "")
	RegisterProject(database, "b", "")
	CreateTask(database, "Task", "a", "b", staticPath("/tmp/1.md"))

	r1, err := CreateReply(database, "ERIC-1", "b", "a", replyPathFn)
	if err != nil {
		t.Fatalf("creating reply 1: %v", err)
	}
	if r1.Seq != 1 {
		t.Fatalf("expected seq 1, got %d", r1.Seq)
	}
	if r1.FilePath != "/tmp/ERIC-1-1.md" {
		t.Fatalf("expected /tmp/ERIC-1-1.md, got %s", r1.FilePath)
	}

	r2, err := CreateReply(database, "ERIC-1", "a", "b", replyPathFn)
	if err != nil {
		t.Fatalf("creating reply 2: %v", err)
	}
	if r2.Seq != 2 {
		t.Fatalf("expected seq 2, got %d", r2.Seq)
	}

	replies, err := ListReplies(database, "ERIC-1")
	if err != nil {
		t.Fatalf("listing replies: %v", err)
	}
	if len(replies) != 2 {
		t.Fatalf("expected 2 replies, got %d", len(replies))
	}
	if replies[0].Origin != "b" || replies[1].Origin != "a" {
		t.Fatalf("unexpected reply origins: %s, %s", replies[0].Origin, replies[1].Origin)
	}
}

func TestReplySequencePerTask(t *testing.T) {
	dbPath := t.TempDir() + "/test.db"
	database, err := Open(dbPath)
	if err != nil {
		t.Fatalf("opening database: %v", err)
	}
	defer database.Close()

	RegisterProject(database, "a", "")
	RegisterProject(database, "b", "")
	CreateTask(database, "Task 1", "a", "b", staticPath("/tmp/1.md"))
	CreateTask(database, "Task 2", "a", "b", staticPath("/tmp/2.md"))

	r1, _ := CreateReply(database, "ERIC-1", "b", "a", replyPathFn)
	r2, _ := CreateReply(database, "ERIC-2", "b", "a", replyPathFn)
	r3, _ := CreateReply(database, "ERIC-1", "a", "b", replyPathFn)

	if r1.Seq != 1 {
		t.Fatalf("ERIC-1 reply 1: expected seq 1, got %d", r1.Seq)
	}
	if r2.Seq != 1 {
		t.Fatalf("ERIC-2 reply 1: expected seq 1, got %d", r2.Seq)
	}
	if r3.Seq != 2 {
		t.Fatalf("ERIC-1 reply 2: expected seq 2, got %d", r3.Seq)
	}
}

func TestListRepliesEmpty(t *testing.T) {
	dbPath := t.TempDir() + "/test.db"
	database, err := Open(dbPath)
	if err != nil {
		t.Fatalf("opening database: %v", err)
	}
	defer database.Close()

	RegisterProject(database, "a", "")
	RegisterProject(database, "b", "")
	CreateTask(database, "Task", "a", "b", staticPath("/tmp/1.md"))

	replies, err := ListReplies(database, "ERIC-1")
	if err != nil {
		t.Fatalf("listing replies: %v", err)
	}
	if len(replies) != 0 {
		t.Fatalf("expected 0 replies, got %d", len(replies))
	}
}
