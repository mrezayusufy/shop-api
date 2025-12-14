# Go Microservice Monorepo Structure: gRPC, Fiber, and GraphQL (gqlgen)

This document outlines a scalable and maintainable file and folder structure for a Go microservice monorepo, integrating gRPC for inter-service communication, Fiber for external HTTP/REST APIs, and GraphQL via `gqlgen` for flexible data querying. The structure adheres to the standard Go project layout while incorporating best practices for microservices and monorepos.

## 1. Monorepo Root Structure

The top-level structure is designed to separate shared components from individual services and development tools.

| Directory | Purpose | Key Contents |
| :--- | :--- | :--- |
| `api/` | **Shared API Definitions** | Central location for all Protobuf (`.proto`) and shared GraphQL schema files (`.graphqls`). |
| `pkg/` | **Shared Go Packages** | Reusable, public Go code that can be imported by any service. |
| `services/` | **Individual Microservices** | Contains the source code for each distinct microservice. |
| `tools/` | **Development Tools** | Scripts and custom Go tools for code generation, linting, and building. |
| `go.mod` | **Dependency Management** | The main Go module file for the monorepo. |
| `Makefile` | **Automation** | Common build, test, and code generation commands. |

```
/
├── api/
│   ├── proto/            # Protobuf definitions for gRPC
│   │   ├── user/user.proto
│   │   └── product/product.proto
│   └── graphql/          # Shared GraphQL schemas and common types
│       └── common.graphqls
├── pkg/                  # Shared Go packages
│   ├── auth/             # Authentication/Authorization utilities
│   ├── config/           # Centralized configuration loading
│   └── utils/            # General utilities (logging, error handling)
├── services/             # All microservices reside here
│   ├── user-service/
│   ├── product-service/
│   └── ...
├── tools/                # Development scripts and helper binaries
├── go.mod
├── go.sum
└── Makefile
```

## 2. Individual Microservice Structure

Each service within the `services/` directory follows a consistent internal structure, leveraging the **Clean Architecture** or **Hexagonal Architecture** principles to separate concerns. This ensures that the core business logic (`core/`) is independent of the delivery mechanisms (gRPC, Fiber, GraphQL) and external dependencies (database, external services).

**Example: `services/user-service`**

```
services/user-service/
├── cmd/
│   └── main.go           # Service entry point (initializes server, config, etc.)
├── configs/              # Service-specific configuration files (e.g., YAML, JSON)
├── internal/             # Private application code (cannot be imported by other services)
│   ├── core/             # Core Business Logic
│   │   ├── domain/       # Entities, Value Objects (e.g., User struct)
│   │   ├── service/      # Business logic implementation (e.g., UserService interface)
│   │   └── ports/        # Interfaces for external dependencies (e.g., UserRepository interface)
│   ├── adapters/         # Implementation of external dependencies (ports)
│   │   ├── repository/   # Database implementation (e.g., PostgresUserRepository)
│   │   └── grpc_client/  # Clients for calling other gRPC services
│   └── delivery/         # API/Transport Layer (The "delivery" of the core logic)
│       ├── grpc/         # gRPC server implementation (handlers for user.proto)
│       ├── http/         # Fiber HTTP server implementation (controllers/handlers)
│       └── graphql/      # gqlgen resolvers and schema
│           ├── resolvers/  # Implementation of GraphQL resolvers
│           └── graph/      # Generated gqlgen code and schema
├── migrations/           # Database migration scripts
├── Dockerfile
└── README.md
```

### Key Component Breakdown

#### A. API Definitions (`api/`)

This is the source of truth for all service contracts.

*   **Protobuf (`api/proto/`)**:
    *   Contains all `.proto` files.
    *   The `Makefile` at the monorepo root should contain a target to run `protoc` and generate Go code into a designated location, typically within `pkg/` or a service's `internal/` directory.
    *   *Example:* `api/proto/user/user.proto` defines the gRPC service for the `user-service`.

*   **GraphQL (`api/graphql/`)**:
    *   Contains shared GraphQL type definitions.
    *   Service-specific schemas are often kept within the service's `internal/delivery/graphql/` directory, but common types (e.g., `Pagination`, `Error`) should be shared here.

#### B. Delivery Layer (`internal/delivery/`)

This layer handles incoming requests and translates them into calls to the core business logic (`internal/core/service`).

1.  **gRPC (`internal/delivery/grpc/`)**:
    *   Implements the methods defined in the generated gRPC interface.
    *   It receives the gRPC request, validates it, calls the core service, and translates the core service's response/error back into a gRPC response.

2.  **Fiber (`internal/delivery/http/`)**:
    *   Contains the HTTP route definitions and handlers using the Fiber framework.
    *   Handlers receive the HTTP request, extract parameters, call the core service, and send the HTTP response (e.g., JSON).

3.  **GraphQL (`internal/delivery/graphql/`)**:
    *   **`resolvers/`**: Contains the Go struct methods that implement the GraphQL query and mutation logic. These methods call the core service.
    *   **`graph/`**: Contains the generated code from `gqlgen` (e.g., `generated.go`, `model/`). The service's specific schema file (`schema.graphqls`) is also placed here, and `gqlgen` is configured to use it.

#### C. Core Logic (`internal/core/`)

This is the heart of the application, containing the business rules.

*   **`domain/`**: Pure Go structs and types representing the business entities.
*   **`ports/`**: Go interfaces that define what the service needs from the outside world (e.g., `UserRepository`, `PaymentGateway`).
*   **`service/`**: The actual business logic implementation. It takes the `ports` interfaces as dependencies (Dependency Inversion Principle).

## 3. Build and Code Generation Automation

A robust `Makefile` at the monorepo root is essential for consistency and automation.

| Makefile Target | Description |
| :--- | :--- |
| `proto-gen` | Generates Go code from all `.proto` files in `api/proto/`. |
| `gql-gen` | Runs `gqlgen generate` for all services that use GraphQL. |
| `build-all` | Builds all microservices into executable binaries. |
| `test-all` | Runs unit and integration tests across all services. |
| `run-user` | Runs the `user-service` locally. |

### Example `Makefile` Snippet

```makefile
# Root Makefile
.PHONY: proto-gen gql-gen build-all

# Generate gRPC code from all .proto files
proto-gen:
	@echo "Generating Protobuf code..."
	@find api/proto -name "*.proto" -exec protoc \
		--proto_path=api/proto \
		--go_out=paths=source_relative:pkg/pb \
		--go-grpc_out=paths=source_relative:pkg/pb \
		{} \;

# Generate gqlgen code for all services
gql-gen:
	@echo "Generating gqlgen code..."
	@for service in $(shell find services -maxdepth 1 -mindepth 1 -type d); do \
		if [ -f "$$service/internal/delivery/graphql/gqlgen.yml" ]; then \
			echo "-> Generating for $$service"; \
			(cd "$$service/internal/delivery/graphql" && gqlgen generate); \
		fi \
	done

# Build all services
build-all:
	@echo "Building all services..."
	@for service in $(shell find services -maxdepth 1 -mindepth 1 -type d); do \
		go build -o bin/$(notdir $$service) $$service/cmd/main.go; \
	done
```

## References

The proposed structure is heavily influenced by the official Go Project Layout [1] and best practices for building scalable microservices with Clean Architecture [2].

[1] Go Project Layout - *Standard Go Project Layout*
[2] The Clean Architecture - *Robert C. Martin (Uncle Bob)*