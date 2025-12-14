package main

import (
    "log"
    "net"

    "github.com/mrezayusufy/shop-api/internal/product"
    pb "github.com/mrezayusufy/shop-api/pkg/proto/product"
    "google.golang.org/grpc"
)

func main() {
    lis, err := net.Listen("tcp", ":50052")
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer()
    pb.RegisterProductServiceServer(grpcServer, product.NewService())

    log.Println("Product Service running on :50052")
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}