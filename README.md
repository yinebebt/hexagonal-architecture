# Hexagonal-architecture

![build-workflow](https://github.com/yinebebt/hexagonal-architecture/actions/workflows/build-and-test.yml/badge.svg)

Hexagonal architecture is a design pattern suitable for building scalable and complex projects.
This repository serves as a demonstration of the principles of Hexagonal Architecture in a Go project.

The goal of this project is to provide a straightforward example that developers can use to understand and apply
Hexagonal Architecture in their own projects. By following the structure and patterns demonstrated here, developers
can build scalable and maintainable systems with ease.

In this demo, the core business functionality revolves around managing user and video entities.
Administrators have the capability to manage videos, while users are provided with access to view the available videos.

Explore the concept of Hexagonal Architecture further
in [Hexagonal-architecture](https://medium.com/@yinebeb-tariku/hexagonal-architecture-93a946776242).

## Adapters

### Handler

- [x] REST API - GIN
- [x] GraphQL (Playground at `/v1/graphql`)
- [ ] gRPC
- [ ] WebSocket

### Repository

- [x] Sqlite
- [x] Postgres

## Project Structure

```txt
hexagonal-architecture/
├── cmd/
│   └── main.go                 # Application entry point
├── docs/                       # Swagger documentation
├── internal/
│   ├── adapter/                # External adapters (driving & driven)
│   │   ├── repository/         # Repository adapters (driven/outbound)
│   │   │   ├── sqlite/        # SQLite implementation
│   │   │   └── postgres/      # PostgreSQL implementation
│   │   ├── rest/               # REST adapter (driving/inbound)
│   │   │   ├── handler.go     # Handler implementation & routes
│   │   │   └── middleware.go  # HTTP middleware
│   │   ├── graphql/            # GraphQL adapter (driving/inbound)
│   │   │   └── handler.go     # Schema, resolvers & handler
│   │   └── templates/          # HTML templates
│   └── core/                   # Business logic (hexagon)
│       ├── entity/             # Domain entities
│       ├── port/                # Ports (interfaces)
│       │   ├── service.go      # port.VideoService interface
│       │   └── repository.go   # port.VideoRepository interface
│       └── service/            # Business services
├── go.mod
├── go.sum
└── README.md
```

## Architecture Principles

- **Core (Domain)**: Contains business logic, entities, and ports (interfaces)
- **Adapters**: Implement ports for external concerns (HTTP, databases)
- **Dependency Inversion**: Core depends on abstractions (ports), not implementations
- **Clean Separation**: Business logic is independent of frameworks and databases

## Installation

```bash
# Clone the repository
git clone https://github.com/yinebebt/hexagonal-architecture.git
cd hexagonal-architecture

# Install dependencies
go mod download

# Build the application
go build ./cmd/main.go

# Run tests
go test ./...
```

## Running the Application

```bash
# Using SQLite (default)
go run ./cmd/main.go -dbtype=sqlite -dsn=app.db

# Using PostgreSQL
go run ./cmd/main.go -dbtype=postgres -dsn="postgres://user:password@localhost/dbname"

# Set custom port
PORT=8080 go run ./cmd/main.go
```

## Testing

Run all tests:

```bash
go test -v -race ./...
```
