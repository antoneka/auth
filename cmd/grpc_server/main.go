package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/antoneka/auth/internal/model"
	"github.com/antoneka/auth/internal/storage"
	"github.com/antoneka/auth/internal/storage/user"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/antoneka/auth/internal/config"
	desc "github.com/antoneka/auth/pkg/user_v1"
)

type server struct {
	desc.UnimplementedUserV1Server
	userStorage storage.UserStorage
}

// Create creates a new user.
func (s *server) Create(
	ctx context.Context,
	req *desc.CreateRequest,
) (*desc.CreateResponse, error) {
	userInfo := &model.UserInfo{
		Name:     req.GetName(),
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
		Role:     desc.Role_name[int32(req.Role)],
	}

	id, err := s.userStorage.Create(ctx, userInfo)
	if err != nil {
		log.Fatal(err)
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

	user, err := s.userStorage.Get(ctx, id)
	if err != nil {
		log.Fatal(err)
	}

	var updatedAt *timestamppb.Timestamp
	if user.UpdatedAt.Valid {
		updatedAt = timestamppb.New(user.UpdatedAt.Time)
	}

	return &desc.GetResponse{
		Id:        user.ID,
		Name:      user.UserInfo.Name,
		Email:     user.UserInfo.Email,
		Role:      desc.Role(desc.Role_value[user.UserInfo.Role]),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: updatedAt,
	}, nil
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
	desc.RegisterUserV1Server(s, &server{userStorage: user.NewStorage(pool)})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
