package order

import (
	"context"
	"fmt"
	"sync"
	"time"

	pb "github.com/mrezayusufy/shop-api/pkg/proto/order"
)

type Service struct {
    pb.UnimplementedOrderServiceServer
    orders map[string]*pb.OrderResponse
    mu     sync.RWMutex
}

func NewService() *Service {
    return &Service{
        orders: make(map[string]*pb.OrderResponse),
    }
}

func (s *Service) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.OrderResponse, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    if order, exists := s.orders[req.Id]; exists {
        return order, nil
    }
    return nil, fmt.Errorf("order not found")
}

func (s *Service) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.OrderResponse, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    var total float64
    for _, item := range req.Items {
        total += item.Price * float64(item.Quantity)
    }

    id := fmt.Sprintf("order_%d", len(s.orders)+1)
    order := &pb.OrderResponse{
        Id:        id,
        UserId:    req.UserId,
        Items:     req.Items,
        Total:     total,
        Status:    "pending",
        CreatedAt: time.Now().Format(time.RFC3339),
    }
    s.orders[id] = order
    return order, nil
}

func (s *Service) GetUserOrders(ctx context.Context, req *pb.GetUserOrdersRequest) (*pb.GetUserOrdersResponse, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    var orders []*pb.OrderResponse
    for _, order := range s.orders {
        if order.UserId == req.UserId {
            orders = append(orders, order)
        }
    }
    return &pb.GetUserOrdersResponse{Orders: orders}, nil
}