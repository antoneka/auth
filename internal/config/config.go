package config

import (
	"flag"

	"github.com/antoneka/auth/internal/config/env"
	"github.com/joho/godotenv"
)

// Config represents the overall configuration for the app.
type Config struct {
	GRPC *env.GRPCConfig
	PG   *env.PGConfig
}

// MustLoad loads the configuration for the app from the .env file.
func MustLoad() *Config {
	var configPath string

	flag.StringVar(&configPath, "config", ".env", "path to config file")
	flag.Parse()

	err := godotenv.Load(configPath)
	if err != nil {
		panic("env file was not found")
	}

	grpcConfig, err := env.NewGRPCConfig()
	if err != nil {
		panic(err)
	}

	pgConfig, err := env.NewPGConfig()
	if err != nil {
		panic(err)
	}

	return &Config{
		GRPC: grpcConfig,
		PG:   pgConfig,
	}
}
