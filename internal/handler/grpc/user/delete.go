package user

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/types/known/emptypb"

	desc "github.com/antoneka/auth/pkg/user_v1"
)

// Delete deletes the user from the system.
func (s *Implementation) Delete(
	ctx context.Context,
	req *desc.DeleteRequest,
) (*emptypb.Empty, error) {
	const op = "handler.grpc.user.Delete"

	id := req.GetId()

	err := s.userService.Delete(ctx, id)
	if err != nil {
		return &emptypb.Empty{}, fmt.Errorf("%s: %w", op, err)
	}

	return &emptypb.Empty{}, nil
}
