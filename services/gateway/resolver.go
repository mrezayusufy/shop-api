package gateway

import (
	orderpb "github.com/mrezayusufy/shop-api/pkg/proto/order"
	productpb "github.com/mrezayusufy/shop-api/pkg/proto/product"
	userpb "github.com/mrezayusufy/shop-api/pkg/proto/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Resolver struct {
	UserClient    userpb.UserServiceClient
	ProductClient productpb.ProductServiceClient
	OrderClient   orderpb.OrderServiceClient
}

func NewResolver() (*Resolver, error) {
	userConn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	productConn, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	orderConn, err := grpc.Dial("localhost:50053", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Resolver{
		UserClient:    userpb.NewUserServiceClient(userConn),
		ProductClient: productpb.NewProductServiceClient(productConn),
		OrderClient:   orderpb.NewOrderServiceClient(orderConn),
	}, nil
}
