# Makefile for CacheProxy CLI application

APP_NAME := cacheproxy
MAIN_SRC := .
CMD_SRC := ./cmd
BINARY := $(APP_NAME)

# Commands
GO := go
RM := rm -f

# Default target: build the application
all: build

# Build the CLI application
build:
	@echo "Building $(APP_NAME)..."
	$(GO) build -o $(BINARY) $(MAIN_SRC)

# Run the application
run: build
	@echo "Running $(APP_NAME)..."
	./$(BINARY) start --port 8080 --origin http://dummyjson.com

# Clean up the binary and cache files
clean:
	@echo "Cleaning up..."
	$(RM) $(BINARY)
	$(RM) -r cache/

# Install dependencies (if any)
deps:
	@echo "Installing dependencies..."
	$(GO) get ./...

# Format code using gofmt
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

# Display help
help:
	@echo "Makefile commands:"
	@echo "  build   - Build the CLI application"
	@echo "  run     - Run the CLI application (default port 8080 and origin)"
	@echo "  clean   - Remove the binary and cache files"
	@echo "  deps    - Install dependencies"
	@echo "  fmt     - Format the Go code"

.PHONY: all build run clean deps fmt help
