package mssql

import (
	"context"
	"fmt"
	"strings"
)

func BackupDatabase(host, user, password, database, filename string) error {
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

	db, err := openDB(host, user, password)
	if err != nil {
		return err
	}
	defer db.Close()

	query := fmt.Sprintf(
		"BACKUP DATABASE [%s] TO DISK = N'/var/opt/mssql/backup/%s'",
		strings.ReplaceAll(database, "]", "]]"),
		strings.ReplaceAll(filename, "'", "''"),
	)

	_, err = db.ExecContext(context.Background(), query)
	if err != nil {
		return fmt.Errorf("failed to backup database: %w", err)
	}
	return nil
}
