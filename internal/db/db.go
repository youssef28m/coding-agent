// internal/db/db.go
package db

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func NewSQLiteDB(path string) (*sql.DB, error) {
	conn, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	conn.SetMaxOpenConns(1)
	if _, err := conn.Exec("PRAGMA foreign_keys=ON;"); err != nil {
		return nil, err
	}
	return conn, conn.Ping()
}