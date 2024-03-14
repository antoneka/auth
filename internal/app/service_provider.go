package app

import (
	"context"
	"log"

	"github.com/antoneka/auth/internal/api/user"
	"github.com/antoneka/auth/internal/client/db"
	"github.com/antoneka/auth/internal/client/db/pg"
	"github.com/antoneka/auth/internal/closer"
	"github.com/antoneka/auth/internal/config"
	"github.com/antoneka/auth/internal/service"
	userServ "github.com/antoneka/auth/internal/service/user"
	"github.com/antoneka/auth/internal/storage"
	userStore "github.com/antoneka/auth/internal/storage/user"
)

type serviceProvider struct {
	config *config.Config

	dbClient    db.Client
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

// DBClient returns the database client.
func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		client, err := pg.New(ctx, s.Config().PG.DSN)
		if err != nil {
			log.Panicf("failed to create db client %v", err)
		}

		err = client.DB().Ping(ctx)
		if err != nil {
			log.Panicf("ping error: %v", err)
		}

		closer.Add(client.Close)

		s.dbClient = client
	}

	return s.dbClient
}

// UserStorage returns the user storage instance.
func (s *serviceProvider) UserStorage(ctx context.Context) storage.UserStorage {
	if s.userStorage == nil {
		s.userStorage = userStore.NewStorage(s.DBClient(ctx))
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
