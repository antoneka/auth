package model

import (
	"database/sql"
	"time"
)

// Role ...
type Role string

const (
	// UNKNOWN ...
	UNKNOWN Role = "UNKNOWN_CHANGE_TYPE"
	// USER ...
	USER Role = "USER"
	// ADMIN ...
	ADMIN Role = "ADMIN"
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
	Role     Role
}
