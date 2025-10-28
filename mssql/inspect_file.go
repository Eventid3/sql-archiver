package mssql

import (
	"fmt"
	"os/exec"
	"strings"
)

func InspectBackupFile(container, user, password, file string) (BackupEntry, error) {
	if password == "" {
		return BackupEntry{}, fmt.Errorf("SA password required. Use the -p command to set the pw")
	}

	if file == "" {
		return BackupEntry{}, fmt.Errorf("filename must be provoded. Use the -f command to set the filename")
	}

	if len(file) < 5 || file[len(file)-4:] != ".bak" {
		return BackupEntry{}, fmt.Errorf("filename must be of type .bak")
	}

	query := fmt.Sprintf("RESTORE FILELISTONLY FROM DISK = N'/var/opt/mssql/backup/%s'", file)

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
		return BackupEntry{}, fmt.Errorf("error listing entries from .bak file: %w\nOutput: %s", err, output)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")

	columns := make(map[int]string)

	logicalNameCol := 0
	typeCol := 2
	sizeCol := 4
	backupSizeCol := 12

	result := BackupEntry{}

	for i, line := range lines {
		if i == 0 {
			fields := strings.Fields(line)
			for j, field := range fields {
				columns[j] = field
			}
			continue
		}

		if i == 1 {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 13 {
			continue
		}
		switch fields[typeCol] {
		case "D":
			result.MdfFile = MdfEntry{
				Name:       fields[logicalNameCol],
				Size:       fields[sizeCol],
				BackupSize: fields[backupSizeCol],
			}
		case "L":
			result.LdfFile = LdfEntry{
				Name: fields[logicalNameCol],
				Size: fields[sizeCol],
			}
		}
	}

	return result, nil
}
