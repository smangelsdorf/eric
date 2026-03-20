package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/smangelsdorf/eric/internal/db"
	"github.com/smangelsdorf/eric/internal/storage"
)

type server struct {
	db         *sql.DB
	storageDir string
}

func main() {
	dbPath, err := db.DefaultDBPath()
	if err != nil {
		log.Fatalf("Failed to determine database path: %v", err)
	}

	database, err := db.Open(dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer database.Close()

	storageDir, err := storage.DefaultStoragePath()
	if err != nil {
		log.Fatalf("Failed to determine storage path: %v", err)
	}

	s := &server{db: database, storageDir: storageDir}

	mcpServer := mcp.NewServer(
		&mcp.Implementation{Name: "eric", Version: "0.1.0"},
		&mcp.ServerOptions{
			Instructions: "Eric is a private, local task management system for the developer. " +
				"ERIC- task IDs and task content must NEVER appear in git commits, " +
				"PR descriptions, code comments, or any externally visible content. " +
				"Use Eric only for internal task tracking between projects.",
		},
	)

	s.registerTools(mcpServer)

	if err := mcpServer.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func (s *server) registerTools(mcpServer *mcp.Server) {
	type registerProjectArgs struct {
		Name        string `json:"name" jsonschema:"Name of the project to register"`
		Description string `json:"description" jsonschema:"Brief description of the project"`
	}
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "register_project",
		Description: "Register a new project with Eric. Projects must be registered before they can be used as task origins or destinations.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args registerProjectArgs) (*mcp.CallToolResult, any, error) {
		if strings.TrimSpace(args.Name) == "" {
			return toolError(fmt.Errorf("project name is required")), nil, nil
		}
		updated, err := db.RegisterProject(s.db, args.Name, args.Description)
		if err != nil {
			return toolError(err), nil, nil
		}
		if updated {
			return toolText(fmt.Sprintf("Project %q updated.", args.Name)), nil, nil
		}
		return toolText(fmt.Sprintf("Project %q registered.", args.Name)), nil, nil
	})

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "list_projects",
		Description: "List all projects registered with Eric.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args struct{}) (*mcp.CallToolResult, any, error) {
		projects, err := db.ListProjects(s.db)
		if err != nil {
			return toolError(err), nil, nil
		}
		if len(projects) == 0 {
			return toolText("No projects registered."), nil, nil
		}
		text := ""
		for _, p := range projects {
			text += fmt.Sprintf("- **%s**: %s\n", p.Name, p.Description)
		}
		return toolText(text), nil, nil
	})

	type createTaskArgs struct {
		Summary     string `json:"summary" jsonschema:"Brief one-line summary of the task"`
		Origin      string `json:"origin" jsonschema:"Name of the project creating this task"`
		Destination string `json:"destination" jsonschema:"Name of the project this task is for"`
		Content     string `json:"content" jsonschema:"Detailed task description in Markdown"`
	}
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "create_task",
		Description: "Create a new task, passing work from one project to another.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args createTaskArgs) (*mcp.CallToolResult, any, error) {
		if strings.TrimSpace(args.Summary) == "" {
			return toolError(fmt.Errorf("task summary is required")), nil, nil
		}
		if strings.TrimSpace(args.Content) == "" {
			return toolError(fmt.Errorf("task content is required")), nil, nil
		}
		for _, proj := range []string{args.Origin, args.Destination} {
			exists, err := db.ProjectExists(s.db, proj)
			if err != nil {
				return toolError(err), nil, nil
			}
			if !exists {
				return toolError(fmt.Errorf("project %q is not registered — use register_project first", proj)), nil, nil
			}
		}

		task, err := db.CreateTask(s.db, args.Summary, args.Origin, args.Destination, func(id string) string {
			return storage.TaskFilePath(s.storageDir, id)
		})
		if err != nil {
			return toolError(err), nil, nil
		}

		if _, err := storage.WriteTaskFile(s.storageDir, task.ID, args.Summary, args.Origin, args.Destination, args.Content); err != nil {
			return toolError(err), nil, nil
		}

		return toolText(fmt.Sprintf("Created task **%s**: %s\n\nOrigin: %s → Destination: %s", task.ID, args.Summary, args.Origin, args.Destination)), nil, nil
	})

	type listTasksArgs struct {
		Origin      string `json:"origin,omitempty" jsonschema:"Filter by origin project"`
		Destination string `json:"destination,omitempty" jsonschema:"Filter by destination project"`
		Status      string `json:"status,omitempty" jsonschema:"Filter by status (open, in_progress, or closed)"`
	}
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "list_tasks",
		Description: "List tasks, optionally filtered by origin, destination, or status.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args listTasksArgs) (*mcp.CallToolResult, any, error) {
		if args.Status != "" && args.Status != "open" && args.Status != "in_progress" && args.Status != "closed" {
			return toolError(fmt.Errorf("status must be \"open\", \"in_progress\", or \"closed\"")), nil, nil
		}
		tasks, err := db.ListTasks(s.db, db.TaskFilter{
			Origin:      args.Origin,
			Destination: args.Destination,
			Status:      args.Status,
		})
		if err != nil {
			return toolError(err), nil, nil
		}
		if len(tasks) == 0 {
			return toolText("No tasks found."), nil, nil
		}
		text := ""
		for _, t := range tasks {
			text += fmt.Sprintf("- **%s** [%s]: %s (%s → %s)\n", t.ID, t.Status, t.Summary, t.Origin, t.Destination)
		}
		return toolText(text), nil, nil
	})

	type getTaskArgs struct {
		ID string `json:"id" jsonschema:"Task ID (e.g. ERIC-1)"`
	}
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "get_task",
		Description: "Get full details and content for a task by its ID.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args getTaskArgs) (*mcp.CallToolResult, any, error) {
		task, err := db.GetTask(s.db, args.ID)
		if err != nil {
			return toolError(err), nil, nil
		}
		if task == nil {
			return toolError(fmt.Errorf("task %s not found", args.ID)), nil, nil
		}

		content, err := storage.ReadTaskFile(task.FilePath)
		if err != nil {
			return toolError(fmt.Errorf("task %s: %w", task.ID, err)), nil, nil
		}
		return toolText(content), nil, nil
	})

	type updateTaskArgs struct {
		ID      string `json:"id" jsonschema:"Task ID (e.g. ERIC-1)"`
		Content string `json:"content" jsonschema:"Update content to append in Markdown"`
	}
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "update_task",
		Description: "Append an update to an existing task.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args updateTaskArgs) (*mcp.CallToolResult, any, error) {
		if strings.TrimSpace(args.Content) == "" {
			return toolError(fmt.Errorf("update content is required")), nil, nil
		}
		task, err := db.GetTask(s.db, args.ID)
		if err != nil {
			return toolError(err), nil, nil
		}
		if task == nil {
			return toolError(fmt.Errorf("task %s not found", args.ID)), nil, nil
		}

		if err := storage.AppendToTaskFile(task.FilePath, args.Content); err != nil {
			return toolError(fmt.Errorf("task %s: %w", task.ID, err)), nil, nil
		}
		return toolText(fmt.Sprintf("Updated task %s.", args.ID)), nil, nil
	})

	type searchTasksArgs struct {
		Query string `json:"query" jsonschema:"Text to search for in task content (case-insensitive)"`
	}
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "search_tasks",
		Description: "Search task content for a query string. Returns tasks whose Markdown content contains the query (case-insensitive).",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args searchTasksArgs) (*mcp.CallToolResult, any, error) {
		if strings.TrimSpace(args.Query) == "" {
			return toolError(fmt.Errorf("search query is required")), nil, nil
		}
		tasks, err := db.ListTasks(s.db, db.TaskFilter{})
		if err != nil {
			return toolError(err), nil, nil
		}
		var matches []db.Task
		for _, t := range tasks {
			if t.FilePath == "" {
				continue
			}
			found, err := storage.FileContains(t.FilePath, args.Query)
			if err != nil {
				continue
			}
			if found {
				matches = append(matches, t)
			}
		}
		if len(matches) == 0 {
			return toolText(fmt.Sprintf("No tasks found matching %q.", args.Query)), nil, nil
		}
		text := fmt.Sprintf("Tasks matching %q:\n", args.Query)
		for _, t := range matches {
			text += fmt.Sprintf("- **%s** [%s]: %s (%s → %s)\n", t.ID, t.Status, t.Summary, t.Origin, t.Destination)
		}
		return toolText(text), nil, nil
	})

	type closeTaskArgs struct {
		ID string `json:"id" jsonschema:"Task ID (e.g. ERIC-1)"`
	}
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "close_task",
		Description: "Close a task, marking it as done.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args closeTaskArgs) (*mcp.CallToolResult, any, error) {
		task, err := db.GetTask(s.db, args.ID)
		if err != nil {
			return toolError(err), nil, nil
		}
		if task == nil {
			return toolError(fmt.Errorf("task %s not found", args.ID)), nil, nil
		}
		if task.Status == "closed" {
			return toolError(fmt.Errorf("task %s is already closed", task.ID)), nil, nil
		}
		if err := db.UpdateTaskStatus(s.db, args.ID, "closed"); err != nil {
			return toolError(err), nil, nil
		}
		return toolText(fmt.Sprintf("Task %s closed.", args.ID)), nil, nil
	})

	type startTaskArgs struct {
		ID string `json:"id" jsonschema:"Task ID (e.g. ERIC-1)"`
	}
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "start_task",
		Description: "Mark a task as in-progress. Multiple tasks can be in-progress simultaneously.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args startTaskArgs) (*mcp.CallToolResult, any, error) {
		if err := db.StartTask(s.db, args.ID); err != nil {
			return toolError(err), nil, nil
		}
		return toolText(fmt.Sprintf("Task %s is now in progress.", args.ID)), nil, nil
	})
}

func toolText(text string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: text}},
	}
}

func toolError(err error) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Error: %v", err)}},
		IsError: true,
	}
}

