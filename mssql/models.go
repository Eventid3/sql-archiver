/*
Package mssql defines data models for Microsoft SQL Server interactions.
*/
package mssql

type DBItem struct {
	Name, ID, Created, State string
}

type bakFile struct {
	size, date, name string
}

type backupEntry struct {
	mdfName, logName, size, backupsize string
}
