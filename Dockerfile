FROM golang:latest

WORKDIR /app

COPY go.mod .

COPY go.sum .

RUN go mod download

COPY . .

# Set rest port
ENV PORT 5000

# Build the app
RUN go build -o app cmd/main.go

RUN find . -name "*.go" -type f -delete

EXPOSE $PORT

CMD ["./app"]
