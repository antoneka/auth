package user

import (
	"context"
	"fmt"
)

// Delete deletes the user with the specified ID from the user storage.
func (s *serv) Delete(
	ctx context.Context,
	id int64,
) error {
	const op = "service.user.Delete"

	err := s.userStorage.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
