# ==================================================================================== #
# PROTOBUF
# ==================================================================================== #

# Define the paths
PROTO_DIR := proto
INVOICE_PROTO := $(PROTO_DIR)/invoice.proto
USER_PROTO := $(PROTO_DIR)/user.proto

INVOICE_OUT_DIR := invoiceservice/genproto
USER_OUT_DIR := userservice/genproto

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
	
.PHONY: build
build:
	@echo "Building services..."
	docker-compose build

.PHONY: up
up:
	@echo "Starting services..."
	docker-compose up -d

.PHONY: down
down:
	@echo "Stopping services..."
	docker-compose down

.PHONY: migrate
migrate:
	@echo "Applying migrations for invoice service..."
	docker-compose run --rm invoiceservice golang-migrate -path /migrations -database "postgres://${INVOICE_DB_USER}:${INVOICE_DB_PASSWORD}@${INVOICE_DB_HOST}:${INVOICE_DB_PORT}/${INVOICE_DB_NAME}?sslmode=disable" up
	@echo "Applying migrations for user service..."
	docker-compose run --rm userservice golang-migrate -path /migrations -database "postgres://${USER_DB_USER}:${USER_DB_PASSWORD}@${USER_DB_HOST}:${USER_DB_PORT}/${USER_DB_NAME}?sslmode=disable" up

.PHONY: test
test:
	@echo "Running tests for invoice service..."
	cd invoiceservice && go test ./...

.PHONY: clean
clean:
	@echo "Removing containers..."
	docker-compose down -v
