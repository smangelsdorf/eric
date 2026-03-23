# Eric — Development Guide

## Build & Test

```bash
go build ./...          # Build all packages
go test ./...           # Run all tests
go test -race ./...     # Run tests with race detector
go vet ./...            # Static analysis
```

## Go Conventions

- Format code with `gofmt` (editor should handle this, but `gofmt -w .` if needed)
- Use `go vet ./...` before committing to catch common issues
- Errors are values — return and handle them, don't panic
- Keep packages small and focused; avoid package `utils`
- Test files live alongside their code: `foo.go` / `foo_test.go`
- Use table-driven tests where appropriate
- Prefer the standard library over external dependencies

## Project Structure

- `go.mod` — module: `github.com/smangelsdorf/eric`
- Task content stored as Markdown flat files on disk
- SQLite used for the task/project index

## Dependencies

- **MCP**: [modelcontextprotocol/go-sdk](https://github.com/modelcontextprotocol/go-sdk) — official Go SDK for MCP
- **SQLite**: [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite) — pure Go, no CGO required
