.PHONY: build run proto test clean help

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Binary names
BINARY_NAME=syncer
MAIN_PATH=cmd/syncer/main.go

# Protocol buffer settings
PROTO_PATH=control/proto
PROTO_FILE=$(PROTO_PATH)/control.proto
PROTO_OUT_GO=.
PROTO_OUT_WEB=control/web

# Default target
all: build

# Build the binary
build:
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_PATH)

# Run the application
run: build
	./$(BINARY_NAME)

# Run directly without building binary
run-go:
	$(GOCMD) run $(MAIN_PATH)

# Generate protocol buffer files
proto:
	@echo "Generating Go protobuf files..."
	protoc -I=$(PROTO_PATH) \
		--go_out=$(PROTO_OUT_GO) \
		--go-grpc_out=$(PROTO_OUT_GO) \
		$(PROTO_FILE)
	@echo "Generating JavaScript/gRPC-Web files..."
	protoc -I=$(PROTO_PATH) \
		--js_out=import_style=commonjs,binary:$(PROTO_OUT_WEB) \
		--grpc-web_out=import_style=commonjs,mode=grpcwebtext:$(PROTO_OUT_WEB) \
		$(PROTO_FILE)
	@echo "Protocol buffer generation complete"

# Run tests
test:
	$(GOTEST) -v ./...

# Run tests for specific package
test-extractor:
	$(GOTEST) -v ./extractor

test-injector:
	$(GOTEST) -v ./injector

test-control:
	$(GOTEST) -v ./control

test-helpers:
	$(GOTEST) -v ./helpers/...

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

# Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Install protoc plugins (run once)
install-proto-deps:
	@echo "Installing protobuf compiler plugins..."
	$(GOGET) google.golang.org/protobuf/cmd/protoc-gen-go
	$(GOGET) google.golang.org/grpc/cmd/protoc-gen-go-grpc
	@echo "Make sure protoc is installed: https://grpc.io/docs/protoc-installation/"

# Help
help:
	@echo "Available targets:"
	@echo "  build              - Build the syncer binary"
	@echo "  run                - Build and run the application"
	@echo "  run-go             - Run directly with 'go run'"
	@echo "  proto              - Generate protobuf files"
	@echo "  test               - Run all tests"
	@echo "  test-<package>     - Run tests for specific package"
	@echo "  clean              - Clean build artifacts"
	@echo "  deps               - Download and tidy dependencies"
	@echo "  install-proto-deps - Install protobuf compiler plugins"
	@echo "  help               - Show this help message"