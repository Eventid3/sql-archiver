# SQL Archiver

WORK IN PROCESS

SQL Archiver is a simpler interactive TUI tool, that allows the user to backup and restore MSSQL databases from Docker MSSQL containers.

## Installation

SQL Archiver is best installed from via the go cli install command:

```go
go install "github.com/Eventid3/sql-archiver"
```

Alternatively, it can be build from the source code:

```go
git clone https://github.com/Eventid3/sql-archiver
cd sql-archiver
go build -o sql-archiver main.go
```

## Setup

It is possible to setup a config file, in order to avoid logging in to the Docker container every time the tool is used. This is only recommended on dev environments, as the password will be stored in clear text.

The config file must be located at `$HOME/.config/sql-archiver/` and must be named `config.yaml`, `config.json` or `config.toml`. Then in the config file, supply your container name, username and password. As an example, using .yaml:

```yaml
container: mssql-container
user: sa
password: SomeStrongP@ssw0rd
```

## Usage

Simply run the tool, enter the container, user and password, and choose either the Backup or Restore action. At the time of writing, backup files must be places in `/var/opt/mssql/backup/` inside the Docker container, in order for the tool to gain access to the .bak files.

The two flows will guide the user through the process.
