package log

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/antoneka/auth/internal/model"
)

// Log inserts a log entry of the user`s action into the database.
func (l *log) Log(ctx context.Context, log *model.LogUser) error {
	const op = "storage.log.Log"

	builder := sq.Insert(tableLogs).
		Columns(userIDColumn, actionColumn).
		Values(log.UserID, log.Action).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = l.db.DB().Pool().Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
