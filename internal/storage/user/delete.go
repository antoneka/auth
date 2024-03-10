package user

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

// Delete ...
func (r *store) Delete(
	ctx context.Context,
	id int64,
) error {
	const op = "storage.user.Delete"

	builder := sq.Delete(tableUsers).
		Where(sq.Eq{idColumn: id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
