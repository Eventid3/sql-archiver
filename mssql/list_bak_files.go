package mssql

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Eventid3/sql-archiver/domain"
)

func ListBackupFilesInContainer(host, user, password string) ([]domain.BakFile, error) {
	if password == "" {
		return nil, fmt.Errorf("SA password required. Use the -p command to set the pw")
	}

	db, err := openDB(host, user, password)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// sys.dm_os_enumerate_filesystem is available on SQL Server 2019+.
	// It lists files with size and timestamp without requiring Docker exec.
	rows, err := db.QueryContext(context.Background(), `
		SELECT file_or_directory_name, size_in_bytes, last_write_time
		FROM sys.dm_os_enumerate_filesystem('/var/opt/mssql/backup', '*')
		WHERE is_directory = 0
		  AND file_or_directory_name LIKE '%.bak'
		ORDER BY last_write_time DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to list backup files: %w", err)
	}
	defer rows.Close()

	var result []domain.BakFile
	for rows.Next() {
		var name string
		var sizeBytes int64
		var modTime time.Time

		if err := rows.Scan(&name, &sizeBytes, &modTime); err != nil {
			return nil, fmt.Errorf("failed to scan backup file row: %w", err)
		}

		// Only include .bak files (filter is already in SQL but double-check)
		if !strings.HasSuffix(strings.ToLower(name), ".bak") {
			continue
		}

		result = append(result, domain.BakFile{
			Name: name,
			Size: formatBytes(sizeBytes),
			Date: modTime.Format("2006-01-02 15:04"),
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating backup file rows: %w", err)
	}

	return result, nil
}

func formatBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}
