package model

import (
	"database/sql"
	"time"
)

// User represents a user entity in the database.
type User struct {
	ID        int64        `db:"id"`
	UserInfo  UserInfo     `db:""`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}

// UserInfo represents detailed information about a user.
type UserInfo struct {
	Name     string `db:"name"`
	Email    string `db:"email"`
	Password string `db:"password"`
}
