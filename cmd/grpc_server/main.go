package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/brianvoe/gofakeit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"

	desc "github.com/withsoull/auth-grpc/pkg/auth_v1"
)

var ErrPasswordDoNotMatch = errors.New("auth-grpc: passwords do not match")

const grpcPort = 50051

type server struct {
	desc.UnimplementedAuthV1Server
}

func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	log.Printf("User id: %d", req.GetId())

	return &desc.GetResponse{
		Id:    req.GetId(),
		Name:  gofakeit.Name(),
		Email: gofakeit.Email(),
		Role:  desc.Role(1),
	}, nil
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	id := int64(gofakeit.Number(1, 100))

	if req.GetPassword() != req.GetPasswordConfirm() {
		return nil, ErrPasswordDoNotMatch
	}

	// TODO: Add user to database

	fmt.Println("Was created user with:")
	fmt.Println("Name:\t", req.GetName())
	fmt.Println("Email:\t", req.GetEmail())
	fmt.Println("Password:\t", req.GetPassword())
	fmt.Println("Role:\t")

	log.Printf("User id: %d", id)
	return &desc.CreateResponse{
		Id: id,
	}, nil
}

func (s *server) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	if req.Name != nil {
		log.Printf("Updating user(ID: %d) name to: %s", req.GetId(), req.GetName())
	}
	if req.Email != nil {
		log.Printf("Updating user(ID: %d) email to: %s", req.GetId(), req.GetEmail())
	}

	return &emptypb.Empty{}, nil
}

func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	fmt.Printf("User(ID:%d) has been deleted", req.GetId())
	return &emptypb.Empty{}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterAuthV1Server(s, &server{})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
