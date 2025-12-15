package main

import (
	"log"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderPb "github.com/mrezayusufy/shop-api/pkg/proto/order"
	productPb "github.com/mrezayusufy/shop-api/pkg/proto/product"
	userPb "github.com/mrezayusufy/shop-api/pkg/proto/user"
	"github.com/mrezayusufy/shop-api/services/gateway/graph"
	"github.com/mrezayusufy/shop-api/services/gateway/graph/generated"
)

// Config holds the gRPC client connections
type Config struct {
	UserClient    userPb.UserServiceClient
	OrderClient   orderPb.OrderServiceClient
	ProductClient productPb.ProductServiceClient
}

func main() {
	// 1. Initialize gRPC Clients
	cfg, err := initClients()
	if err != nil {
		log.Fatalf("could not initialize gRPC clients: %v", err)
	}

	// 2. Initialize GraphQL Server
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: &graph.Resolver{Config: cfg},
	}))

	// 3. Initialize Fiber App
	app := fiber.New()

	// 4. Setup GraphQL Endpoint
	app.Post("/query", adaptor.HTTPHandler(srv))

	// 5. Setup GraphQL Playground (for development)
	app.Get("/", adaptor.HTTPHandler(playground.Handler("GraphQL Playground", "/query")))

	log.Printf("connect to http://localhost:4000/ for GraphQL playground")
	log.Fatal(app.Listen(":4000"))
}

func initClients() (*Config, error) {
	// Set up a connection to the gRPC servers.
	// In a real-world scenario, use service discovery (e.g., Consul, Kubernetes)
	// and secure connections (TLS).

	// User Service
	userConn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	// defer userConn.Close() // In a real app, manage connection lifecycle properly

	// Order Service
	orderConn, err := grpc.Dial("localhost:50053", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	// Product Service
	productConn, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Config{
		UserClient:    userPb.NewUserServiceClient(userConn),
		OrderClient:   orderPb.NewOrderServiceClient(orderConn),
		ProductClient: productPb.NewProductServiceClient(productConn),
	}, nil
}