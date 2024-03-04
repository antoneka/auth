package env

import (
	"errors"
	"os"
)

const (
	dsnEnvName = "PG_DSN"
)

// PGConfig represents the configuration for a PostgreSQL database.
type PGConfig struct {
	DSN string
}

// NewPGConfig creates a configuration for a PostgreSQL database.
func NewPGConfig() (*PGConfig, error) {
	dsn := os.Getenv(dsnEnvName)
	if dsn == "" {
		return nil, errors.New("postgres dsn was not found")
	}

	return &PGConfig{
		DSN: dsn,
	}, nil
}
