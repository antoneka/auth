package env

import (
	"os"

	"errors"
)

const (
	swaggerPortEnvName = "SWAGGER_CONTAINER_PORT"
)

// SwaggerConfig represents the configuration for a swagger server.
type SwaggerConfig struct {
	Port string
}

// NewSwaggerConfig creates a configuration for a swagger server.
func NewSwaggerConfig() (*SwaggerConfig, error) {
	port := os.Getenv(swaggerPortEnvName)
	if port == "" {
		return nil, errors.New("swagger port was not found")
	}

	return &SwaggerConfig{
		Port: port,
	}, nil
}
