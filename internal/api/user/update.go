package user

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/antoneka/auth/internal/converter"
	desc "github.com/antoneka/auth/pkg/user_v1"
)

// Update updates user information.
func (s *Implementation) Update(
	ctx context.Context,
	req *desc.UpdateRequest,
) (*emptypb.Empty, error) {
	user := converter.UpdateRequestToService(req)

	err := s.userService.Update(ctx, user)
	if err != nil {
		return &emptypb.Empty{}, nil
	}

	return &emptypb.Empty{}, nil
}
