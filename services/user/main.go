package main

import (
	"context"
	"log"
	"net"

	pb "github.com/mrezayusufy/shop-api/pkg/proto/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const port = ":50051"

type userServer struct {
	pb.UnimplementedUserServiceServer
	users map[string]*pb.User
}

func newUserServer() *userServer {
	return &userServer{
		users: map[string]*pb.User{
			"user-1": {Id: "user-1", Username: "alice", Email: "alice@example.com"},
			"user-2": {Id: "user-2", Username: "bob", Email: "bob@example.com"},
		},
	}
}

func (s *userServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, ok := s.users[req.GetId()]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "user %s not found", req.GetId())
	}
	return &pb.GetUserResponse{User: user}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	pb.RegisterUserServiceServer(srv, newUserServer())

	log.Printf("User Service listening at %s", port)
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
