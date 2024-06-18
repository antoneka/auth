package app

import (
	"context"
	"log"

	"github.com/antoneka/auth/internal/config"
	"github.com/antoneka/auth/internal/handler/grpc/user"
	"github.com/antoneka/auth/internal/service"
	userServ "github.com/antoneka/auth/internal/service/user"
	"github.com/antoneka/auth/internal/storage/postgres"
	logStore "github.com/antoneka/auth/internal/storage/postgres/log"
	userStore "github.com/antoneka/auth/internal/storage/postgres/user"
	"github.com/antoneka/auth/pkg/client/db"
	"github.com/antoneka/auth/pkg/client/db/pg"
	"github.com/antoneka/auth/pkg/client/db/transaction"
	"github.com/antoneka/auth/pkg/closer"
)

// serviceProvider is a DI container that manages service dependencies.
type serviceProvider struct {
	config *config.Config

	dbClient    db.Client
	txManager   db.TxManager
	logStorage  postgres.LogStorage
	userStorage postgres.UserStorage

	userService service.UserService

	userAPI *user.Implementation
}

// newServiceProvider creates a new instance of serviceProvider.
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
		client, err := pg.NewDBClient(ctx, s.Config().PG.DSN)
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

// TxManager returns the transaction manager.
func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) LogStorage(ctx context.Context) postgres.LogStorage {
	if s.logStorage == nil {
		s.logStorage = logStore.NewLogStorage(s.DBClient(ctx))
	}

	return s.logStorage
}

// UserStorage returns the user storage instance.
func (s *serviceProvider) UserStorage(ctx context.Context) postgres.UserStorage {
	if s.userStorage == nil {
		s.userStorage = userStore.NewStorage(s.DBClient(ctx))
	}

	return s.userStorage
}

// UserService returns the user service instance.
func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = userServ.NewService(
			s.UserStorage(ctx),
			s.LogStorage(ctx),
			s.TxManager(ctx),
		)
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
