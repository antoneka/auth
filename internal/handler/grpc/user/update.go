package user

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/antoneka/auth/internal/handler/grpc/user/converter"
	desc "github.com/antoneka/auth/pkg/user_v1"
)

// Update updates user information.
func (s *Implementation) Update(
	ctx context.Context,
	req *desc.UpdateRequest,
) (*emptypb.Empty, error) {
	const op = "handler.grpc.user.Update"

	user := converter.UpdateRequestToService(req)

	err := s.userService.Update(ctx, user)
	if err != nil {
		return &emptypb.Empty{}, fmt.Errorf("%s: %w", op, err)
	}

	return &emptypb.Empty{}, nil
}
