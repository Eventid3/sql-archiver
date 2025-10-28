package mssql

import (
	"fmt"
	"os/exec"
	"strings"
)

func InspectBackupFile(container, user, password, file string) (string, error) {
	if password == "" {
		return "", fmt.Errorf("SA password required. Use the -p command to set the pw")
	}

	if file == "" {
		return "", fmt.Errorf("filename must be provoded. Use the -f command to set the filename")
	}

	if len(file) < 5 || file[len(file)-4:] != ".bak" {
		return "", fmt.Errorf("filename must be of type .bak")
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
		return "", fmt.Errorf("error listing entries from .bak file: %w\nOutput: %s", err, output)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")

	fmt.Printf("\n%-30s %-12s %-12s\n", "DATABASE NAME", "SIZE", "BACKUPSIZE")
	fmt.Println(strings.Repeat("-", 80))

	// columns := make(map[int]string)

	// logicalNameCol := 0
	// typeCol := 2
	// sizeCol := 4
	// backupSizeCol := 12

	result := ""

	for i, line := range lines {
		result += fmt.Sprintf("Line %d:\n%s\n", i, line)
		// if i == 0 {
		// 	fields := strings.Fields(line)
		// 	for j, field := range fields {
		// 		columns[j] = field
		// 	}
		// 	continue
		// }
		//
		// if i == 1 {
		// 	continue
		// }
		//
		// fields := strings.Fields(line)
	}

	return "", nil
}
