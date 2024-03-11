package user

import (
	"context"

	"github.com/antoneka/auth/internal/converter"
	desc "github.com/antoneka/auth/pkg/user_v1"
)

// Get gets information about the user.
func (s *Implementation) Get(
	ctx context.Context,
	req *desc.GetRequest,
) (*desc.GetResponse, error) {
	id := req.GetId()

	user, err := s.userService.Get(ctx, id)
	if err != nil {
		return &desc.GetResponse{}, nil
	}

	return converter.ServiceToGetResponse(user), nil
}
