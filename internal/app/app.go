package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/antoneka/platform-common/pkg/closer"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	desc "github.com/antoneka/auth/pkg/user_v1"
)

// App represents the application instance.
type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
	httpServer      *http.Server
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

	eg := new(errgroup.Group)

	runActions := []func() error{
		a.runGRPCServer,
		a.runHTTPServer,
	}

	for _, runAction := range runActions {
		eg.Go(runAction)
	}

	return eg.Wait()
}

// initDeps initializes the application dependencies.
func (a *App) initDeps(ctx context.Context) error {
	inits := []func(ctx context.Context) error{
		a.initServiceProvider,
		a.initGRPCServer,
		a.initHTTPServer,
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

// initHTTPServer initializes the HTTP server.
func (a *App) initHTTPServer(ctx context.Context) error {
	mux := runtime.NewServeMux()

	grpcAddress := fmt.Sprintf(":%s", a.serviceProvider.Config().GRPC.Port)
	httpAddress := fmt.Sprintf(":%s", a.serviceProvider.Config().HTTP.Port)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	err := desc.RegisterUserV1HandlerFromEndpoint(ctx, mux, grpcAddress, opts)
	if err != nil {
		return nil
	}

	a.httpServer = &http.Server{
		Addr:              httpAddress,
		Handler:           mux,
		ReadHeaderTimeout: 2 * time.Second,
	}

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

// runHTTPServer starts the HTTP server.
func (a *App) runHTTPServer() error {
	log.Printf("HTTP server is running on localhost:%s", a.serviceProvider.Config().HTTP.Port)

	err := a.httpServer.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
