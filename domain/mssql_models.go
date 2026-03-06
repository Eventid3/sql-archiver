/*
Package domain
*/
package domain

// Connection holds a saved SQL Server connection profile.
type Connection struct {
	Name     string
	Host     string
	User     string
	Password string
}

type DBItem struct {
	Name, ID, Created, State string
}

type BakFile struct {
	Size, Date, Name string
}

type BackupEntry struct {
	Filename string
	MdfFile  MdfEntry
	LdfFile  LdfEntry
}

type MdfEntry struct {
	Name, Size, BackupSize string
}

type LdfEntry struct {
	Name, Size string
}
