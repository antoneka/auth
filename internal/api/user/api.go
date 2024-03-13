package user

import (
	"github.com/antoneka/auth/internal/service"
	desc "github.com/antoneka/auth/pkg/user_v1"
)

// Implementation represents the implementation of the user API.
type Implementation struct {
	desc.UnimplementedUserV1Server
	userService service.UserService
}

// NewImplementation creates a new instance of the user API implementation.
func NewImplementation(userService service.UserService) *Implementation {
	return &Implementation{
		userService: userService,
	}
}
