package mssql

import (
	"context"
	"fmt"
	"time"

	"github.com/Eventid3/sql-archiver/domain"
)

func GetDatabases(host, user, password string) ([]domain.DBItem, error) {
	if user == "" {
		return nil, fmt.Errorf("user required")
	}
	if password == "" {
		return nil, fmt.Errorf("password required")
	}

	db, err := openDB(host, user, password)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.QueryContext(
		context.Background(),
		"SELECT name, database_id, create_date, state_desc FROM sys.databases ORDER BY database_id",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list databases: %w", err)
	}
	defer rows.Close()

	var result []domain.DBItem
	for rows.Next() {
		var name, state string
		var id int
		var created time.Time
		if err := rows.Scan(&name, &id, &created, &state); err != nil {
			return nil, fmt.Errorf("failed to scan database row: %w", err)
		}
		result = append(result, domain.DBItem{
			Name:    name,
			ID:      fmt.Sprintf("%d", id),
			Created: created.Format("2006-01-02"),
			State:   state,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating database rows: %w", err)
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("no databases found")
	}
	return result, nil
}
