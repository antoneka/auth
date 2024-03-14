package db

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

// Client represents a client for interacting with the database.
type Client interface {
	DB() DB
	Close() error
}

// Query represents a database query.
type Query struct {
	Name     string
	QueryRaw string
}

// SQLExecer represents an interface for executing SQL queries.
type SQLExecer interface {
	NamedExecer
	QueryExecer
}

// NamedExecer represents an interface for executing named SQL queries using tags in structures.
type NamedExecer interface {
	ScanOneContext(ctx context.Context, dest interface{}, q Query, args ...interface{}) error
	ScanAllContext(ctx context.Context, dest interface{}, q Query, args ...interface{}) error
}

// QueryExecer represents an interface for executing SQL queries.
type QueryExecer interface {
	ExecContext(ctx context.Context, q Query, args ...interface{}) (pgconn.CommandTag, error)
	QueryContext(ctx context.Context, q Query, args ...interface{}) (pgx.Rows, error)
	QueryRowContext(ctx context.Context, q Query, args ...interface{}) pgx.Row
}

// Pinger represents an interface for checking database connectivity.
type Pinger interface {
	Ping(ctx context.Context) error
}

// DB represents a database connection.
type DB interface {
	SQLExecer
	Pinger
	Close()
}
