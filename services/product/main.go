package main

import (
	"context"
	"log"
	"net"

	pb "github.com/mrezayusufy/shop-api/pkg/proto/product"

	"google.golang.org/grpc"
)

const (
	port = ":50052"
)

// server is used to implement ecommerce.ProductServiceServer.
type server struct {
	pb.UnimplementedProductServiceServer
}

// GetProduct implements ecommerce.ProductServiceServer
func (s *server) GetProduct(ctx context.Context, in *pb.GetProductRequest) (*pb.ProductResponse, error) {
	log.Printf("Received: GetProduct request for ID: %s", in.GetId())
	// Mock data for demonstration
	product := &pb.ProductResponse{
		Id:          in.GetId(),
		Name:        "Product " + in.GetId(),
		Description: "A mock product description.",
		Price:       19.99,
		Stock: 			 12,
	}
	return &pb.ProductResponse{
		Id: product.Id,
		Name: product.Name,
		Price: product.Price,
		Stock: product.Stock,
		Description: product.Description,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterProductServiceServer(s, &server{})
	log.Printf("Product Service listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}