package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	userAPI "github.com/antoneka/auth/internal/api/user"
	userService "github.com/antoneka/auth/internal/service/user"
	userStorage "github.com/antoneka/auth/internal/storage/user"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/antoneka/auth/internal/config"
	desc "github.com/antoneka/auth/pkg/user_v1"
)

func main() {
	ctx := context.Background()
	cfg := config.MustLoad()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPC.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	pool, err := pgxpool.Connect(ctx, cfg.PG.DSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	s := grpc.NewServer()
	reflection.Register(s)

	userStorage := userStorage.NewStorage(pool)
	userService := userService.NewService(userStorage)

	desc.RegisterUserV1Server(s, userAPI.NewImplementation(userService))

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
