package user

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/antoneka/auth/internal/model"
)

// Update ...
func (r *store) Update(
	ctx context.Context,
	user *model.User,
) error {
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
		return err
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
