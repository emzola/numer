package handler

import (
	"context"

	"github.com/emzola/numer/userservice/internal/models"
	"github.com/emzola/numer/userservice/internal/service"
	pb "github.com/emzola/numer/userservice/proto"
)

type UserHandler struct {
	userService service.UserService
	pb.UnimplementedUserServiceServer
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// User Endpoints
func (h *UserHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	user, err := h.userService.CreateUser(ctx, req.Email, req.Password, req.Role)
	if err != nil {
		return nil, err
	}

	return &pb.UserResponse{
		Id:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}

func (h *UserHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	user, err := h.userService.GetUserByID(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	return &pb.UserResponse{
		Id:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	user := models.User{
		ID:    req.UserId,
		Email: req.Email,
		Role:  req.Role,
	}

	if req.Password != "" {
		user.HashedPassword = req.Password
	}

	err := h.userService.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return &pb.UserResponse{
		Id:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}

func (h *UserHandler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.Empty, error) {
	err := h.userService.DeleteUser(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}

// Customer Endpoints
func (h *UserHandler) CreateCustomer(ctx context.Context, req *pb.CreateCustomerRequest) (*pb.CustomerResponse, error) {
	customer := models.Customer{
		UserID:  req.UserId,
		Name:    req.Name,
		Email:   req.Email,
		Address: req.Address,
	}

	customer, err := h.userService.CreateCustomer(ctx, customer)
	if err != nil {
		return nil, err
	}

	return &pb.CustomerResponse{
		Id:      customer.ID,
		UserId:  customer.UserID,
		Name:    customer.Name,
		Email:   customer.Email,
		Address: customer.Address,
	}, nil
}

func (h *UserHandler) GetCustomer(ctx context.Context, req *pb.GetCustomerRequest) (*pb.CustomerResponse, error) {
	customer, err := h.userService.GetCustomerByID(ctx, req.CustomerId)
	if err != nil {
		return nil, err
	}

	return &pb.CustomerResponse{
		Id:      customer.ID,
		UserId:  customer.UserID,
		Name:    customer.Name,
		Email:   customer.Email,
		Address: customer.Address,
	}, nil
}

func (h *UserHandler) UpdateCustomer(ctx context.Context, req *pb.UpdateCustomerRequest) (*pb.CustomerResponse, error) {
	customer := models.Customer{
		ID:      req.CustomerId,
		Name:    req.Name,
		Email:   req.Email,
		Address: req.Address,
	}

	err := h.userService.UpdateCustomer(ctx, customer)
	if err != nil {
		return nil, err
	}

	return &pb.CustomerResponse{
		Id:      customer.ID,
		UserId:  customer.UserID,
		Name:    customer.Name,
		Email:   customer.Email,
		Address: customer.Address,
	}, nil
}

func (h *UserHandler) DeleteCustomer(ctx context.Context, req *pb.DeleteCustomerRequest) (*pb.Empty, error) {
	err := h.userService.DeleteCustomer(ctx, req.CustomerId)
	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}
