# ==================================================================================== #
# PROTO
# ==================================================================================== #

# Define the paths
PROTO_DIR := proto
INVOICE_PROTO := $(PROTO_DIR)/invoice.proto
USER_PROTO := user-service/$(PROTO_DIR)/user.proto

INVOICE_OUT_DIR := invoiceservice/genproto
USER_OUT_DIR := user-service/$(PROTO_DIR)

# Check if output directories exist, if not create them
.PHONY: create_dirs
create_dirs:
	@mkdir -p $(INVOICE_OUT_DIR)
	@mkdir -p $(USER_OUT_DIR)

# Define the protoc command
PROTOC := protoc
PROTOC_GEN_GO := protoc-gen-go
PROTOC_GEN_GRPC_GO := protoc-gen-go-grpc

# Generate the protobuf files
.PHONY: proto
proto: create_dirs invoice_proto user_proto

invoice_proto: $(INVOICE_PROTO)
	$(PROTOC) --go_out=. --go-grpc_out=. $(INVOICE_PROTO)

user_proto: $(USER_PROTO)
	$(PROTOC) --go_out=. --go-grpc_out=. $(USER_PROTO)
	
# Clean the generated files
.PHONY: cleanproto
cleanproto:
	rm -f $(INVOICE_OUT_DIR)/*.pb.go
	rm -f $(USER_OUT_DIR)/*.pb.go

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #
	
# Variables
DOCKER_COMPOSE=docker-compose
GOOSE_CMD=docker-compose run --rm
GOOSE_BIN=goose
MIGRATION_DIR=migrations

# Load environment variables from .env files
# include user-service/.env
# include invoice-service/.env

# Targets
.PHONY: all build up down test migrate-user migrate-invoice stop restart

# Build the services
build:
	$(DOCKER_COMPOSE) build

# Start all services
up:
	$(DOCKER_COMPOSE) up -d

# Stop all services
down:
	$(DOCKER_COMPOSE) down

# Stop and remove all containers
stop:
	$(DOCKER_COMPOSE) down --volumes --remove-orphans

# Restart services
restart:
	$(DOCKER_COMPOSE) down
	$(DOCKER_COMPOSE) up -d

# Apply Goose migrations for the user service
migrate-user:
	$(GOOSE_CMD) user-service $(GOOSE_BIN) -dir $(MIGRATION_DIR) postgres "$(USER_DB_URL)" up

# Apply Goose migrations for the invoice service
# migrate-invoice:
# 	$(GOOSE_CMD) invoice-service $(GOOSE_BIN) -dir $(MIGRATION_DIR) postgres "$(DATABASE_URL)" up

# Apply migrations for all services
migrate-all: migrate-user

# Show logs for all services
logs:
	$(DOCKER_COMPOSE) logs -f

# Run tests for all services
test:
	@echo "Running tests for user service..."
	cd user-service && go test ./...

