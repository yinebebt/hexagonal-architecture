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