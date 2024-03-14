# Hexagonal-architecture
![build-workflow](https://github.com/Yinebeb-01/hexagonal-architecture/actions/workflows/build-and-test.yml/badge.svg)

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
- [ ] gRPC
- [ ] GraphQL
-  [ ] WebSocket

### Repository

- [x] Gorm/sqlite
- [ ] Mongodb
- [ ] Postgres

## Project structure

```
/app
|-- /cmd
|   |-- main.go
|-- /docs   
|-- /internal
|   |-- /adapter
|   |   |-- /handler
|   |   |   |-- /rest
|   |   |   |-- /gRPC
|   |   |-- /reository
|   |   |   |-- /gorm
|   |   |   |-- /mongo
|   |   |-- /glue
|   |   |   |-- /route
|   |   |   |-- route.go
|   |   |-- /dto
|   |   |-- /templates
|   |-- /core
|   |   |-- /entity
|   |   |-- /port
|   |   |-- /service
|   |   |   |-- /test
|   |   |-- /util   
```

## Installation

Install **godog** binary:

```bash
go install github.com/cucumber/godog/cmd/godog@latest
```

Use `go test` command to run feature tests since godog's cli is deprecated.