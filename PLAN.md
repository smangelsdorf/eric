# Project Plan

## Phase 1: Foundation

- [ ] SQLite database setup and connection management
- [ ] Projects table and registration (name, description)
- [ ] Tasks table (id, summary, origin project, destination project, status, file path)
- [ ] Markdown file storage (directory structure, date-stamped file format)
- [ ] Task creation: write Markdown file + insert index row

## Phase 2: MCP Server

- [ ] MCP server scaffolding (stdio transport)
- [ ] `register_project` tool — teach Eric about a new project
- [ ] `list_projects` tool
- [ ] `create_task` tool — creates Markdown file and index entry
- [ ] `list_tasks` tool — filter by project, status, origin/destination
- [ ] `get_task` tool — read task metadata + Markdown content
- [ ] `update_task` tool — append to Markdown file, update index
- [ ] `close_task` tool

## Phase 3: Cross-Project Workflow

- [ ] Query tasks by destination (what's incoming for this project?)
- [ ] Query tasks by origin (what did this project send out?)
- [ ] Search task content

## Phase 4: Polish

- [ ] Error handling and validation
- [ ] Documentation for MCP client configuration
- [ ] Testing
