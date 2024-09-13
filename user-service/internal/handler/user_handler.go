package handler

import (
	"context"

	"github.com/emzola/numer/user-service/internal/models"
	"github.com/emzola/numer/user-service/internal/service"
	pb "github.com/emzola/numer/user-service/proto"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	userService *service.UserService
	pb.UnimplementedUserServiceServer
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// User Endpoints
func (h *UserHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	user, err := h.userService.CreateUser(ctx, req.Email, req.Password, req.Role)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UserResponse{User: models.ConvertUserToProto(user)}, nil
}

func (h *UserHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	user, err := h.userService.GetUserByID(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UserResponse{User: models.ConvertUserToProto(user)}, nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	user := &models.User{
		ID:    req.UserId,
		Email: req.Email,
		Role:  req.Role,
	}

	if req.Password != "" {
		user.HashedPassword = req.Password
	}

	err := h.userService.UpdateUser(ctx, user)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UserResponse{User: models.ConvertUserToProto(user)}, nil
}

func (h *UserHandler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	err := h.userService.DeleteUser(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.DeleteUserResponse{Message: "user successfully deleted"}, nil
}

func (h *UserHandler) AuthenticateUser(ctx context.Context, req *pb.AuthenticateUserRequest) (*pb.AuthenticateUserResponse, error) {
	user, err := h.userService.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(req.Password))
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return &pb.AuthenticateUserResponse{
		Valid:  true,
		UserId: user.ID,
		Email:  user.Email,
		Role:   user.Role,
	}, nil
}

// Customer Endpoints
func (h *UserHandler) CreateCustomer(ctx context.Context, req *pb.CreateCustomerRequest) (*pb.CustomerResponse, error) {
	customer := &models.Customer{
		UserID:  req.UserId,
		Name:    req.Name,
		Email:   req.Email,
		Address: req.Address,
	}

	customer, err := h.userService.CreateCustomer(ctx, customer)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CustomerResponse{Customer: models.ConvertCustomerToProto(customer)}, nil
}

func (h *UserHandler) GetCustomer(ctx context.Context, req *pb.GetCustomerRequest) (*pb.CustomerResponse, error) {
	customer, err := h.userService.GetCustomerByID(ctx, req.CustomerId)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &pb.CustomerResponse{Customer: models.ConvertCustomerToProto(customer)}, nil
}

func (h *UserHandler) UpdateCustomer(ctx context.Context, req *pb.UpdateCustomerRequest) (*pb.CustomerResponse, error) {
	customer := &models.Customer{
		ID:      req.CustomerId,
		Name:    req.Name,
		Email:   req.Email,
		Address: req.Address,
	}

	err := h.userService.UpdateCustomer(ctx, customer)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CustomerResponse{Customer: models.ConvertCustomerToProto(customer)}, nil
}

func (h *UserHandler) DeleteCustomer(ctx context.Context, req *pb.DeleteCustomerRequest) (*pb.DeleteCustomerResponse, error) {
	err := h.userService.DeleteCustomer(ctx, req.CustomerId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.DeleteCustomerResponse{Message: "customer successfully deleted"}, nil
}
