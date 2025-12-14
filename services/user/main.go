package user

import (
    "context"
    "fmt"
    "sync"

    pb "github.com/mrezayusufy/shop-api/pkg/proto/user"
)

type Service struct {
    pb.UnimplementedUserServiceServer
    users map[string]*pb.UserResponse
    mu    sync.RWMutex
}

func NewService() *Service {
    svc := &Service{
        users: make(map[string]*pb.UserResponse),
    }
    // Seed data
    svc.users["user_1"] = &pb.UserResponse{
        Id:    "user_1",
        Name:  "John Doe",
        Email: "john@example.com",
    }
    return svc
}

func (s *Service) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    if user, exists := s.users[req.Id]; exists {
        return user, nil
    }
    return nil, fmt.Errorf("user not found")
}

func (s *Service) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    id := fmt.Sprintf("user_%d", len(s.users)+1)
    user := &pb.UserResponse{
        Id:    id,
        Name:  req.Name,
        Email: req.Email,
    }
    s.users[id] = user
    return user, nil
}

func (s *Service) GetUsers(ctx context.Context, req *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    var users []*pb.UserResponse
    for _, id := range req.Ids {
        if user, exists := s.users[id]; exists {
            users = append(users, user)
        }
    }
    return &pb.GetUsersResponse{Users: users}, nil
}