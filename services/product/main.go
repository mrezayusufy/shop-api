package product

import (
    "context"
    "fmt"
    "sync"

    pb "github.com/mrezayusufy/shop-api/pkg/proto/product"
)

type Service struct {
    pb.UnimplementedProductServiceServer
    products map[string]*pb.ProductResponse
    mu       sync.RWMutex
}

func NewService() *Service {
    svc := &Service{
        products: make(map[string]*pb.ProductResponse),
    }
    // Seed data
    svc.products["prod_1"] = &pb.ProductResponse{
        Id:          "prod_1",
        Name:        "Laptop",
        Description: "High-performance laptop",
        Price:       999.99,
        Stock:       10,
    }
    svc.products["prod_2"] = &pb.ProductResponse{
        Id:          "prod_2",
        Name:        "Mouse",
        Description: "Wireless mouse",
        Price:       29.99,
        Stock:       50,
    }
    return svc
}

func (s *Service) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.ProductResponse, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    if product, exists := s.products[req.Id]; exists {
        return product, nil
    }
    return nil, fmt.Errorf("product not found")
}

func (s *Service) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.ProductResponse, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    id := fmt.Sprintf("prod_%d", len(s.products)+1)
    product := &pb.ProductResponse{
        Id:          id,
        Name:        req.Name,
        Description: req.Description,
        Price:       req.Price,
        Stock:       req.Stock,
    }
    s.products[id] = product
    return product, nil
}

func (s *Service) GetProducts(ctx context.Context, req *pb.GetProductsRequest) (*pb.GetProductsResponse, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    var products []*pb.ProductResponse
    for _, id := range req.Ids {
        if product, exists := s.products[id]; exists {
            products = append(products, product)
        }
    }
    return &pb.GetProductsResponse{Products: products}, nil
}

func (s *Service) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    var products []*pb.ProductResponse
    for _, product := range s.products {
        products = append(products, product)
    }
    return &pb.ListProductsResponse{
        Products: products,
        Total:    int32(len(products)),
    }, nil
}