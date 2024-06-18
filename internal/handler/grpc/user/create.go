package user

import (
	"context"
	"fmt"

	"github.com/antoneka/auth/internal/handler/grpc/user/converter"
	desc "github.com/antoneka/auth/pkg/user_v1"
)

// Create creates a new user.
func (s *Implementation) Create(
	ctx context.Context,
	req *desc.CreateRequest,
) (*desc.CreateResponse, error) {
	const op = "handler.grpc.user.Create"

	userInfo := converter.CreateRequestToService(req)

	id, err := s.userService.Create(ctx, userInfo)
	if err != nil {
		return &desc.CreateResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	return &desc.CreateResponse{
		Id: id,
	}, nil
}
