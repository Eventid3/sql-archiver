package mssql

import (
	"fmt"
	"os/exec"
	"strings"
)

func ListBackupFilesInContainer(container, user, password string) ([]BakFile, error) {
	if password == "" {
		return nil, fmt.Errorf("SA password required. Use the -p command to set the pw")
	}

	dockerCmd := exec.Command(
		"docker", "exec", "-i", container,
		"ls", "-lhgG", "/var/opt/mssql/backup",
	)

	output, err := dockerCmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to list backup files: %w\nOutput: %s", err, string(output))
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")

	result := []BakFile{}

	if len(lines) == 0 {
		return result, nil
	}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "total") {
			continue
		}

		// Use Fields instead of Split to handle multiple spaces
		fields := strings.Fields(line)

		// ls -lhgG output format:
		// [0] permissions, [1] links, [2] size, [3] month, [4] day, [5] time/year, [6] filename
		if len(fields) >= 7 {
			size := fields[2]
			date := fmt.Sprintf("%s %s %s", fields[3], fields[4], fields[5])
			name := fields[6]

			// Handle filenames with spaces (join remaining fields)
			if len(fields) > 7 {
				name = strings.Join(fields[6:], " ")
			}
			result = append(result, BakFile{
				Size: size,
				Date: date,
				Name: name,
			})
		}
	}
	return result, nil
}
