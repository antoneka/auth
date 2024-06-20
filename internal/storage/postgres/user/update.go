package user

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"github.com/antoneka/platform-common/pkg/db"

	"github.com/antoneka/auth/internal/model"
	"github.com/antoneka/auth/internal/storage/postgres/user/converter"
)

// Update updates a user record in the database with the provided user information.
func (r *store) Update(
	ctx context.Context,
	user *model.User,
) error {
	const op = "storage.user.Update"

	storeUser := converter.ServiceUserToStorage(user)

	builder := sq.Update(tableUsers).
		Set(nameColumn, storeUser.UserInfo.Name).
		Set(emailColumn, storeUser.UserInfo.Email).
		Set(passwordColumn, storeUser.UserInfo.Password).
		Set(roleColumn, storeUser.UserInfo.Role).
		Set(updatedColumn, storeUser.UpdatedAt).
		Where(sq.Eq{idColumn: storeUser.ID}).
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
