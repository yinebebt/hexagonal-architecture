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


## Project structure

- `bin/`: Compiled binaries are stored in this directory after building the project.

- `cmd/`: Contains command-line application entry points.

- `internal/`: Internal packages of the project are located here.

    - `adapter/`: Contains adapters for external systems.

        - `handler/`: Handler adapters for different protocols (e.g., HTTP(REST), gRPC) are stored here.

            - `rest/`: HTTP handler adapter implementation is located in this directory.

        - `repository/`: Repository adapters for various databases are stored here.

            - `gorm(sqlite)/`: Gorm(sqlite) repository adapter implementation resides here.

    - `core/`: Contains the core business logic of the application.

        - `entity/`: Domain entities representing specific objects within the application's domain are stored here.

        - `port/`: Ports (interfaces) defining interactions with adapters are stored in this directory.

        - `service/`: Core application services are located here.

        - `util/`: Utility functions and helpers used across the project are stored in this directory.