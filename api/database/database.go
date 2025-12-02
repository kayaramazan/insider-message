package database

import (
	"context"
	"database/sql"
)

type Database interface {
	Connect(ctx context.Context) error
	Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row
	Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Close() error
}
