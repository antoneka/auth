package storage

import (
	"context"

	"github.com/antoneka/auth/internal/model"
)

// UserStorage defines the interface for interacting with user data in the storage.
type UserStorage interface {
	Create(ctx context.Context, info *model.UserInfo) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id int64) error
}
