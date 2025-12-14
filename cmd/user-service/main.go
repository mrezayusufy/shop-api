package main

import (
    "log"
    "net"

    "github.com/mrezayusufy/shop-api/internal/user"
    pb "github.com/mrezayusufy/shop-api/pkg/proto/user"
    "google.golang.org/grpc"
)

func main() {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer()
    pb.RegisterUserServiceServer(grpcServer, user.NewService())

    log.Println("User Service running on :50051")
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}