package user

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/antoneka/auth/internal/model"
	"github.com/antoneka/auth/internal/storage/user/converter"
	modelStore "github.com/antoneka/auth/internal/storage/user/model"
)

// Get ...
func (r *store) Get(
	ctx context.Context,
	id int64,
) (*model.User, error) {
	builder := sq.Select(idColumn, nameColumn, emailColumn, passwordColumn, roleColumn, createdColumn, updatedColumn).
		From(tableUsers).
		Where(sq.Eq{idColumn: id}).
		Limit(1).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	var user modelStore.User
	err = r.db.QueryRow(ctx, query, args...).Scan(&user.ID, &user.UserInfo.Name, &user.UserInfo.Email, &user.UserInfo.Password, &user.UserInfo.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return converter.StorageToServiceUser(&user), nil
}
