# AGENTS.md — sql-archiver

A Go CLI tool for backing up and restoring MSSQL databases via `.bak` files, with a
Bubble Tea TUI. Module path: `github.com/Eventid3/sql-archiver`.

---

## Repository Layout

```
cmd/            CLI entry point (cobra)
config/         Config file I/O — LoadConnections, SaveConnection, DeleteConnection
domain/         Shared data models (no external dependencies)
mssql/          MSSQL driver logic; testable via unexported helpers + sqlmock
interactive/    Bubble Tea TUI models, styles, messages
main.go         Delegates to cmd.Execute()
Makefile        Primary task runner
```

---

## Build / Lint / Test Commands

### Build

```sh
make build          # compiles to build/sql-archiver
go build .          # used in CI (go.yml)
```

### Run all tests

```sh
make test           # go test ./...
make test-verbose   # go test -v ./...
make test-cover     # go test -coverprofile=coverage.out ./... + HTML report
```

### Run a single test

```sh
# One test function by name (regex):
go test ./mssql/ -run TestInspectBackupFile_WindowsPathWithSpaces

# One test function in interactive package:
go test ./interactive/ -run TestLoginModel_Update_TabAdvancesFocus

# Verbose single test:
go test -v ./mssql/ -run TestFormatBytes

# All tests in one package:
go test ./mssql/
go test ./interactive/
```

### Install binary

```sh
make install        # copies build/sql-archiver to ~/go/bin
```

### Lint / Format

No linter is configured in the Makefile or CI. Use standard Go tooling:

```sh
go vet ./...
gofmt -l .          # list files needing formatting
gofmt -w .          # reformat in place
```

---

## Code Style Guidelines

### General

- Follow standard Go idioms and `gofmt` formatting. No project-specific formatter
  configuration exists; default `gofmt` output is canonical.
- Keep functions small and focused. Prefer explicit error handling over panics.
- No third-party linter config (no `.golangci.yml`); `go vet` is the baseline.

### Imports

- Group imports in the standard Go order (enforced by `gofmt`/`goimports`):
  1. Standard library
  2. Third-party packages
  3. Internal packages (`github.com/Eventid3/sql-archiver/...`)
- Use blank imports (`_ "..."`) only for database driver side-effects (see `mssql/db.go`).
- Alias third-party packages only when the default name is ambiguous, e.g.:
  ```go
  sqlmock "github.com/DATA-DOG/go-sqlmock"
  tea     "github.com/charmbracelet/bubbletea"
  ```

### Naming Conventions

- **Packages**: short, lowercase, single word (`mssql`, `interactive`, `domain`, `cmd`).
- **Exported functions/types**: `PascalCase` (`GetDatabases`, `BackupDatabase`,
  `NewLoginModel`, `InitialModel`).
- **Unexported helpers**: `camelCase`; expose a public wrapper that opens the real DB and
  an unexported core function that accepts `*sql.DB` for testability (see
  `inspectBackupFile` / `InspectBackupFile` pattern in `mssql/inspect_file.go`).
- **Types**: `PascalCase` structs (`DBItem`, `BakFile`, `BackupEntry`, `MdfEntry`).
- **Message types** (Bubble Tea): unexported, `camelCase` with `Msg` suffix
  (`loginDoneMsg`, `actionSelectedMsg`, `restoreExecMsg`).
- **Model constructors**: `New<Name>` (`NewLoginModel`, `NewActionModel`) or
  `Initial<Name>` for top-level entry points (`InitialModel`, `InitialModelWithConfig`).
- **Style variables**: exported, `PascalCase` with `Style` suffix (`DocStyle`,
  `HeadingStyle`, `ErrorTextStyle`).

### Types and Structs

- Domain models live exclusively in `domain/`; no business logic in that package.
- Prefer `string` fields in domain structs for display data that will be rendered in
  the TUI (avoids repetitive `fmt.Sprintf` at render time).
- Use `[]interface{}` with a column-index map when scanning dynamic SQL result sets
  to stay compatible across SQL Server versions (see `inspect_file.go`).

### Error Handling

- Always return errors; never silently discard them.
- Wrap errors with context using `fmt.Errorf("description: %w", err)`.
- Validate all required inputs at the top of public functions before touching the DB:
  ```go
  if password == "" {
      return nil, fmt.Errorf("password required")
  }
  ```
- Use `defer db.Close()` and `defer rows.Close()` immediately after successful open/query.
- Check `rows.Err()` after iterating result sets.
- In the TUI, surface errors by passing them into model constructors (e.g.
  `NewLoginModel(err)`); never `log.Fatal` from a model.

### Bubble Tea (TUI) Patterns

- Each screen is a separate model type implementing `tea.Model` (`Init`, `Update`, `View`).
- `parentModel` (`interactive/parent_model.go`) owns navigation: it intercepts
  message types and swaps `activeModel`; child models emit typed messages instead of
  navigating directly.
- Message types are unexported structs defined in `interactive/messages.go`.
- Styles are centralized in `interactive/style.go` as package-level `var` declarations
  using `lipgloss`.
- Return `tea.Quit` from `Update` on `"ctrl+c"` or `"esc"` where appropriate.
- Use `tea.Batch` to combine multiple `tea.Cmd` values.
- SQL calls are always made inside a `tea.Cmd` closure (never inside `Update` directly)
  so the TUI remains responsive.

### SQL / Database Patterns

- All DB access lives in the `mssql` package. No SQL in `interactive/` or `domain/`.
- Use `context.Background()` for all query/exec calls.
- SQL identifiers (database names) are escaped with `[` `]` brackets and `]]` doubling;
  string literals use `''` doubling — never use `fmt.Sprintf` with unescaped user input
  for identifiers.
- Default host is `localhost:1433`; always provide a fallback in `openDB`.

### Testing

- Tests use Go's standard `testing` package — no external test framework.
- SQL logic is tested with `github.com/DATA-DOG/go-sqlmock` by injecting a `*sql.DB`
  into unexported helper functions.
- Test files are in the same package as the code under test (white-box testing).
- Table-driven tests use anonymous structs with `input`/`want` fields.
- Helper functions use `t.Helper()` and accept `(t *testing.T, ...)`.
- Test function names follow `Test<Unit>_<Scenario>` (e.g.
  `TestInspectBackupFile_EmptyResult`, `TestLoginModel_Update_EscQuitsApp`).
- Always verify `mock.ExpectationsWereMet()` at the end of sqlmock tests.
- TUI model tests call `View()` and assert on rendered string content; they send
  `tea.KeyMsg` values via `Update()` to exercise keyboard interactions.

### File / Package Conventions

- One package per directory; package name matches directory name.
- Package-level doc comments go on the `package` declaration in one file per package
  (e.g. `parent_model.go` for `interactive`, `check_connection.go` for `mssql`).
- No generated code or build tags are currently used.

---

## CI

GitHub Actions (`.github/workflows/go.yml`) runs on push/PR to `master`:
- Go version: `1.25`
- Steps: `go build .` only (no test step in CI currently).
