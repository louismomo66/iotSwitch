# Use an official Go runtime as a base
FROM golang:1.22.2 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the local package files and the .env file to the container's workspace.
COPY . .

# Install godotenv to manage .env files
RUN go get github.com/joho/godotenv

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod tidy

# Build the application to run in a scratch container.
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o myapp ./cmd/myapp

# Use a lightweight Alpine image.
FROM alpine:latest

# Install ca-certificates in case your application makes external HTTPS calls
RUN apk --no-cache add ca-certificates

# Set the working directory in the container
WORKDIR /root/

# Copy the pre-built binary file and the .env file from the previous stage
COPY --from=builder /app/myapp .
COPY --from=builder /app/.env .

# Expose port (if your app uses a port)
EXPOSE 9001

# Command to run the executable
CMD ["./myapp"]
