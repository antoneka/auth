package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/antoneka/auth/internal/converter"
	"github.com/antoneka/auth/internal/service"
	userService "github.com/antoneka/auth/internal/service/user"
	userStorage "github.com/antoneka/auth/internal/storage/user"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/antoneka/auth/internal/config"
	desc "github.com/antoneka/auth/pkg/user_v1"
)

type server struct {
	desc.UnimplementedUserV1Server
	userService service.UserService
}

// Create creates a new user.
func (s *server) Create(
	ctx context.Context,
	req *desc.CreateRequest,
) (*desc.CreateResponse, error) {
	userInfo := converter.CreateRequestToService(req)

	id, err := s.userService.Create(ctx, userInfo)
	if err != nil {
		return &desc.CreateResponse{}, err
	}

	return &desc.CreateResponse{
		Id: id,
	}, nil
}

// Get gets information about the user.
func (s *server) Get(
	ctx context.Context,
	req *desc.GetRequest,
) (*desc.GetResponse, error) {
	id := req.GetId()

	user, err := s.userService.Get(ctx, id)
	if err != nil {
		return &desc.GetResponse{}, nil
	}

	return converter.ServiceToGetResponse(user), nil
}

// Update updates user information.
func (s *server) Update(
	ctx context.Context,
	req *desc.UpdateRequest,
) (*emptypb.Empty, error) {
	user := converter.UpdateRequestToService(req)

	err := s.userService.Update(ctx, user)
	if err != nil {
		return &emptypb.Empty{}, nil
	}

	return &emptypb.Empty{}, nil
}

// Delete deletes the user from the system.
func (s *server) Delete(
	ctx context.Context,
	req *desc.DeleteRequest,
) (*emptypb.Empty, error) {
	id := req.GetId()

	err := s.userService.Delete(ctx, id)
	if err != nil {
		return &emptypb.Empty{}, err
	}

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

	userStorage := userStorage.NewStorage(pool)
	userService := userService.NewService(userStorage)

	desc.RegisterUserV1Server(s, &server{
		userService: userService,
	})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
