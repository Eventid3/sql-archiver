package mssql

import (
	"context"
	"fmt"
	"strings"
)

func RestoreDatabase(host, user, password, backupFile, newDBName, mdfName, ldfName string) (string, error) {
	if password == "" {
		return "", fmt.Errorf("SA password required. Use the -p command to set the pw")
	}
	if backupFile == "" {
		return "", fmt.Errorf("backupFilename must be provided. Use the -f command to set the filename")
	}
	if len(backupFile) < 5 || backupFile[len(backupFile)-4:] != ".bak" {
		return "", fmt.Errorf("filename must be of type .bak")
	}
	if mdfName == "" || ldfName == "" {
		return "", fmt.Errorf("mdf and ldf names must be provided")
	}

	db, err := openDB(host, user, password)
	if err != nil {
		return "", err
	}
	defer db.Close()

	query := fmt.Sprintf(
		`RESTORE DATABASE [%s]
FROM DISK = N'/var/opt/mssql/backup/%s'
WITH
MOVE '%s' TO '/var/opt/mssql/data/%s.mdf',
MOVE '%s' TO '/var/opt/mssql/data/%s_log.ldf'`,
		strings.ReplaceAll(newDBName, "]", "]]"),
		strings.ReplaceAll(backupFile, "'", "''"),
		strings.ReplaceAll(mdfName, "'", "''"),
		strings.ReplaceAll(newDBName, "'", "''"),
		strings.ReplaceAll(ldfName, "'", "''"),
		strings.ReplaceAll(newDBName, "'", "''"),
	)

	_, err = db.ExecContext(context.Background(), query)
	if err != nil {
		return "", fmt.Errorf("failed to restore database: %w", err)
	}
	return query, nil
}
