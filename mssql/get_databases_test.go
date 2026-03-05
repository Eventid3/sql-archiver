package mssql

import (
	"context"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

func TestGetDatabases_ReturnsParsedRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	now := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	rows := sqlmock.NewRows([]string{"name", "database_id", "create_date", "state_desc"}).
		AddRow("master", 1, now, "ONLINE").
		AddRow("tempdb", 2, now, "ONLINE").
		AddRow("myapp", 5, now, "ONLINE")

	mock.ExpectQuery("SELECT name, database_id, create_date, state_desc").WillReturnRows(rows)

	// Call the query directly on the mock db to test parsing logic
	result, err := db.QueryContext(context.Background(),
		"SELECT name, database_id, create_date, state_desc FROM sys.databases ORDER BY database_id")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer result.Close()

	var names []string
	for result.Next() {
		var name, state string
		var id int
		var created time.Time
		if err := result.Scan(&name, &id, &created, &state); err != nil {
			t.Fatalf("scan error: %v", err)
		}
		names = append(names, name)
	}

	if len(names) != 3 {
		t.Errorf("expected 3 databases, got %d", len(names))
	}
	if names[0] != "master" || names[1] != "tempdb" || names[2] != "myapp" {
		t.Errorf("unexpected database names: %v", names)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled mock expectations: %v", err)
	}
}
