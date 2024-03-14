package service

import (
	"context"

	"github.com/antoneka/auth/internal/model"
)

// UserService defines the interface for user management operations within the service layer.
type UserService interface {
	Create(ctx context.Context, info *model.UserInfo) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id int64) error
}
