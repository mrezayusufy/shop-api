package main

import (
	"context"
	"log"
	"net"

	pb "github.com/mrezayusufy/shop-api/pkg/proto/order"

	"google.golang.org/grpc"
)

const (
	port = ":50053"
)

// server is used to implement ecommerce.OrderServiceServer.
type server struct {
	pb.UnimplementedOrderServiceServer
}

// GetOrder implements ecommerce.OrderServiceServer
func (s *server) GetOrder(ctx context.Context, in *pb.GetOrderRequest) (*pb.OrderResponse, error) {
	log.Printf("Received: GetOrder request for ID: %s", in.GetId())
	// Mock data for demonstration
	order := &pb.OrderResponse{
		Id:     in.GetId(),
		UserId: "user-123",
		Items: []*pb.OrderItem{
			{ProductId: "prod-1", Quantity: 2},
			{ProductId: "prod-2", Quantity: 1},
		},
		Total: 59.97,
		Status: "paid",
	}
	return &pb.OrderResponse{
		Id: order.Id,
		UserId: order.UserId,
		Items: order.Items,
		Total: order.Total,
		Status: order.Status,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterOrderServiceServer(s, &server{})
	log.Printf("Order Service listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}