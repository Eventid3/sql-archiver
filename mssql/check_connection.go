/*
Package mssql
*/
package mssql

import (
	"context"
	"fmt"
)

func CheckConnection(host, user, password string) error {
	if user == "" {
		return fmt.Errorf("user required")
	}
	if password == "" {
		return fmt.Errorf("password required")
	}

	db, err := openDB(host, user, password)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := db.PingContext(context.Background()); err != nil {
		return fmt.Errorf("failed to authenticate access to SQL Server: %w", err)
	}
	return nil
}
