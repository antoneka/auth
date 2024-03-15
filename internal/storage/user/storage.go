package user

import (
	"github.com/antoneka/auth/internal/client/db"
	"github.com/antoneka/auth/internal/storage"
)

const (
	tableUsers = "users"

	idColumn       = "id"
	nameColumn     = "name"
	emailColumn    = "email"
	passwordColumn = "password"
	roleColumn     = "role"
	createdColumn  = "created_at"
	updatedColumn  = "updated_at"
)

var _ storage.UserStorage = (*store)(nil)

// store represents the implementation of the UserStorage interface.
type store struct {
	db db.Client
}

// NewStorage creates a new instance of UserStorage.
func NewStorage(db db.Client) storage.UserStorage {
	return &store{db: db}
}
