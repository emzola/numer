# ==================================================================================== #
# PROTO
# ==================================================================================== #

# Define the paths
PROTO_DIR := proto
USER_PROTO := user-service/$(PROTO_DIR)/user.proto
INVOICE_PROTO := invoice-service/$(PROTO_DIR)/invoice.proto
STATS_PROTO := stats-service/$(PROTO_DIR)/stats.proto
NOTIFICATION_PROTO := notification-service/$(PROTO_DIR)/notification.proto
REMINDER_PROTO := reminder-service/$(PROTO_DIR)/reminder.proto
ACTIVITY_PROTO := activity-service/$(PROTO_DIR)/activity.proto

USER_OUT_DIR := user-service/$(PROTO_DIR)
INVOICE_OUT_DIR := invoice-service/$(PROTO_DIR)
STATS_OUT_DIR := stats-service/$(PROTO_DIR)
NOTIFICATION_OUT_DIR := notification-service/$(PROTO_DIR)
REMINDER_OUT_DIR := reminder-service/$(PROTO_DIR)
ACTIVITY_OUT_DIR := activity-service/$(PROTO_DIR)

# Check if output directories exist, if not create them
.PHONY: create_dirs
create_dirs:
	@mkdir -p $(USER_OUT_DIR)
	@mkdir -p $(INVOICE_OUT_DIR)
	@mkdir -p $(STATS_OUT_DIR)
	@mkdir -p $(NOTIFICATION_OUT_DIR)
	@mkdir -p $(REMINDER_OUT_DIR)
	@mkdir -p $(ACTIVITY_OUT_DIR)

# Define the protoc command
PROTOC := protoc
PROTOC_GEN_GO := protoc-gen-go
PROTOC_GEN_GRPC_GO := protoc-gen-go-grpc

# Generate the protobuf files
.PHONY: proto
proto: create_dirs user_proto invoice_proto stats_proto notification_proto reminder_proto activity_proto 

user_proto: $(USER_PROTO)
	$(PROTOC) --go_out=. --go-grpc_out=. $(USER_PROTO)

invoice_proto: $(INVOICE_PROTO)
	$(PROTOC) --go_out=. --go-grpc_out=. $(INVOICE_PROTO)

stats_proto: $(STATS_PROTO)
	$(PROTOC) --go_out=. --go-grpc_out=. $(STATS_PROTO)

notification_proto: $(NOTIFICATION_PROTO)
	$(PROTOC) --go_out=. --go-grpc_out=. $(NOTIFICATION_PROTO)

reminder_proto: $(REMINDER_PROTO)
	$(PROTOC) --go_out=. --go-grpc_out=. $(REMINDER_PROTO)

activity_proto: $(ACTIVITY_PROTO)
	$(PROTOC) --go_out=. --go-grpc_out=. $(ACTIVITY_PROTO)
	
# Clean the generated files
.PHONY: cleanproto
cleanproto:
	rm -f $(USER_OUT_DIR)/*.pb.go
	rm -f $(INVOICE_OUT_DIR)/*.pb.go
	rm -f $(STATS_OUT_DIR)/*.pb.go
	rm -f $(NOTIFICATION_OUT_DIR)/*.pb.go
	rm -f $(REMINDER_OUT_DIR)/*.pb.go
	rm -f $(ACTIVITY_OUT_DIR)/*.pb.go

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #
	
# Variables
DOCKER_COMPOSE=docker-compose
GOOSE_CMD=docker-compose run --rm
GOOSE_BIN=goose
MIGRATION_DIR=migrations

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
migrate-invoice:
	$(GOOSE_CMD) invoice-service $(GOOSE_BIN) -dir $(MIGRATION_DIR) postgres "$(INVOICE_DB_URL)" up

# Apply Goose migrations for the activity service
migrate-activity:
	$(GOOSE_CMD) activity-service $(GOOSE_BIN) -dir $(MIGRATION_DIR) postgres "$(ACTIVITY_DB_URL)" up

# Apply migrations for all services
migrate-all: migrate-user migrate-invoice migrate-activity

# Show logs for all services
logs:
	$(DOCKER_COMPOSE) logs -f

# Run tests for all services
test:
	@echo "Running tests for user service..."
	cd user-service && go test ./...

	@echo "Running tests for invoice service..."
	cd invoice-service && go test ./...

	@echo "Running tests for stats service..."
	cd stats-service && go test ./...

	@echo "Running tests for notification service..."
	cd notification-service && go test ./...

	@echo "Running tests for reminder service..."
	cd reminder-service && go test ./...

	@echo "Running tests for activity service..."
	cd activity-service && go test ./...

