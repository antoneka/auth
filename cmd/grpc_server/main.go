package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/antoneka/auth/internal/config"
	desc "github.com/antoneka/auth/pkg/user_v1"
)

type server struct {
	desc.UnimplementedUserV1Server
	pool *pgxpool.Pool
}

// Create creates a new user.
func (s *server) Create(
	ctx context.Context,
	req *desc.CreateRequest,
) (*desc.CreateResponse, error) {
	_, cancel := context.WithCancel(ctx)
	defer cancel()

	fmt.Printf("%+v\n", req)

	return &desc.CreateResponse{}, nil
}

// Get gets information about the user.
func (s *server) Get(
	ctx context.Context,
	req *desc.GetRequest,
) (*desc.GetResponse, error) {
	_, cancel := context.WithCancel(ctx)
	defer cancel()

	fmt.Printf("+%v\n", req)

	return &desc.GetResponse{}, nil
}

// Update updates user information.
func (s *server) Update(
	ctx context.Context,
	req *desc.UpdateRequest,
) (*emptypb.Empty, error) {
	_, cancel := context.WithCancel(ctx)
	defer cancel()

	fmt.Printf("+%v\n", req)

	return &emptypb.Empty{}, nil
}

// Delete deletes the user from the system.
func (s *server) Delete(
	ctx context.Context,
	req *desc.DeleteRequest,
) (*emptypb.Empty, error) {
	_, cancel := context.WithCancel(ctx)
	defer cancel()

	fmt.Printf("%+v\n", req)

	return &emptypb.Empty{}, nil
}

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
	desc.RegisterUserV1Server(s, &server{pool: pool})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
