package pg

import (
	"fmt"
	"log"

	"golang.org/x/net/context"

	"github.com/antoneka/auth/internal/client/db"
	"github.com/antoneka/auth/internal/client/db/prettier"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type pg struct {
	dbc *pgxpool.Pool
}

// NewDB creates a new PostgreSQL database instance.
func NewDB(dbc *pgxpool.Pool) db.DB {
	return &pg{
		dbc: dbc,
	}
}

// ScanOneContext executes a query and scans a single row into the destination struct.
func (p *pg) ScanOneContext(ctx context.Context, dest interface{}, q db.Query, args ...interface{}) error {
	row, err := p.QueryContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return pgxscan.ScanOne(dest, row)
}

// ScanAllContext executes a query and scans all rows into the destination slice.
func (p *pg) ScanAllContext(ctx context.Context, dest interface{}, q db.Query, args ...interface{}) error {
	rows, err := p.QueryContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return pgxscan.ScanAll(dest, rows)
}

// ExecContext executes a query that doesn't return rows.
func (p *pg) ExecContext(ctx context.Context, q db.Query, args ...interface{}) (pgconn.CommandTag, error) {
	logQuery(ctx, q, args...)

	return p.dbc.Exec(ctx, q.QueryRaw, args...)
}

// QueryContext executes a query that returns rows.
func (p *pg) QueryContext(ctx context.Context, q db.Query, args ...interface{}) (pgx.Rows, error) {
	logQuery(ctx, q, args...)

	return p.dbc.Query(ctx, q.QueryRaw, args...)
}

// QueryRowContext executes a query that is expected to return at most one row.
func (p *pg) QueryRowContext(ctx context.Context, q db.Query, args ...interface{}) pgx.Row {
	logQuery(ctx, q, args...)

	return p.dbc.QueryRow(ctx, q.QueryRaw, args...)
}

// Ping checks if the database connection is alive.
func (p *pg) Ping(ctx context.Context) error {
	return p.dbc.Ping(ctx)
}

// Close closes the database connection.
func (p *pg) Close() {
	p.dbc.Close()
}

// logQuery logs the SQL query.
func logQuery(ctx context.Context, q db.Query, args ...interface{}) {
	prettyQuery := prettier.Pretty(q.QueryRaw, prettier.PlaceholderDollar, args...)
	log.Println(
		ctx,
		fmt.Sprintf("sql: %s", q.Name),
		fmt.Sprintf("query: %s", prettyQuery),
	)
}
