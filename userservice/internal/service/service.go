package service

import (
	"context"
	"errors"

	"github.com/emzola/numer/userservice/internal/repository"
	"github.com/emzola/numer/userservice/pkg/model"
	"github.com/google/uuid"
)

var (
	ErrNotFound               = errors.New("not found")
	ErrInvalidRequest         = errors.New("the request is invalid")
	ErrInvoiceNumberIncrement = errors.New("failed to increment invoice number")
)

// UserRepository defines the interface for interacting with the user data store.
type UserRepository interface {
	CreateUser(ctx context.Context, user model.User) error
	GetUserByID(ctx context.Context, id string) (*model.User, error)
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, user model.User) (*model.User, error) {
	// Generate user ID
	user.ID = uuid.New().String()

	// Activate user. In production, set to true only after email verification
	user.Active = true

	// Save to repository
	err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserService) GetUser(ctx context.Context, id string) (*model.User, error) {
	// Call repository to fetch the user by its ID
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return user, nil
}
