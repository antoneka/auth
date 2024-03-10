package user

import (
	"context"
	"fmt"
)

// Delete ...
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
