package mssql

import (
	"fmt"
	"os/exec"
)

func RestoreDatabase(container, user, password, backupFile, newDBName, mdfName, ldfName string) (string, error) {
	if password == "" {
		return "", fmt.Errorf("SA password required. Use the -p command to set the pw")
	}

	if backupFile == "" {
		return "", fmt.Errorf("backupFilename must be provoded. Use the -f command to set the filename")
	}

	if len(backupFile) < 5 || backupFile[len(backupFile)-4:] != ".bak" {
		return "", fmt.Errorf("filename must be of type .bak")
	}

	if mdfName == "" || ldfName == "" {
		return "", fmt.Errorf("mdf and ldf names must be provoded")
	}

	query := fmt.Sprintf(
		`
RESTORE DATABASE %s 
FROM DISK = N'/var/opt/mssql/backup/%s'
WITH 
MOVE '%s' TO '/var/opt/mssql/data/%s.mdf',
MOVE '%s' TO '/var/opt/mssql/data/%s_log.ldf'
`,
		newDBName, backupFile,
		mdfName, newDBName,
		ldfName, newDBName,
	)

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
		return "", fmt.Errorf("failed to restore database: %w\nOutput: %s", err, string(output))
	}

	return query, nil
}
