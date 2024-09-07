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
RUN go build -o user-service ./cmd/main.go

# Install Goose in the builder stage
RUN go install github.com/pressly/goose/v3/cmd/goose@v3.21.1

# Second stage: build a minimal image
FROM alpine:latest

# Install required dependencies for Goose migration and PostgreSQL client
RUN apk --no-cache add ca-certificates bash postgresql-client

# Set the working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/user-service .

# Copy Goose binary from the builder stage
COPY --from=builder /go/bin/goose /usr/local/bin/goose

# Copy Goose migration files
COPY migrations /app/db/migrations

# Copy the .env file (if necessary)
COPY .env .env

# Expose the gRPC port
EXPOSE 50051

# Command to run the migrations and start the application
CMD ["/bin/sh", "-c", "/usr/local/bin/goose -dir ./db/migrations postgres \"$DATABASE_URL\" up && ./user-service"]
