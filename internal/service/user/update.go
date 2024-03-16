package user

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/antoneka/auth/internal/model"
)

// Update updates the user information based on the provided user object.
func (s *serv) Update(
	ctx context.Context,
	user *model.User,
) error {
	const op = "service.user.Update"

	err := s.txManager.ReadCommitted(ctx, func(context.Context) error {
		currentUser, errTx := s.userStorage.Get(ctx, user.ID)
		if errTx != nil {
			return errTx
		}

		if user.UserInfo.Name == "" {
			user.UserInfo.Name = currentUser.UserInfo.Name
		}
		if user.UserInfo.Email == "" {
			user.UserInfo.Email = currentUser.UserInfo.Email
		}
		if user.UserInfo.Password == "" {
			user.UserInfo.Password = currentUser.UserInfo.Password
		}
		if user.UserInfo.Role == model.UNKNOWN {
			user.UserInfo.Role = currentUser.UserInfo.Role
		}

		user.UpdatedAt = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}

		errTx = s.userStorage.Update(ctx, user)
		if errTx != nil {
			return errTx
		}

		errTx = s.logStorage.Log(ctx, &model.LogUser{
			UserID: user.ID,
			Action: model.LogActionUpdateUser,
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
