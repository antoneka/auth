package model

import (
	"database/sql"
	"time"
)

// User represents a user entity.
type User struct {
	ID        int64
	UserInfo  UserInfo
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

// UserInfo contains user-specific information.
type UserInfo struct {
	Name     string
	Email    string
	Password string
}
