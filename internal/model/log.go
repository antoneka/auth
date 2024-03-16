package model

// Log actions constants.
const (
	LogActionCreateUser = "CREATE_USER"
	LogActionDeleteUser = "DELETE_USER"
	LogActionUpdateUser = "UPDATE_USER"
	LogActionGetUser    = "GET_USER"
)

// LogUser represents a log entry of the user`s action.
type LogUser struct {
	UserID int64
	Action string
}
