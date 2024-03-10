package user

import (
	"github.com/jackc/pgx/v4/pgxpool"

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

type store struct {
	db *pgxpool.Pool
}

// NewStorage ...
func NewStorage(db *pgxpool.Pool) storage.UserStorage {
	return &store{db: db}
}
