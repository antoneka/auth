package user

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"github.com/antoneka/platform-common/pkg/db"

	"github.com/antoneka/auth/internal/model"
	"github.com/antoneka/auth/internal/storage/postgres/user/converter"
)

// Create creates a new user record in the database.
func (r *store) Create(
	ctx context.Context,
	info *model.UserInfo,
) (int64, error) {
	const op = "storage.user.Create"

	storeUserInfo := converter.ServiceUserInfoToStorage(info)

	builder := sq.Insert(tableUsers).
		Columns(nameColumn, emailColumn, passwordColumn).
		Values(storeUserInfo.Name, storeUserInfo.Email, storeUserInfo.Password).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	q := db.Query{
		Name:     op,
		QueryRaw: query,
	}

	var id int64
	err = r.db.DB().ScanOneContext(ctx, &id, q, args...)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
