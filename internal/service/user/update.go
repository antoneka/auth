package user

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/antoneka/auth/internal/model"
)

// Update ...
func (s *serv) Update(
	ctx context.Context,
	user *model.User,
) error {
	const op = "service.user.Update"

	currentUser, err := s.userStorage.Get(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
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

	err = s.userStorage.Update(ctx, user)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
