# Variables
APP_NAME := receipt-processor
GOOS := linux
GOARCH := amd64
OUTPUT_DIR := build

# Default target
all: build

# Clean the build directory
clean:
	rm -rf $(OUTPUT_DIR)

# Build the application for Linux
build:
	mkdir -p $(OUTPUT_DIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(OUTPUT_DIR)/$(APP_NAME)

# Run the application (for local testing)
run:
	go run main.go

# Format code
fmt:
	go fmt ./...

# Run tests
test:
	go test ./...

# Help command to list available targets
help:
	@echo "Available commands:"
	@echo "  make build       Build the Go executable for Linux"
	@echo "  make clean       Remove build artifacts"
	@echo "  make run         Run the application locally"
	@echo "  make fmt         Format Go source code"
	@echo "  make test        Run tests"
