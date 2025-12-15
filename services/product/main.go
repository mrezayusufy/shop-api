package main

import (
	"context"
	"log"
	"net"

	"github.com/google/uuid"

	pb "github.com/mrezayusufy/shop-api/pkg/proto/product"
	"google.golang.org/grpc"
)

const port = ":50052"

type productServer struct {
	pb.UnimplementedProductServiceServer
	products map[string]*pb.Product
}

func newProductServer() *productServer {
	products := map[string]*pb.Product{
		"product-1": {Id: "product-1", Name: "Wireless Mouse", Description: "A smooth wireless mouse", Price: 19.99, Stock: 42},
		"product-2": {Id: "product-2", Name: "Mechanical Keyboard", Description: "Backlit keyboard", Price: 39.99, Stock: 18},
		"product-3": {Id: "product-3", Name: "USB-C Cable", Description: "1m braided cable", Price: 9.99, Stock: 120},
	}
	return &productServer{products: products}
}

func (s *productServer) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.ProductResponse, error) {
	product, ok := s.products[req.GetId()]
	if !ok {
		return &pb.ProductResponse{}, nil
	}
	return &pb.ProductResponse{Product: product}, nil
}

func (s *productServer) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.ProductResponse, error) {
	id := uuid.NewString()
	product := &pb.Product{
		Id:          id,
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Price:       req.GetPrice(),
		Stock:       req.GetStock(),
	}
	s.products[id] = product
	return &pb.ProductResponse{Product: product}, nil
}

func (s *productServer) GetProducts(ctx context.Context, req *pb.GetProductsRequest) (*pb.GetProductsResponse, error) {
	responses := make([]*pb.ProductResponse, 0, len(req.GetIds()))
	for _, id := range req.GetIds() {
		if p, ok := s.products[id]; ok {
			responses = append(responses, &pb.ProductResponse{Product: p})
		}
	}
	return &pb.GetProductsResponse{Products: responses}, nil
}

func (s *productServer) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	responses := make([]*pb.ProductResponse, 0, len(s.products))
	for _, p := range s.products {
		responses = append(responses, &pb.ProductResponse{Product: p})
	}
	return &pb.ListProductsResponse{Products: responses, Total: int32(len(responses))}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	pb.RegisterProductServiceServer(srv, newProductServer())

	log.Printf("Product Service listening at %s", port)
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
