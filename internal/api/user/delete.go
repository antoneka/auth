package user

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	desc "github.com/antoneka/auth/pkg/user_v1"
)

// Delete deletes the user from the system.
func (s *Implementation) Delete(
	ctx context.Context,
	req *desc.DeleteRequest,
) (*emptypb.Empty, error) {
	id := req.GetId()

	err := s.userService.Delete(ctx, id)
	if err != nil {
		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, nil
}
