package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/emzola/numer/gateway-service/internal/grpcutil"
	userpb "github.com/emzola/numer/user-service/proto"
)

func (h *Handler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Decode the JSON body into the HTTP request struct
	var httpReq CreateUserHTTPReq
	err := h.decodeJSON(w, r, &httpReq)
	if err != nil {
		h.badRequestResponse(w, r, err)
		return
	}

	// Convert the HTTP request into the gRPC CreateUserRequest
	grpcReq := &userpb.CreateUserRequest{
		Email:    httpReq.Email,
		Password: httpReq.Password,
		Role:     httpReq.Role,
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
	grpcRes, err := client.CreateUser(ctx, grpcReq)
	if err != nil {
		h.serverErrorResponse(w, r, err)
		return
	}

	// Map the gRPC UserResponse back to the HTTP response
	userResp := UserHTTPResp{
		ID:    grpcRes.User.Id,
		Email: grpcRes.User.Email,
		Role:  grpcRes.User.Role,
	}

	err = h.encodeJSON(w, http.StatusCreated, envelope{"user": userResp}, nil)
	if err != nil {
		h.serverErrorResponse(w, r, err)
	}
}

func (h *Handler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	// Extract ID param
	userId, err := h.readIDParam(r)
	if err != nil {
		h.notFoundResponse(w, r)
		return
	}

	// Decode the JSON body into the HTTP request struct
	var httpReq GetUserHTTPReq
	err = h.decodeJSON(w, r, &httpReq)
	if err != nil {
		h.badRequestResponse(w, r, err)
		return
	}

	// Convert the HTTP request into the gRPC GetUserRequest
	grpcReq := &userpb.GetUserRequest{
		UserId: userId,
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
	grpcRes, err := client.GetUser(ctx, grpcReq)
	if err != nil {
		h.serverErrorResponse(w, r, err)
		return
	}

	// Map the gRPC UserResponse back to the HTTP response
	userResp := UserHTTPResp{
		ID:    grpcRes.User.Id,
		Email: grpcRes.User.Email,
		Role:  grpcRes.User.Role,
	}

	err = h.encodeJSON(w, http.StatusOK, envelope{"user": userResp}, nil)
	if err != nil {
		h.serverErrorResponse(w, r, err)
	}
}

func (h *Handler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Extract user from context
	user := h.contextGetUser(r)

	// Extract ID param
	userId, err := h.readIDParam(r)
	if err != nil {
		h.notFoundResponse(w, r)
		return
	}

	// Basic validation to ensure user only updates own profile
	if user.Id != userId {
		h.notPermittedResponse(w, r)
		return
	}

	// Decode the JSON body into the HTTP request struct
	var httpReq UpdateUserHTTPReq
	err = h.decodeJSON(w, r, &httpReq)
	if err != nil {
		h.badRequestResponse(w, r, err)
		return
	}

	// Convert the HTTP request into the gRPC CreateUserRequest
	grpcReq := &userpb.UpdateUserRequest{
		UserId:   user.Id,
		Email:    httpReq.Email,
		Password: httpReq.Password,
		Role:     httpReq.Role,
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
	grpcRes, err := client.UpdateUser(ctx, grpcReq)
	if err != nil {
		h.serverErrorResponse(w, r, err)
		return
	}

	// Map the gRPC UserResponse back to the HTTP response
	userResp := UserHTTPResp{
		ID:    grpcRes.User.Id,
		Email: grpcRes.User.Email,
		Role:  grpcRes.User.Role,
	}

	err = h.encodeJSON(w, http.StatusOK, envelope{"user": userResp}, nil)
	if err != nil {
		h.serverErrorResponse(w, r, err)
	}
}

func (h *Handler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	// Extract user from context
	user := h.contextGetUser(r)

	// Extract ID param
	userId, err := h.readIDParam(r)
	if err != nil {
		h.notFoundResponse(w, r)
	}

	// Basic validation to ensure user only deletes own profile
	if user.Id != userId {
		h.notPermittedResponse(w, r)
		return
	}

	// Decode the JSON body into the HTTP request struct
	var httpReq DeleteUserHTTPReq
	err = h.decodeJSON(w, r, &httpReq)
	if err != nil {
		h.badRequestResponse(w, r, err)
		return
	}

	// Convert the HTTP request into the gRPC CreateUserRequest
	grpcReq := &userpb.DeleteUserRequest{
		UserId: user.Id,
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
	grpcRes, err := client.DeleteUser(ctx, grpcReq)
	if err != nil {
		h.serverErrorResponse(w, r, err)
		return
	}

	// Map the gRPC UserResponse back to the HTTP response
	userResp := DeleteUserHTTPResp{
		Message: grpcRes.Message,
	}

	err = h.encodeJSON(w, http.StatusOK, envelope{"message": userResp}, nil)
	if err != nil {
		h.serverErrorResponse(w, r, err)
	}
}

// Struct to capture the HTTP request JSON data
type CreateUserHTTPReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// Struct to capture the HTTP response
type UserHTTPResp struct {
	ID        int64  `json:"user_id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}

// Struct to capture the HTTP request JSON data
type GetUserHTTPReq struct {
	ID int64 `json:"user_id"`
}

// Struct to capture the HTTP request JSON data
type UpdateUserHTTPReq struct {
	ID       int64  `json:"user_id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// Struct to capture the HTTP request JSON data
type DeleteUserHTTPReq struct {
	ID int64 `json:"user_id"`
}

// Struct to capture the HTTP response
type DeleteUserHTTPResp struct {
	Message string `json:"message"`
}
