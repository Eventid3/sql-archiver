package mssql

import (
	"context"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

func TestFormatBytes(t *testing.T) {
	cases := []struct {
		input int64
		want  string
	}{
		{0, "0 B"},
		{512, "512 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{66060288, "63.0 MB"},
	}
	for _, c := range cases {
		got := formatBytes(c.input)
		if got != c.want {
			t.Errorf("formatBytes(%d) = %q, want %q", c.input, got, c.want)
		}
	}
}

func TestListBackupFiles_ParsesRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	modTime := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

	rows := sqlmock.NewRows([]string{"file_or_directory_name", "size_in_bytes", "last_write_time"}).
		AddRow("vej-eu.bak", int64(66060288), modTime).
		AddRow("mydb.bak", int64(8388608), modTime)

	mock.ExpectQuery("sys.dm_os_enumerate_filesystem").WillReturnRows(rows)

	result, err := db.QueryContext(context.Background(), "SELECT * FROM sys.dm_os_enumerate_filesystem('/var/opt/mssql/backup', '*') WHERE is_directory = 0 AND file_or_directory_name LIKE '%.bak'")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer result.Close()

	var files []string
	for result.Next() {
		var name string
		var size int64
		var ts time.Time
		if err := result.Scan(&name, &size, &ts); err != nil {
			t.Fatalf("scan error: %v", err)
		}
		files = append(files, name)
	}

	if len(files) != 2 {
		t.Errorf("expected 2 files, got %d", len(files))
	}
	if files[0] != "vej-eu.bak" {
		t.Errorf("expected first file to be 'vej-eu.bak', got %q", files[0])
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled mock expectations: %v", err)
	}
}
