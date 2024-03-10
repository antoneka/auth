package model

import (
	"database/sql"
	"time"
)

// User ...
type User struct {
	ID        int64
	UserInfo  UserInfo
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

// UserInfo ...
type UserInfo struct {
	Name     string
	Email    string
	Password string
	Role     string
}
