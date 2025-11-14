/*
Package mssql
*/
package mssql

import (
	"fmt"
	"os/exec"
)

func CheckConnection(container, user, password string) error {
	if user == "" {
		return fmt.Errorf("user required")
	}
	if password == "" {
		return fmt.Errorf("password required")
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
	)

	output, err := dockerCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to authenticate access to the docker container: %w\nOutput: %s", err, string(output))
	}
	return nil
}
