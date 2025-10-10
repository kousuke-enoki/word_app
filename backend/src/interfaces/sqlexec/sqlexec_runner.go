// infrastructure/sqlexec/std_sql_runner.go
package sqlexec

import (
	"context"
	"database/sql"
)

type StdSQLRunner struct {
	db *sql.DB
}

func NewStdSQLRunner(db *sql.DB) *StdSQLRunner {
	return &StdSQLRunner{db: db}
}

func (r *StdSQLRunner) QueryRowContext(ctx context.Context, q string, args ...any) *sql.Row {
	return r.db.QueryRowContext(ctx, q, args...)
}

func (r *StdSQLRunner) ExecContext(ctx context.Context, q string, args ...any) (sql.Result, error) {
	return r.db.ExecContext(ctx, q, args...)
}

var _ Runner = (*StdSQLRunner)(nil)
