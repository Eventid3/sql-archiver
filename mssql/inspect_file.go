package mssql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/Eventid3/sql-archiver/domain"
)

func InspectBackupFile(host, user, password, file string) (domain.BackupEntry, error) {
	if password == "" {
		return domain.BackupEntry{}, fmt.Errorf("SA password required. Use the -p command to set the pw")
	}
	if file == "" {
		return domain.BackupEntry{}, fmt.Errorf("filename must be provided. Use the -f command to set the filename")
	}
	if len(file) < 5 || file[len(file)-4:] != ".bak" {
		return domain.BackupEntry{}, fmt.Errorf("filename must be of type .bak")
	}

	db, err := openDB(host, user, password)
	if err != nil {
		return domain.BackupEntry{}, err
	}
	defer db.Close()

	return inspectBackupFile(db, file)
}

// inspectBackupFile is the testable core: it accepts any *sql.DB so tests can inject a mock.
func inspectBackupFile(db *sql.DB, file string) (domain.BackupEntry, error) {
	query := fmt.Sprintf(
		"RESTORE FILELISTONLY FROM DISK = N'/var/opt/mssql/backup/%s'",
		strings.ReplaceAll(file, "'", "''"),
	)

	rows, err := db.QueryContext(context.Background(), query)
	if err != nil {
		return domain.BackupEntry{}, fmt.Errorf("error reading entries from .bak file: %w", err)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return domain.BackupEntry{}, fmt.Errorf("error reading column names: %w", err)
	}

	result := domain.BackupEntry{}

	for rows.Next() {
		// Scan all columns into interface{} to be flexible across SQL Server versions.
		values := make([]interface{}, len(cols))
		valuePtrs := make([]interface{}, len(cols))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return domain.BackupEntry{}, fmt.Errorf("error scanning file list row: %w", err)
		}

		colIndex := make(map[string]int, len(cols))
		for i, c := range cols {
			colIndex[strings.ToLower(c)] = i
		}

		logicalName := stringVal(values, colIndex, "logicalname")
		fileType := stringVal(values, colIndex, "type")
		size := int64Val(values, colIndex, "size")
		backupSize := int64Val(values, colIndex, "backupsizeinbytes")

		switch fileType {
		case "D":
			result.MdfFile = domain.MdfEntry{
				Name:       logicalName,
				Size:       fmt.Sprintf("%d", size),
				BackupSize: fmt.Sprintf("%d", backupSize),
			}
		case "L":
			result.LdfFile = domain.LdfEntry{
				Name: logicalName,
				Size: fmt.Sprintf("%d", size),
			}
		}
	}
	if err := rows.Err(); err != nil {
		return domain.BackupEntry{}, fmt.Errorf("error iterating file list: %w", err)
	}

	return result, nil
}

func stringVal(values []interface{}, colIndex map[string]int, col string) string {
	i, ok := colIndex[col]
	if !ok || i >= len(values) || values[i] == nil {
		return ""
	}
	switch v := values[i].(type) {
	case string:
		return strings.TrimSpace(v)
	case []byte:
		return strings.TrimSpace(string(v))
	default:
		return fmt.Sprintf("%v", v)
	}
}

func int64Val(values []interface{}, colIndex map[string]int, col string) int64 {
	i, ok := colIndex[col]
	if !ok || i >= len(values) || values[i] == nil {
		return 0
	}
	switch v := values[i].(type) {
	case int64:
		return v
	case int32:
		return int64(v)
	case int:
		return int64(v)
	case float64:
		return int64(v)
	default:
		return 0
	}
}
