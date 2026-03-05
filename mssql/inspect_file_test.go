package mssql

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/Eventid3/sql-archiver/domain"
)

func TestInspectBackupFile_WindowsPathWithSpaces(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	// Simulate RESTORE FILELISTONLY output for vej-eu.bak (Windows path with spaces in PhysicalName)
	cols := []string{
		"LogicalName", "PhysicalName", "Type", "FileGroupName",
		"Size", "MaxSize", "FileId", "CreateLSN", "DropLSN",
		"UniqueId", "ReadOnlyLSN", "ReadWriteLSN", "BackupSizeInBytes",
		"SourceBlockSize", "FileGroupId", "LogGroupGUID",
		"DifferentialBaseLSN", "DifferentialBaseGUID",
		"IsReadOnly", "IsPresent", "TDEThumbprint",
	}

	rows := sqlmock.NewRows(cols).
		AddRow(
			"VEJ-EU.Staging.Umbraco",
			`C:\Program Files\Microsoft SQL Server\MSSQL16.SQLEXPRESS\MSSQL\DATA\VEJ-EU.Production.Umbraco_newV2.mdf`,
			"D", "PRIMARY",
			int64(276824064), int64(35184372080640), int64(1),
			int64(0), int64(0), "A152F04E-E847-4E19-A76C-3CC0DD5A714C",
			int64(0), int64(0), int64(65404928),
			int64(4096), int64(1), nil,
			int64(1493000008938400138), "28EDAC35-8140-4DF9-A90F-2B5CDC73D198",
			false, true, nil,
		).
		AddRow(
			"VEJ-EU.Staging.Umbraco_log",
			`C:\Program Files\Microsoft SQL Server\MSSQL16.SQLEXPRESS\MSSQL\DATA\VEJ-EU.Production.Umbraco_newV2_log.ldf`,
			"L", "NULL",
			int64(1954545664), int64(2199023255552), int64(2),
			int64(0), int64(0), "56711F9E-3202-4635-A504-12F106A56B65",
			int64(0), int64(0), int64(0),
			int64(4096), int64(0), nil,
			int64(0), "00000000-0000-0000-0000-000000000000",
			false, true, nil,
		)

	mock.ExpectQuery("RESTORE FILELISTONLY").WillReturnRows(rows)

	entry, err := inspectBackupFile(db, "vej-eu.bak")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// This is the core regression test: Windows paths with spaces must not break parsing
	assertEntry(t, entry, domain.BackupEntry{
		MdfFile: domain.MdfEntry{
			Name:       "VEJ-EU.Staging.Umbraco",
			Size:       "276824064",
			BackupSize: "65404928",
		},
		LdfFile: domain.LdfEntry{
			Name: "VEJ-EU.Staging.Umbraco_log",
			Size: "1954545664",
		},
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled mock expectations: %v", err)
	}
}

func TestInspectBackupFile_LinuxPath(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	cols := []string{
		"LogicalName", "PhysicalName", "Type", "FileGroupName",
		"Size", "MaxSize", "FileId", "CreateLSN", "DropLSN",
		"UniqueId", "ReadOnlyLSN", "ReadWriteLSN", "BackupSizeInBytes",
		"SourceBlockSize", "FileGroupId", "LogGroupGUID",
		"DifferentialBaseLSN", "DifferentialBaseGUID",
		"IsReadOnly", "IsPresent", "TDEThumbprint",
	}

	rows := sqlmock.NewRows(cols).
		AddRow(
			"mydb",
			"/var/opt/mssql/data/mydb.mdf",
			"D", "PRIMARY",
			int64(8388608), int64(35184372080640), int64(1),
			int64(0), int64(0), "AAAAAAAA-BBBB-CCCC-DDDD-EEEEEEEEEEEE",
			int64(0), int64(0), int64(1048576),
			int64(512), int64(1), nil,
			int64(0), "00000000-0000-0000-0000-000000000000",
			false, true, nil,
		).
		AddRow(
			"mydb_log",
			"/var/opt/mssql/data/mydb_log.ldf",
			"L", nil,
			int64(2097152), int64(2199023255552), int64(2),
			int64(0), int64(0), "FFFFFFFF-1111-2222-3333-444444444444",
			int64(0), int64(0), int64(0),
			int64(512), int64(0), nil,
			int64(0), "00000000-0000-0000-0000-000000000000",
			false, true, nil,
		)

	mock.ExpectQuery("RESTORE FILELISTONLY").WillReturnRows(rows)

	entry, err := inspectBackupFile(db, "mydb.bak")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertEntry(t, entry, domain.BackupEntry{
		MdfFile: domain.MdfEntry{
			Name:       "mydb",
			Size:       "8388608",
			BackupSize: "1048576",
		},
		LdfFile: domain.LdfEntry{
			Name: "mydb_log",
			Size: "2097152",
		},
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled mock expectations: %v", err)
	}
}

func TestInspectBackupFile_EmptyResult(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	cols := []string{"LogicalName", "PhysicalName", "Type", "FileGroupName", "Size", "MaxSize", "FileId"}
	mock.ExpectQuery("RESTORE FILELISTONLY").WillReturnRows(sqlmock.NewRows(cols))

	entry, err := inspectBackupFile(db, "empty.bak")
	if err != nil {
		t.Fatalf("unexpected error on empty result: %v", err)
	}
	if entry.MdfFile.Name != "" || entry.LdfFile.Name != "" {
		t.Errorf("expected empty entry, got MDF=%q LDF=%q", entry.MdfFile.Name, entry.LdfFile.Name)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled mock expectations: %v", err)
	}
}

func assertEntry(t *testing.T, got, want domain.BackupEntry) {
	t.Helper()
	if got.MdfFile.Name != want.MdfFile.Name {
		t.Errorf("MdfFile.Name: got %q, want %q", got.MdfFile.Name, want.MdfFile.Name)
	}
	if got.MdfFile.Size != want.MdfFile.Size {
		t.Errorf("MdfFile.Size: got %q, want %q", got.MdfFile.Size, want.MdfFile.Size)
	}
	if got.MdfFile.BackupSize != want.MdfFile.BackupSize {
		t.Errorf("MdfFile.BackupSize: got %q, want %q", got.MdfFile.BackupSize, want.MdfFile.BackupSize)
	}
	if got.LdfFile.Name != want.LdfFile.Name {
		t.Errorf("LdfFile.Name: got %q, want %q", got.LdfFile.Name, want.LdfFile.Name)
	}
	if got.LdfFile.Size != want.LdfFile.Size {
		t.Errorf("LdfFile.Size: got %q, want %q", got.LdfFile.Size, want.LdfFile.Size)
	}
}
