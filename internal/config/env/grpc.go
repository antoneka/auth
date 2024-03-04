package env

import (
	"errors"
	"os"
)

const (
	grpcHostEnvName = "GRPC_HOST"
	grpcPortEnvName = "GRPC_PORT"
)

// GRPCConfig represents the configuration for a gRPC server.
type GRPCConfig struct {
	Host string
	Port string
}

// NewGRPCConfig creates a configuration for a gRPC server.
func NewGRPCConfig() (*GRPCConfig, error) {
	host := os.Getenv(grpcHostEnvName)
	if host == "" {
		return nil, errors.New("grpc host was not found")
	}

	port := os.Getenv(grpcPortEnvName)
	if port == "" {
		return nil, errors.New("grpc port was not found")
	}

	return &GRPCConfig{
		Host: host,
		Port: port,
	}, nil
}
