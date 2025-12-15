package main

import (
	"log"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"

	"github.com/mrezayusufy/shop-api/services/gateway"
	"github.com/mrezayusufy/shop-api/services/gateway/graph"
	"github.com/mrezayusufy/shop-api/services/gateway/graph/generated"
)

func main() {
	resolver, err := gateway.NewResolver()
	if err != nil {
		log.Fatalf("could not initialize gRPC clients: %v", err)
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: &graph.Resolver{Resolver: resolver},
	}))

	app := fiber.New()
	app.Post("/query", adaptor.HTTPHandler(srv))
	app.Get("/", adaptor.HTTPHandler(playground.Handler("GraphQL Playground", "/query")))

	log.Printf("connect to http://localhost:4000/ for GraphQL playground")
	if err := app.Listen(":4000"); err != nil {
		log.Fatalf("failed to start Fiber: %v", err)
	}
}
