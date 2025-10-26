package mssql

import (
	"fmt"
	"os/exec"
	"strings"
)

func GetDatabases(container, user, password string) ([]DBItem, error) {
	if user == "" {
		return nil, fmt.Errorf("user required")
	}
	if password == "" {
		return nil, fmt.Errorf("password required")
	}

	query := "SELECT name, database_id, create_date, state_desc FROM sys.databases ORDER BY database_id"

	dockerCmd := exec.Command(
		"docker", "exec", "-i", container,
		"/opt/mssql-tools18/bin/sqlcmd",
		"-S", "localhost",
		"-U", user,
		"-P", password,
		"-Q", query,
		"-C",
		"-h", "-1", // remove headers
		"-s", "|",
		"-W", // remove trailing whitespace
	)

	output, err := dockerCmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to list databases: %w\nOutput: %s", err, string(output))
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")

	if len(lines) == 0 {
		return nil, fmt.Errorf("no databases found")
	}

	var result []DBItem

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		fields := strings.Split(line, "|")
		if len(fields) >= 4 {
			name := fields[0]
			id := fields[1]
			created := fields[2]
			state := fields[3]

			// only get date part
			if len(created) > 10 {
				created = created[:10]
			}

			result = append(result, DBItem{name, id, created, state})
		}
	}
	return result, nil
}
