package user

import (
	"github.com/antoneka/auth/internal/service"
	"github.com/antoneka/auth/internal/storage/postgres"
	"github.com/antoneka/auth/pkg/client/db"
)

var _ service.UserService = (*serv)(nil)

// serv represents the implementation of the UserService interface.
type serv struct {
	userStorage postgres.UserStorage
	logStorage  postgres.LogStorage
	txManager   db.TxManager
}

// NewService creates a new instance of the UserService interface.
func NewService(
	userStorage postgres.UserStorage,
	logStorage postgres.LogStorage,
	txManager db.TxManager,
) service.UserService {
	return &serv{
		userStorage: userStorage,
		logStorage:  logStorage,
		txManager:   txManager,
	}
}
