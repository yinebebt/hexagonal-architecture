swag:
	swag init -g cmd/main.go

test:
	go test ./...

run:
	go run ./cmd/main.go

build:
	go build -o app ./cmd/main.go

.PHONY: swag test run build