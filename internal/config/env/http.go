package env

import (
	"errors"
	"os"
)

const (
	httpPortEnvName = "HTTP_CONTAINER_PORT"
)

// HTTPConfig represents the configuration for a HTTP server.
type HTTPConfig struct {
	Port string
}

// NewHTTPConfig creates a configuration for a HTTP server.
func NewHTTPConfig() (*HTTPConfig, error) {
	port := os.Getenv(httpPortEnvName)
	if port == "" {
		return nil, errors.New("http port was not found")
	}

	return &HTTPConfig{
		Port: port,
	}, nil
}
