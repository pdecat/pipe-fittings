package backend

import (
	"context"
	"database/sql"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v5/sperr"
)

const sqliteConnectionStringPrefix = "sqlite:"

type SqliteBackend struct {
	connectionString string
	rowreader        RowReader
}

func NewSqliteBackend(connString string) *SqliteBackend {
	connString = strings.TrimSpace(connString) // remove any leading or trailing whitespace
	connString = strings.TrimPrefix(connString, sqliteConnectionStringPrefix)
	return &SqliteBackend{
		connectionString: connString,
		rowreader:        NewSqliteRowReader(),
	}
}

// Connect implements Backend.
func (s *SqliteBackend) Connect(_ context.Context, options ...ConnectOption) (*sql.DB, error) {
	config := NewConnectConfig(options)
	db, err := sql.Open("sqlite3", s.connectionString)
	if err != nil {
		return nil, sperr.WrapWithMessage(err, "could not connect to sqlite backend")
	}
	db.SetConnMaxIdleTime(config.MaxConnIdleTime)
	db.SetConnMaxLifetime(config.MaxConnLifeTime)
	db.SetMaxOpenConns(config.MaxOpenConns)
	return db, nil
}

// RowReader implements Backend.
func (s *SqliteBackend) RowReader() RowReader {
	return s.rowreader
}

type sqliteRowReader struct {
	BasicRowReader
}

func NewSqliteRowReader() *sqliteRowReader {
	return &sqliteRowReader{
		// use the generic row reader - there's no real difference between sqlite and generic
		BasicRowReader: *NewBasicRowReader(),
	}
}
