package user

import (
	"context"
	"fmt"

	"github.com/antoneka/auth/internal/model"
)

// Delete deletes the user with the specified ID from the user storage.
func (s *serv) Delete(
	ctx context.Context,
	id int64,
) error {
	const op = "service.user.Delete"

	err := s.txManager.ReadCommitted(ctx, func(context.Context) error {
		errTx := s.userStorage.Delete(ctx, id)
		if errTx != nil {
			return errTx
		}

		errTx = s.logStorage.Log(ctx, &model.LogUser{
			UserID: id,
			Action: model.LogActionDeleteUser,
		})
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
