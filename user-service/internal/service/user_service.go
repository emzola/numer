package service

import (
	"context"

	"github.com/emzola/numer/user-service/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type userRepository interface {
	// User management methods
	CreateUser(ctx context.Context, email, password, role string) (*models.User, error)
	GetUserByID(ctx context.Context, userID int64) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, userID int64) error

	// Customer management methods
	CreateCustomer(ctx context.Context, customer *models.Customer) (*models.Customer, error)
	GetCustomerByID(ctx context.Context, customerID int64) (*models.Customer, error)
	UpdateCustomer(ctx context.Context, customer *models.Customer) error
	DeleteCustomer(ctx context.Context, customerID int64) error
}

type UserService struct {
	repo userRepository
}

func NewUserService(repo userRepository) *UserService {
	return &UserService{repo: repo}
}

// User Management
func (s *UserService) CreateUser(ctx context.Context, email, password, role string) (*models.User, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return &models.User{}, err
	}

	// Create user in repository
	user, err := s.repo.CreateUser(ctx, email, string(hashedPassword), role)
	return user, err
}

func (s *UserService) GetUserByID(ctx context.Context, userID int64) (*models.User, error) {
	return s.repo.GetUserByID(ctx, userID)
}

func (s *UserService) UpdateUser(ctx context.Context, user *models.User) error {
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

func (s *UserService) DeleteUser(ctx context.Context, userID int64) error {
	return s.repo.DeleteUser(ctx, userID)
}

// Customer Management
func (s *UserService) CreateCustomer(ctx context.Context, customer *models.Customer) (*models.Customer, error) {
	return s.repo.CreateCustomer(ctx, customer)
}

func (s *UserService) GetCustomerByID(ctx context.Context, customerID int64) (*models.Customer, error) {
	return s.repo.GetCustomerByID(ctx, customerID)
}

func (s *UserService) UpdateCustomer(ctx context.Context, customer *models.Customer) error {
	return s.repo.UpdateCustomer(ctx, customer)
}

func (s *UserService) DeleteCustomer(ctx context.Context, customerID int64) error {
	return s.repo.DeleteCustomer(ctx, customerID)
}
