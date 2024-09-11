package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/emzola/numer/gateway-service/internal/grpcutil"
	userpb "github.com/emzola/numer/user-service/proto"
)

func (h *Handler) CreateCustomerHandler(w http.ResponseWriter, r *http.Request) {
	// Extract user from context
	user := h.contextGetUser(r)

	// Decode the JSON body into the HTTP request struct
	var httpReq CreateCustomerHTTPReq
	err := h.decodeJSON(w, r, &httpReq)
	if err != nil {
		h.badRequestResponse(w, r, err)
		return
	}

	// Convert the HTTP request into the gRPC CreateUserRequest
	grpcReq := &userpb.CreateCustomerRequest{
		UserId:  user.Id,
		Name:    httpReq.Name,
		Email:   httpReq.Email,
		Address: httpReq.Address,
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// Create gRPC connection to user service
	conn, err := grpcutil.ServiceConnection(ctx, "user-service", h.registry)
	if err != nil {
		h.serverErrorResponse(w, r, err)
		return
	}
	defer conn.Close()

	client := userpb.NewUserServiceClient(conn)
	grpcRes, err := client.CreateCustomer(ctx, grpcReq)
	if err != nil {
		h.serverErrorResponse(w, r, err)
		return
	}

	// Map the gRPC CustomerResponse back to the HTTP response
	cusResp := CustomerHTTPResp{
		ID:      grpcRes.Customer.Id,
		UserId:  grpcRes.Customer.UserId,
		Name:    grpcRes.Customer.Name,
		Email:   grpcRes.Customer.Email,
		Address: grpcRes.Customer.Address,
	}

	err = h.encodeJSON(w, http.StatusCreated, envelope{"customer": cusResp}, nil)
	if err != nil {
		h.serverErrorResponse(w, r, err)
	}
}

func (h *Handler) GetCustomerHandler(w http.ResponseWriter, r *http.Request) {
	// Extract ID param
	customerId, err := h.readIDParam(r)
	if err != nil {
		h.notFoundResponse(w, r)
	}

	// Decode the JSON body into the HTTP request struct
	var httpReq GetCustomerHTTPReq
	err = h.decodeJSON(w, r, &httpReq)
	if err != nil {
		h.badRequestResponse(w, r, err)
		return
	}

	// Convert the HTTP request into the gRPC CreateUserRequest
	grpcReq := &userpb.GetCustomerRequest{
		CustomerId: customerId,
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// Create gRPC connection to user service
	conn, err := grpcutil.ServiceConnection(ctx, "user-service", h.registry)
	if err != nil {
		h.serverErrorResponse(w, r, err)
		return
	}
	defer conn.Close()

	client := userpb.NewUserServiceClient(conn)
	grpcRes, err := client.GetCustomer(ctx, grpcReq)
	if err != nil {
		h.serverErrorResponse(w, r, err)
		return
	}

	// Map the gRPC CustomerResponse back to the HTTP response
	cusResp := CustomerHTTPResp{
		ID:      grpcRes.Customer.Id,
		UserId:  grpcRes.Customer.UserId,
		Name:    grpcRes.Customer.Name,
		Email:   grpcRes.Customer.Email,
		Address: grpcRes.Customer.Address,
	}

	err = h.encodeJSON(w, http.StatusCreated, envelope{"customer": cusResp}, nil)
	if err != nil {
		h.serverErrorResponse(w, r, err)
	}
}

func (h *Handler) UpdateCustomerHandler(w http.ResponseWriter, r *http.Request) {
	// Extract ID param
	customerId, err := h.readIDParam(r)
	if err != nil {
		h.notFoundResponse(w, r)
	}

	// Decode the JSON body into the HTTP request struct
	var httpReq UpdateCustomerHTTPReq
	err = h.decodeJSON(w, r, &httpReq)
	if err != nil {
		h.badRequestResponse(w, r, err)
		return
	}

	// Convert the HTTP request into the gRPC CreateUserRequest
	grpcReq := &userpb.UpdateCustomerRequest{
		CustomerId: customerId,
		Name:       httpReq.Name,
		Email:      httpReq.Email,
		Address:    httpReq.Address,
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// Create gRPC connection to user service
	conn, err := grpcutil.ServiceConnection(ctx, "user-service", h.registry)
	if err != nil {
		h.serverErrorResponse(w, r, err)
		return
	}
	defer conn.Close()

	client := userpb.NewUserServiceClient(conn)
	grpcRes, err := client.UpdateCustomer(ctx, grpcReq)
	if err != nil {
		h.serverErrorResponse(w, r, err)
		return
	}

	// Map the gRPC CustomerResponse back to the HTTP response
	cusResp := CustomerHTTPResp{
		ID:      grpcRes.Customer.Id,
		UserId:  grpcRes.Customer.UserId,
		Name:    grpcRes.Customer.Name,
		Email:   grpcRes.Customer.Email,
		Address: grpcRes.Customer.Address,
	}

	err = h.encodeJSON(w, http.StatusCreated, envelope{"customer": cusResp}, nil)
	if err != nil {
		h.serverErrorResponse(w, r, err)
	}
}

func (h *Handler) DeleteCustomerHandler(w http.ResponseWriter, r *http.Request) {
	// Extract ID param
	customerId, err := h.readIDParam(r)
	if err != nil {
		h.notFoundResponse(w, r)
	}

	// Decode the JSON body into the HTTP request struct
	var httpReq DeleteCustomerHTTPReq
	err = h.decodeJSON(w, r, &httpReq)
	if err != nil {
		h.badRequestResponse(w, r, err)
		return
	}

	// Convert the HTTP request into the gRPC CreateUserRequest
	grpcReq := &userpb.DeleteCustomerRequest{
		CustomerId: customerId,
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// Create gRPC connection to user service
	conn, err := grpcutil.ServiceConnection(ctx, "user-service", h.registry)
	if err != nil {
		h.serverErrorResponse(w, r, err)
		return
	}
	defer conn.Close()

	client := userpb.NewUserServiceClient(conn)
	grpcRes, err := client.DeleteCustomer(ctx, grpcReq)
	if err != nil {
		h.serverErrorResponse(w, r, err)
		return
	}

	// Map the gRPC CustomerResponse back to the HTTP response
	cusResp := DeleteCustomerHTTPResp{
		Message: grpcRes.Message,
	}

	err = h.encodeJSON(w, http.StatusCreated, envelope{"message": cusResp}, nil)
	if err != nil {
		h.serverErrorResponse(w, r, err)
	}
}

// Struct to capture the HTTP request JSON data
type CreateCustomerHTTPReq struct {
	UserId  int64  `json:"user_id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Address string `json:"address"`
}

// Struct to capture the HTTP response
type CustomerHTTPResp struct {
	ID      int64  `json:"customer_id"`
	UserId  int64  `json:"user_id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Address string `json:"address"`
}

// Struct to capture the HTTP request JSON data
type GetCustomerHTTPReq struct {
	ID int64 `json:"customer_id"`
}

// Struct to capture the HTTP request JSON data
type UpdateCustomerHTTPReq struct {
	ID      int64  `json:"customer_id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Address string `json:"address"`
}

// Struct to capture the HTTP request JSON data
type DeleteCustomerHTTPReq struct {
	ID int64 `json:"customer_id"`
}

// Struct to capture the HTTP response
type DeleteCustomerHTTPResp struct {
	Message string `json:"message"`
}
