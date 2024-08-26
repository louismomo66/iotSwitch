# .PHONY: all build run test clean tidy run1 pull commit push

# Load environment variables from .env file
# include .env
# export

# # Go parameters
# GO_CMD=go
# GO_BUILD=$(GO_CMD) build
# GO_RUN=$(GO_CMD) run
# GO_TEST=$(GO_CMD) test

# # Binary name
# BINARY_NAME=myapp

# all: build

# # Build the binary
# build:
# 	$(GO_BUILD) -buildvcs=false -o $(BINARY_NAME) ./cmd/myapp

# # Run the binary
# run: build
# 	./$(BINARY_NAME)

# # Run tests
# test:
# 	$(GO_TEST) ./...

# # Clean up build artifacts
# clean:
# 	rm -f $(BINARY_NAME)

# # Tidy up dependencies
# tidy:
# 	$(GO_CMD) mod tidy

# # Run the application directly
# run1:
# 	$(GO_RUN) ./cmd/myapp/main.go

# # Pull latest changes from the repository
# pull:
# 	git pull origin main

# # Commit changes with a message
# commit:
# 	git add .
# 	git commit -m "$(m)"

# # Push changes to the repository
# push:
# 	git push origin main
BINARY_NAME=myapp
DSN=host=localhost port=5434 user=postgres password=XNEHk9iSGp9GItlxVuXYfmbEiTyugBuZ dbname=iotSwitch sslmode=disable timezone=UTC connect_timeout=5

## build: Build binary
build:
	@echo "Building..."
	env CGO_ENABLED=0 go build -ldflags="-s -w" -o ${BINARY_NAME} ./cmd/myapp
	@echo "Built!"

## run: builds and runs the application
run: build
	@echo "Starting..."
	@env DSN=${DSN} ./${BINARY_NAME} &
	@echo "Started!"

## clean: runs go clean and deletes binaries
clean:
	@echo "Cleaning..."
	@go clean
	@rm ${BINARY_NAME}
	@echo "Cleaned!"

## start: an alias to run
start: run

## stop: stops the running application
stop:
	@echo "Stopping..."
	@-pkill -SIGTERM -f "./${BINARY_NAME}"
	@echo "Stopped!"

## restart: stops and starts the application
restart: stop start

## test: runs all tests
test:
	@go test -v ./...

