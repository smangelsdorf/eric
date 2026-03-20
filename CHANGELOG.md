# Changelog

## 2026-03-20

- Add in-progress task status with `start_task` MCP tool, TUI shortcut (`s`), and sort priority above open tasks
- Add terminal UI (`eric-tui`) for browsing and managing tasks
  - Table view with sorting (open first, then by date)
  - Task detail modal with scrollable markdown content
  - Close/reopen tasks directly from the TUI
  - Live refresh via SQLite data_version polling

## 2026-03-16

- Initial release
- MCP server with stdio transport
- Project registration and listing
- Task creation, listing, viewing, updating, closing
- Cross-project task routing (origin/destination)
- Full-text search across task content
- SQLite index with markdown flat-file storage
- WAL journal mode and busy_timeout for concurrent access
