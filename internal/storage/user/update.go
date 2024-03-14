package user

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/antoneka/auth/internal/client/db"
	"github.com/antoneka/auth/internal/model"
)

// Update updates a user record in the database with the provided user information.
func (r *store) Update(
	ctx context.Context,
	user *model.User,
) error {
	const op = "storage.user.Update"

	builder := sq.Update(tableUsers).
		Set(nameColumn, user.UserInfo.Name).
		Set(emailColumn, user.UserInfo.Email).
		Set(passwordColumn, user.UserInfo.Password).
		Set(roleColumn, user.UserInfo.Role).
		Set(updatedColumn, user.UpdatedAt).
		Where(sq.Eq{idColumn: user.ID}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	q := db.Query{
		Name:     op,
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
