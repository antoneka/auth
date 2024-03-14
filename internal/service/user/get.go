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

	user, err := s.userStorage.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}
