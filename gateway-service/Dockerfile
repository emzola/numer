# Use the latest version of Go
FROM golang:latest AS builder

# Set environment variables
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

# Set the working directory inside the container
WORKDIR /app

# Copy Go module files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -o gateway-service ./cmd/main.go

# Second stage: build a minimal image
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/gateway-service .

# Copy the .env file (if necessary)
COPY .env .env

# Expose the gRPC port
EXPOSE 50057

# Command to start the application
CMD ["./gateway-service"]
