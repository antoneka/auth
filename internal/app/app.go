package app

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/antoneka/auth/pkg/closer"
	desc "github.com/antoneka/auth/pkg/user_v1"
)

// App represents the application instance.
type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
}

// NewApp returns a new instance of the application.
func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

// Run starts the application.
func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	return a.runGRPCServer()
}

// initDeps initializes the application dependencies.
func (a *App) initDeps(ctx context.Context) error {
	inits := []func(ctx context.Context) error{
		a.initServiceProvider,
		a.initGRPCServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

// initServiceProvider initializes the service provider.
func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()

	return nil
}

// initGRPCServer initializes the gRPC server.
func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()))

	reflection.Register(a.grpcServer)

	desc.RegisterUserV1Server(a.grpcServer, a.serviceProvider.UserAPI(ctx))

	return nil
}

// runGRPCServer starts the gRPC server.
func (a *App) runGRPCServer() error {
	log.Printf("GRPC server is running on localhost:%s", a.serviceProvider.Config().GRPC.Port)

	list, err := net.Listen("tcp", fmt.Sprintf(":%s", a.serviceProvider.Config().GRPC.Port))
	if err != nil {
		return err
	}

	err = a.grpcServer.Serve(list)
	if err != nil {
		return err
	}

	return nil
}
