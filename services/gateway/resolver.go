package gateway

import (
	orderpb "github.com/mrezayusufy/shop-api/pkg/proto/order"
	productpb "github.com/mrezayusufy/shop-api/pkg/proto/product"
	userpb "github.com/mrezayusufy/shop-api/pkg/proto/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Resolver struct {
    userClient    userpb.UserServiceClient
    productClient productpb.ProductServiceClient
    orderClient   orderpb.OrderServiceClient
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
        userClient:    userpb.NewUserServiceClient(userConn),
        productClient: productpb.NewProductServiceClient(productConn),
        orderClient:   orderpb.NewOrderServiceClient(orderConn),
    }, nil
}