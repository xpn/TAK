.PHONY: all build clean test install deps help

# Variables
BINARY_NAME=tunnel-manager
GO=go
GOFLAGS=
LDFLAGS=

# Default target
all: build

grpc:
	@echo "Generating gRPC code..."
	protoc --go_out=./pkg/grpc/ --go_opt=paths=source_relative --go-grpc_out=./pkg/grpc/ --go-grpc_opt=paths=source_relative -I ./proto -I ~/Projects/gogoproto  ./proto/authservice.proto

# Build the project
build:
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BINARY_NAME) .

# Clean build artifacts
clean:
	$(GO) clean
	rm -f $(BINARY_NAME)

# Run tests
test:
	$(GO) test -v ./...

# Install dependencies
deps:
	$(GO) mod download
	$(GO) mod tidy
