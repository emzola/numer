# ==================================================================================== #
# PROTOBUF
# ==================================================================================== #

# Define the paths
PROTO_DIR := proto
INVOICE_PROTO := $(PROTO_DIR)/invoiceservice/invoice.proto

INVOICE_OUT_DIR := invoiceservice/genproto

# Check if output directories exist, if not create them
.PHONY: create_dirs
create_dirs:
	@mkdir -p $(INVOICE_OUT_DIR)

# Define the protoc command
PROTOC := protoc
PROTOC_GEN_GO := protoc-gen-go
PROTOC_GEN_GRPC_GO := protoc-gen-go-grpc

# Generate the protobuf files
.PHONY: proto
proto: create_dirs invoice_proto

invoice_proto: $(INVOICE_PROTO)
	$(PROTOC) --go_out=. --go-grpc_out=. $(INVOICE_PROTO)
	
# Clean the generated files
.PHONY: clean
clean:
	rm -f $(INVOICE_OUT_DIR)/*.pb.go

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #
	
# Run the invoice service
.PHONY: invoice
rider:
	@go run invoiceservice/cmd/*.go