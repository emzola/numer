package service

import (
	"context"

	"github.com/emzola/numer/userservice/internal/models"
	"github.com/emzola/numer/userservice/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	// User management methods
	CreateUser(ctx context.Context, email, password, role string) (models.User, error)
	GetUserByID(ctx context.Context, userID int64) (models.User, error)
	UpdateUser(ctx context.Context, user models.User) error
	DeleteUser(ctx context.Context, userID int64) error

	// Customer management methods
	CreateCustomer(ctx context.Context, customer models.Customer) (models.Customer, error)
	GetCustomerByID(ctx context.Context, customerID int64) (models.Customer, error)
	UpdateCustomer(ctx context.Context, customer models.Customer) error
	DeleteCustomer(ctx context.Context, customerID int64) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

// User Management
func (s *userService) CreateUser(ctx context.Context, email, password, role string) (models.User, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}

	// Create user in repository
	user, err := s.repo.CreateUser(ctx, email, string(hashedPassword), role)
	return user, err
}

func (s *userService) GetUserByID(ctx context.Context, userID int64) (models.User, error) {
	return s.repo.GetUserByID(ctx, userID)
}

func (s *userService) UpdateUser(ctx context.Context, user models.User) error {
	// Hash password if provided
	if user.HashedPassword != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.HashedPassword), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.HashedPassword = string(hashedPassword)
	}
	return s.repo.UpdateUser(ctx, user)
}

func (s *userService) DeleteUser(ctx context.Context, userID int64) error {
	return s.repo.DeleteUser(ctx, userID)
}

// Customer Management
func (s *userService) CreateCustomer(ctx context.Context, customer models.Customer) (models.Customer, error) {
	return s.repo.CreateCustomer(ctx, customer)
}

func (s *userService) GetCustomerByID(ctx context.Context, customerID int64) (models.Customer, error) {
	return s.repo.GetCustomerByID(ctx, customerID)
}

func (s *userService) UpdateCustomer(ctx context.Context, customer models.Customer) error {
	return s.repo.UpdateCustomer(ctx, customer)
}

func (s *userService) DeleteCustomer(ctx context.Context, customerID int64) error {
	return s.repo.DeleteCustomer(ctx, customerID)
}
