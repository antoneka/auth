package config

import (
	"flag"
	"log"

	"github.com/joho/godotenv"

	"github.com/antoneka/auth/internal/config/env"
)

// Config represents the overall configuration for the app.
type Config struct {
	GRPC    *env.GRPCConfig
	HTTP    *env.HTTPConfig
	PG      *env.PGConfig
	Swagger *env.SwaggerConfig
	Redis   *env.RedisConfig
}

// MustLoad loads the configuration for the app from the .env file.
func MustLoad() *Config {
	var configPath string

	flag.StringVar(&configPath, "config", ".env", "path to config file")
	flag.Parse()

	err := godotenv.Load(configPath)
	if err != nil {
		log.Panicf("failed to load .env file: %v", err)
	}

	grpcConfig, err := env.NewGRPCConfig()
	if err != nil {
		log.Panicf("failed to load gRPC config: %v", err)
	}

	pgConfig, err := env.NewPGConfig()
	if err != nil {
		log.Panicf("failed to load PostgreSQL config: %v", err)
	}

	httpConfig, err := env.NewHTTPConfig()
	if err != nil {
		log.Panicf("failed to load HTTP config: %v", err)
	}

	swaggerConfig, err := env.NewSwaggerConfig()
	if err != nil {
		log.Panicf("failed to load Swagger config: %v", err)
	}

	redisConfig, err := env.NewRedisConfig()
	if err != nil {
		log.Panicf("failed to load Redis config: %v", err)
	}

	return &Config{
		GRPC:    grpcConfig,
		HTTP:    httpConfig,
		PG:      pgConfig,
		Swagger: swaggerConfig,
		Redis:   redisConfig,
	}
}
