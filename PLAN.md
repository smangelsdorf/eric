# Project Plan

## Phase 1: Foundation

- [x] SQLite database setup and connection management
- [x] Projects table and registration (name, description)
- [x] Tasks table (id, summary, origin project, destination project, status, file path)
- [x] Markdown file storage (directory structure, date-stamped file format)
- [x] Task creation: write Markdown file + insert index row

## Phase 2: MCP Server

- [x] MCP server scaffolding (stdio transport)
- [x] `register_project` tool — teach Eric about a new project
- [x] `list_projects` tool
- [x] `create_task` tool — creates Markdown file and index entry
- [x] `list_tasks` tool — filter by project, status, origin/destination
- [x] `get_task` tool — read task metadata + Markdown content
- [x] `update_task` tool — append to Markdown file, update index
- [x] `close_task` tool

## Phase 3: Cross-Project Workflow

- [x] Query tasks by destination (what's incoming for this project?)
- [x] Query tasks by origin (what did this project send out?)
- [x] Search task content

## Phase 4: Polish

- [x] Error handling and validation
- [x] Documentation for MCP client configuration
- [x] Testing (db layer)

## Phase 5: Concurrency & Multi-Process Safety

- [x] Set `busy_timeout` pragma so concurrent writers wait instead of failing
- [ ] Add concurrency tests for simultaneous reads/writes

## Phase 6: TUI

- [x] Add Bubble Tea dependency
- [x] Single table view: task ID, summary, status, created date
- [x] Sort order: open tasks first, then descending by created date
- [x] Styling: cyan borders, white text (open), grey text (closed), green status for open, yellow highlight on dates older than 1 week
- [x] Hint bar at bottom: shortcut key in light green, rest of label in white
- [x] Keyboard actions: close/reopen (immediate, reversible)
- [x] Space to open task content in a centered modal popup (white border, white text, scrollable)
- [x] Live refresh via `PRAGMA data_version` polling
