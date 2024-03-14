package user

import (
	"github.com/antoneka/auth/internal/service"
	"github.com/antoneka/auth/internal/storage"
)

var _ service.UserService = (*serv)(nil)

type serv struct {
	userStorage storage.UserStorage
}

// NewService creates a new instance of the UserService interface.
func NewService(userStorage storage.UserStorage) service.UserService {
	return &serv{
		userStorage: userStorage,
	}
}
