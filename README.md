# eric

An MCP (Model Context Protocol) server for cross-project task management. Eric acts as a private, local ticket system — tasks are stored as Markdown files in a local repository and indexed with SQLite for fast querying.

## Why

When working across multiple projects, context about tasks, decisions, and follow-ups gets lost between conversations. Eric provides a shared task store that any MCP-capable tool can read from and write to, keeping task information accessible regardless of which project you're currently in.

## Key ideas

- **MCP server** — exposes task operations as MCP tools, so AI assistants can create, query, and update tasks naturally
- **Markdown storage** — task content is stored as date-stamped Markdown files on disk, human-readable and easy to inspect
- **SQLite index** — tracks task metadata (origin, destination, status, file path) for fast search and filtering
- **Local and private** — everything stays on your machine, no external services required

## Data model

A **task** represents a piece of work passed from one project to another:

- **ID** — stable identifier for referencing tasks across sessions
- **Summary** — brief one-line description of the task
- **Origin project** — the project that created the task
- **Destination project** — the project the task is intended for
- **Status** — lifecycle state of the task
- **Content** — date-stamped Markdown file on disk, referenced by file path in the index

A **project** is a named entity registered with Eric. Projects are not tied to working directories (because of Git worktrees, among other reasons) — they must be explicitly registered before tasks can reference them.

## Setup

Build the binary, then register it as an MCP server:

```bash
go build -o eric .

claude mcp add --transport stdio --scope user \
  --allow-tools "mcp__eric__list_*,mcp__eric__get_task,mcp__eric__search_tasks" \
  eric -- /path/to/eric
```

This pre-approves the read-only tools (`list_projects`, `list_tasks`, `get_task`, `search_tasks`) so they won't prompt for confirmation each time. Mutating tools (`create_task`, `update_task`, `close_task`, `register_project`) will still require approval.

Data is stored in `~/.eric/` (SQLite database and Markdown task files).

## TUI

A terminal interface for browsing and managing tasks directly:

```bash
go build -o eric-tui ./cmd/eric-tui
./eric-tui
```

The TUI shows all tasks in a navigable table (open tasks first), with keyboard shortcuts to view task details, close, and reopen tasks. It automatically refreshes when tasks are created or modified by the MCP server in another session.
