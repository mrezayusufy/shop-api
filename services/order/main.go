package main

import (
	"context"
	"log"
	"net"

	pb "github.com/mrezayusufy/shop-api/pkg/proto/order"
	"google.golang.org/grpc"
)

const port = ":50053"

type orderServer struct {
	pb.UnimplementedOrderServiceServer
	orders map[string]*pb.Order
}

func newOrderServer() *orderServer {
	return &orderServer{
		orders: map[string]*pb.Order{
			"order-1": {
				Id:     "order-1",
				UserId: "user-1",
				Items: []*pb.OrderItem{
					{ProductId: "product-1", Quantity: 2},
					{ProductId: "product-2", Quantity: 1},
				},
				TotalAmount: 59.97,
			},
			"order-2": {
				Id:          "order-2",
				UserId:      "user-2",
				Items:       []*pb.OrderItem{{ProductId: "product-3", Quantity: 3}},
				TotalAmount: 29.97,
			},
		},
	}
}

func (s *orderServer) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	order, ok := s.orders[req.GetId()]
	if !ok {
		return &pb.GetOrderResponse{}, nil
	}
	return &pb.GetOrderResponse{Order: order}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	pb.RegisterOrderServiceServer(srv, newOrderServer())

	log.Printf("Order Service listening at %s", port)
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
