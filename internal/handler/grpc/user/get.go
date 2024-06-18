package user

import (
	"context"
	"fmt"

	"github.com/antoneka/auth/internal/handler/grpc/user/converter"
	desc "github.com/antoneka/auth/pkg/user_v1"
)

// Get gets information about the user.
func (s *Implementation) Get(
	ctx context.Context,
	req *desc.GetRequest,
) (*desc.GetResponse, error) {
	const op = "handler.grpc.user.Get"

	id := req.GetId()

	user, err := s.userService.Get(ctx, id)
	if err != nil {
		return &desc.GetResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	return converter.ServiceToGetResponse(user), nil
}
