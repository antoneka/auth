package user

import (
	"context"
	"fmt"

	"github.com/antoneka/auth/internal/model"
)

// Create creates a new user based on the provided user information.
func (s *serv) Create(
	ctx context.Context,
	info *model.UserInfo,
) (int64, error) {
	const op = "service.user.Create"

	var id int64

	err := s.txManager.ReadCommitted(ctx, func(context.Context) error {
		var errTx error
		id, errTx = s.userStorage.Create(ctx, info)
		if errTx != nil {
			return errTx
		}

		errTx = s.logStorage.Log(ctx, &model.LogUser{
			UserID: id,
			Action: model.LogActionCreateUser,
		})
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
