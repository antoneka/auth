package env

import (
	"errors"
	"os"
)

const (
	grpcPortEnvName = "GRPC_PORT"
)

// GRPCConfig represents the configuration for a gRPC server.
type GRPCConfig struct {
	Port string
}

// NewGRPCConfig creates a configuration for a gRPC server.
func NewGRPCConfig() (*GRPCConfig, error) {
	port := os.Getenv(grpcPortEnvName)
	if port == "" {
		return nil, errors.New("grpc port was not found")
	}

	return &GRPCConfig{
		Port: port,
	}, nil
}
