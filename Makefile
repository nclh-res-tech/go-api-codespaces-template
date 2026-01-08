.PHONY: all build run test lint clean setup docker-build docker-run help

# Variables
BINARY_NAME=api-server
MAIN_PATH=./cmd/api-server
DOCKER_IMAGE={{SERVICE_NAME}}

# Default target
all: lint test build

# Build the application
build:
	@echo "Building..."
	go build -o bin/$(BINARY_NAME) $(MAIN_PATH)

# Run the application
run:
	@echo "Running..."
	API_ENVIRONMENT=development go run $(MAIN_PATH)

# Run tests
test:
	@echo "Running tests..."
	go test -v -race -cover ./...

# Run linter
lint:
	@echo "Linting..."
	golangci-lint run ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/
	go clean

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

# Generate code (if applicable)
generate:
	@echo "Generating code..."
	go generate ./...

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE):latest .

# Run Docker container
docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 -e API_ENVIRONMENT=development $(DOCKER_IMAGE):latest

# Setup project (run after cloning template)
setup:
	@echo "Setting up project..."
	./setup.sh

# Show help
help:
	@echo "Available targets:"
	@echo "  all          - Run lint, test, and build"
	@echo "  build        - Build the application"
	@echo "  run          - Run the application locally"
	@echo "  test         - Run tests"
	@echo "  lint         - Run linter"
	@echo "  clean        - Clean build artifacts"
	@echo "  deps         - Download and tidy dependencies"
	@echo "  generate     - Run go generate"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run Docker container"
	@echo "  setup        - Run project setup script"
	@echo "  help         - Show this help message"