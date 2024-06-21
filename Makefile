.PHONY: all build run test clean tidy run1

# Load environment variables from .env file
include .env
export

# Go parameters
GO_CMD=go
GO_BUILD=$(GO_CMD) build
GO_RUN=$(GO_CMD) run
GO_TEST=$(GO_CMD) test

# Binary name
BINARY_NAME=myapp

all: build

# Build the binary
build:
	$(GO_BUILD) -o $(BINARY_NAME) ./cmd/myapp

# Run the binary
run: build
	./$(BINARY_NAME)

# Run tests
test:
	$(GO_TEST) ./...

# Clean up build artifacts
clean:
	rm -f $(BINARY_NAME)

# Tidy up dependencies
tidy:
	$(GO_CMD) mod tidy

# Run the application directly
run1:
	$(GO_RUN) ./cmd/myapp/main.go
