package user

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/antoneka/auth/internal/model"
)

// Create ...
func (r *store) Create(
	ctx context.Context,
	info *model.UserInfo,
) (int64, error) {
	builder := sq.Insert(tableUsers).
		Columns(nameColumn, emailColumn, passwordColumn, roleColumn).
		Values(info.Name, info.Email, info.Password, info.Role).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}

	var id int64
	err = r.db.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
