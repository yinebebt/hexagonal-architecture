# Hexagonal-architecture
![build-workflow](https://github.com/Yinebeb-01/hexagonal-architecture/actions/workflows/build-and-test.yml/badge.svg)

RESTful API built with Gin and GORM, demonstrating the principles of Hexagonal Architecture. 
It manages video entities, allowing admins to create, and manage videos. 

Explore the concept of Hexagonal Architecture further in this blog post: 
[Hexagonal-architecture](https://medium.com/@yinebeb-tariku/hexagonal-architecture-93a946776242).

**Installation**

Install **godog** binary:
```bash
go install github.com/cucumber/godog/cmd/godog@latest
```

Use `go test` command to run feature tests since godog's cli is deprecated.


# Project Folder Structure

This project follows a specific folder structure to organize its components effectively. Below is an explanation of each directory:

- `bin/`: Compiled binaries are stored in this directory after building the project.

- `cmd/`: Contains command-line application entry points.

    - `http/`: This subdirectory contains the HTTP server entry point.

- `docs/`: Documentation related to the project is stored here.

- `internal/`: Internal packages of the project are located here.

    - `adapter/`: Contains adapters for external systems.

        - `cache/`: Cache adapters are stored in this directory.

            - `redis/`: Redis cache adapter implementation resides here.

        - `handler/`: Handler adapters for different protocols (e.g., HTTP, gRPC) are stored here.

            - `http/`: HTTP handler adapter implementation is located in this directory.

        - `repository/`: Repository adapters for various databases are stored here.

            - `postgres/`: Postgres repository adapter implementation resides here.

                - `migrations/`: Database migrations related to the Postgres repository adapter are stored here.

        - `token/`: Token adapters for different token types (e.g., JWT, Paseto) are stored here.

            - `paseto/`: Paseto token adapter implementation is located in this directory.

    - `core/`: Contains the core business logic of the application.

        - `domain/`: Domain entities representing specific objects within the application's domain are stored here.

        - `port/`: Ports (interfaces) defining interactions with adapters are stored in this directory.

        - `service/`: Core application services are located here.

        - `util/`: Utility functions and helpers used across the project are stored in this directory.

This structured approach organizes the project components in a clear and systematic manner, making it easier to understand and maintain the codebase.
