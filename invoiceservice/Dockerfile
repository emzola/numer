# Use the latest Go version available
FROM golang:latest AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o invoiceservice ./cmd

# Create a minimal image to run the service
FROM alpine:3.18
WORKDIR /app

# Install bash and curl to download golang-migrate
RUN apk add --no-cache bash curl

# Download and install golang-migrate
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.1/migrate.linux-amd64.tar.gz | tar xvz \
    && mv migrate /usr/local/bin/golang-migrate

COPY --from=builder /app/invoiceservice .
COPY ./migrations /migrations
EXPOSE 8080
CMD ["./invoiceservice"]
