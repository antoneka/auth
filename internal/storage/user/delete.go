package user

import (
	"context"

	sq "github.com/Masterminds/squirrel"
)

// Delete ...
func (r *store) Delete(
	ctx context.Context,
	id int64,
) error {
	builder := sq.Delete(tableUsers).
		Where(sq.Eq{idColumn: id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
