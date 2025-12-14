package main

import (
    "log"
    "net"

    "github.com/mrezayusufy/shop-api/internal/order"
    pb "github.com/mrezayusufy/shop-api/pkg/proto/order"
    "google.golang.org/grpc"
)

func main() {
    lis, err := net.Listen("tcp", ":50053")
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer()
    pb.RegisterOrderServiceServer(grpcServer, order.NewService())

    log.Println("Order Service running on :50053")
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}