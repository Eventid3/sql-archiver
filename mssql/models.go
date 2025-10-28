/*
Package mssql defines data models for Microsoft SQL Server interactions.
*/
package mssql

type DBItem struct {
	Name, ID, Created, State string
}

type BakFile struct {
	Size, Date, Name string
}

type BackupEntry struct {
	MdfName, LogName, Size, Backupsize string
}
