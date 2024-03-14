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

	id, err := s.userStorage.Create(ctx, info)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
