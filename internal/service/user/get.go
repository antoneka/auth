package user

import (
	"context"
	"fmt"

	"github.com/antoneka/auth/internal/model"
)

// Get retrieves the user with the specified ID from the user storage.
func (s *serv) Get(
	ctx context.Context,
	id int64,
) (*model.User, error) {
	const op = "service.user.Get"

	var user *model.User

	err := s.txManager.ReadCommitted(ctx, func(context.Context) error {
		var errTx error
		user, errTx = s.userStorage.Get(ctx, id)
		if errTx != nil {
			return errTx
		}

		errTx = s.logStorage.Log(ctx, &model.LogUser{
			UserID: id,
			Action: model.LogActionGetUser,
		})
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}
