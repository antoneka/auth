package user

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/antoneka/auth/internal/client/db"
	"github.com/antoneka/auth/internal/model"
	"github.com/antoneka/auth/internal/storage/user/converter"
	modelStore "github.com/antoneka/auth/internal/storage/user/model"
)

// Get retrieves a user record from the database based on the provided ID.
func (r *store) Get(
	ctx context.Context,
	id int64,
) (*model.User, error) {
	const op = "storage.user.Get"

	builder := sq.Select(idColumn, nameColumn, emailColumn, passwordColumn, roleColumn, createdColumn, updatedColumn).
		From(tableUsers).
		Where(sq.Eq{idColumn: id}).
		Limit(1).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	q := db.Query{
		Name:     op,
		QueryRaw: query,
	}

	var user modelStore.User
	err = r.db.DB().ScanOneContext(ctx, &user, q, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return converter.StorageToServiceUser(&user), nil
}
