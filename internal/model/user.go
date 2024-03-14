package model

import (
	"database/sql"
	"time"
)

// Role represents the role of a user.
type Role string

const (
	// UNKNOWN is a stub in the role enum.
	UNKNOWN Role = "UNKNOWN_CHANGE_TYPE"
	// USER represents a user role.
	USER Role = "USER"
	// ADMIN represents an admin role.
	ADMIN Role = "ADMIN"
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
	Role     Role
}
