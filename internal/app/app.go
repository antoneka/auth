package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/antoneka/platform-common/pkg/closer"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rakyll/statik/fs"
	"github.com/rs/cors"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/antoneka/auth/internal/interceptor"
	desc "github.com/antoneka/auth/pkg/user_v1"

	// The statik package is used to embed static swagger server files into the Go binary.
	_ "github.com/antoneka/auth/statik"
)

// App represents the application instance.
type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
	httpServer      *http.Server
	swaggerServer   *http.Server
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
		a.runSwaggerServer,
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
		a.initSwaggerServer,
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
	a.grpcServer = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.UnaryInterceptor(interceptor.ValidateInterceptor),
	)

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

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Authorization"},
		AllowCredentials: true,
	})

	a.httpServer = &http.Server{
		Addr:              httpAddress,
		Handler:           corsMiddleware.Handler(mux),
		ReadHeaderTimeout: 2 * time.Second,
	}

	return nil
}

// initSwaggerServer initialized the swagger server.
func (a *App) initSwaggerServer(_ context.Context) error {
	statikFs, err := fs.New()
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.StripPrefix("/", http.FileServer(statikFs)))
	mux.HandleFunc("/user_api.swagger.json", serveSwaggerFile("/user_api.swagger.json"))

	swaggerAddress := fmt.Sprintf(":%s", a.serviceProvider.Config().Swagger.Port)

	a.swaggerServer = &http.Server{
		Addr:              swaggerAddress,
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

// runSwaggerServer starts the swagger server.
func (a *App) runSwaggerServer() error {
	log.Printf("Swagger server is running on localhost:%s", a.serviceProvider.Config().Swagger.Port)

	err := a.swaggerServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

// serveSwaggerFile returns an HTTP handler function that serves a Swagger JSON file
// located at the given path within the embedded filesystem.
func serveSwaggerFile(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		log.Printf("Serving swagger file: %s", path)

		statikFs, err := fs.New()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Open swagger file: %s", path)

		file, err := statikFs.Open(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = file.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Read swagger file: %s", path)

		content, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Write swagger file: %s", path)

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Served swagger file: %s", path)
	}
}
