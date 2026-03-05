package mssql

import (
	"database/sql"
	"fmt"
	"net/url"

	_ "github.com/microsoft/go-mssqldb"
)

// openDB opens a connection to a SQL Server instance at host (e.g. "localhost:1433").
// The caller is responsible for closing the returned *sql.DB.
func openDB(host, user, password string) (*sql.DB, error) {
	if host == "" {
		host = "localhost:1433"
	}

	query := url.Values{}
	query.Add("TrustServerCertificate", "true")

	u := &url.URL{
		Scheme:   "sqlserver",
		User:     url.UserPassword(user, password),
		Host:     host,
		RawQuery: query.Encode(),
	}

	db, err := sql.Open("sqlserver", u.String())
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	return db, nil
}
