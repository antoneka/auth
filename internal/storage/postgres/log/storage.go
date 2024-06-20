package log

import (
	"github.com/antoneka/platform-common/pkg/db"

	"github.com/antoneka/auth/internal/storage/postgres"
)

const (
	tableLogs = "logs"

	userIDColumn = "user_id"
	actionColumn = "action"
)

// log represents the log storage.
type log struct {
	db db.Client
}

// NewLogStorage creates a new instance of log storage.
func NewLogStorage(db db.Client) postgres.LogStorage {
	return &log{
		db: db,
	}
}
