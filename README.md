# SQL Archiver

SQL Archiver is an interactive TUI tool for backing up and restoring MSSQL databases from Docker MSSQL containers, using `.bak` files.

## Installation

Install via the Go CLI:

```sh
go install github.com/Eventid3/sql-archiver@latest
```

Or build from source:

```sh
git clone https://github.com/Eventid3/sql-archiver
cd sql-archiver
make build          # binary written to build/sql-archiver
make install        # copies binary to ~/go/bin
```

## Usage

Run the tool and the TUI will start:

```sh
sql-archiver
```

On first launch you will see a blank login form. Fill in the SQL Server host
(`localhost:1433` by default), username, and password, then press **Enter**.

After a successful login you will be asked to name and save the connection.
Saved connections appear in the picker at the top of the login screen on
subsequent launches, so you can connect with a single key-press.

### Keyboard shortcuts — login screen

| Key | Action |
|-----|--------|
| `↑` / `↓` | Navigate saved-connection list |
| `Enter` | Select a saved connection (pre-fills the form) |
| `n` | New connection (clears the form) |
| `d` | Delete the highlighted saved connection |
| `y` / any key | Confirm / cancel the delete prompt |
| `Tab` / `↓` | Advance focus to the next form field |
| `Shift+Tab` / `↑` | Move focus to the previous form field |
| `Esc` | Return from the form to the saved-connection picker |
| `Ctrl+C` | Quit |

### Backup files

Backup files must be placed in `/var/opt/mssql/backup/` inside the Docker
container for the tool to access them. The tool reads this directory via the
SQL Server `sys.dm_os_enumerate_filesystem` DMV (requires SQL Server 2019+).

## Saved connections

Connections are stored in `$HOME/.config/sql-archiver/config.yaml`:

```yaml
connections:
  - name: local
    host: localhost:1433
    user: sa
    password: SomeStrongP@ssw0rd
  - name: staging
    host: staging.internal:1433
    user: sa
    password: AnotherPassword
```

You can edit this file manually or let the TUI manage it. Passwords are stored
in plain text — only use saved connections in trusted, non-production environments.

## Development

```sh
make test           # run all tests
make test-verbose   # run all tests with verbose output
make test-cover     # generate coverage report (coverage.html)
go vet ./...        # static analysis
```

Run a single test:

```sh
go test ./interactive/ -run TestParentModel_FullBackupFlow
go test ./mssql/      -run TestInspectBackupFile_WindowsPathWithSpaces
```
