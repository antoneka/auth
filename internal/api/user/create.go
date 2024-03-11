package user

import (
	"context"

	"github.com/antoneka/auth/internal/converter"
	desc "github.com/antoneka/auth/pkg/user_v1"
)

// Create creates a new user.
func (s *Implementation) Create(
	ctx context.Context,
	req *desc.CreateRequest,
) (*desc.CreateResponse, error) {
	userInfo := converter.CreateRequestToService(req)

	id, err := s.userService.Create(ctx, userInfo)
	if err != nil {
		return &desc.CreateResponse{}, err
	}

	return &desc.CreateResponse{
		Id: id,
	}, nil
}
