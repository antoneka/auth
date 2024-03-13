package app

import (
	"context"
	"log"

	"github.com/antoneka/auth/internal/api/user"
	"github.com/antoneka/auth/internal/closer"
	"github.com/antoneka/auth/internal/config"
	"github.com/antoneka/auth/internal/service"
	userServ "github.com/antoneka/auth/internal/service/user"
	"github.com/antoneka/auth/internal/storage"
	userStore "github.com/antoneka/auth/internal/storage/user"
	"github.com/jackc/pgx/v4/pgxpool"
)

type serviceProvider struct {
	config *config.Config

	pgPool      *pgxpool.Pool
	userStorage storage.UserStorage

	userService service.UserService

	userAPI *user.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

// Config returns the application configuration.
func (s *serviceProvider) Config() *config.Config {
	if s.config == nil {
		s.config = config.MustLoad()
	}

	return s.config
}

// PgPool returns the database connection pool.
func (s *serviceProvider) PgPool(ctx context.Context) *pgxpool.Pool {
	if s.pgPool == nil {
		pool, err := pgxpool.Connect(ctx, s.Config().PG.DSN)
		if err != nil {
			log.Panicf("failed to connect to database: %v", err)
		}

		err = pool.Ping(ctx)
		if err != nil {
			log.Panicf("ping error: %v", err)
		}

		closer.Add(func() error {
			pool.Close()
			return nil
		})

		s.pgPool = pool
	}

	return s.pgPool
}

// UserStorage returns the user storage instance.
func (s *serviceProvider) UserStorage(ctx context.Context) storage.UserStorage {
	if s.userStorage == nil {
		s.userStorage = userStore.NewStorage(s.PgPool(ctx))
	}

	return s.userStorage
}

// UserService returns the user service instance.
func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = userServ.NewService(s.UserStorage(ctx))
	}

	return s.userService
}

// UserAPI returns the user API implementation.
func (s *serviceProvider) UserAPI(ctx context.Context) *user.Implementation {
	if s.userAPI == nil {
		s.userAPI = user.NewImplementation(s.UserService(ctx))
	}

	return s.userAPI
}
