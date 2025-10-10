// interfaces/sqlexec/sql_runner.go
package sqlexec

import (
	"context"
	"database/sql"
)

type Runner interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}
