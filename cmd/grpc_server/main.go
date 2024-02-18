package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"

	desc "github.com/antoneka/auth/pkg/user_v1"
)

const grpcPort = 50051

type server struct {
	desc.UnimplementedUserV1Server
}

func (s *server) Create(
	ctx context.Context,
	req *desc.CreateRequest,
) (*desc.CreateResponse, error) {
	_, cancel := context.WithCancel(ctx)
	defer cancel()

	fmt.Printf("%+v\n", req)

	return &desc.CreateResponse{}, nil
}

func (s *server) Get(
	ctx context.Context,
	req *desc.GetRequest,
) (*desc.GetResponse, error) {
	_, cancel := context.WithCancel(ctx)
	defer cancel()

	fmt.Printf("+%v\n", req)

	return &desc.GetResponse{}, nil
}

func (s *server) Update(
	ctx context.Context,
	req *desc.UpdateRequest,
) (*emptypb.Empty, error) {
	_, cancel := context.WithCancel(ctx)
	defer cancel()

	fmt.Printf("+%v\n", req)

	return &emptypb.Empty{}, nil
}

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
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterUserV1Server(s, &server{})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
