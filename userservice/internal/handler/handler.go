package handler

import (
	"context"
	"errors"
	"time"

	userpb "github.com/emzola/numer/userservice/genproto"
	"github.com/emzola/numer/userservice/internal/service"
	"github.com/emzola/numer/userservice/pkg/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserServiceServer is the gRPC server for the User service.
type UserServiceServer struct {
	userpb.UnimplementedUserServiceServer
	service *service.UserService
}

func NewUserServiceServer(service *service.UserService) *UserServiceServer {
	return &UserServiceServer{service: service}
}

func (s *UserServiceServer) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.UserResponse, error) {
	if req == nil || req.Email == "" {
		return nil, status.Errorf(codes.InvalidArgument, service.ErrInvalidRequest.Error())
	}

	user := model.User{
		Name:  req.Name,
		Email: req.Email,
		Role:  req.Role,
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	createdUser, err := s.service.CreateUser(ctx, user)
	if err != nil {
		code, errMsg := mapToGRPCErrorCode(err), err.Error()
		return nil, status.Errorf(code, errMsg)
	}

	return &userpb.UserResponse{User: model.ConvertUserToProto(createdUser)}, nil
}

func (s *UserServiceServer) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.UserResponse, error) {
	if req == nil || req.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, service.ErrInvalidRequest.Error())
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	user, err := s.service.GetUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &userpb.UserResponse{User: model.ConvertUserToProto(user)}, nil
}

// mapToGRPCErrorCode maps domain-specific errors to gRPC status codes.
func mapToGRPCErrorCode(err error) codes.Code {
	switch {
	case errors.Is(err, service.ErrNotFound):
		return codes.NotFound
	case errors.Is(err, service.ErrInvalidRequest):
		return codes.InvalidArgument
	default:
		return codes.Internal
	}
}
