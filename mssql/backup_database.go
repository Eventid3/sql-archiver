package mssql

import (
	"fmt"
	"os/exec"
)

func BackupDatabase(container, user, password, database, filename string) error {
	if password == "" {
		return fmt.Errorf("SA password required. Use the -p command to set the pw")
	}

	if database == "" {
		return fmt.Errorf("database name must be provided. Use the -d flag to set the database name")
	}

	if filename == "" {
		return fmt.Errorf("filename must be provided. Use the -f command to set the filename")
	}

	if len(filename) < 5 || filename[len(filename)-4:] != ".bak" {
		return fmt.Errorf("filename must end with .bak")
	}

	// t := time.Now()
	// timeString := t.Format("2006-01-02T15-04-05")

	// if timestamp {
	// 	filename = filename[:len(filename)-4] + "_" + timeString + ".bak"
	// }

	query := fmt.Sprintf("BACKUP DATABASE [%s] TO DISK = N'/var/opt/mssql/backup/%s'", database, filename)

	dockerCmd := exec.Command(
		"docker", "exec", "-i", container,
		"/opt/mssql-tools18/bin/sqlcmd",
		"-S", "localhost",
		"-U", "sa",
		"-P", password,
		"-Q", query,
		"-C",
	)

	output, err := dockerCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to backup database: %w\nOutput: %s", err, string(output))
	}

	return nil
}
