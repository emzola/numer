package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/emzola/numer/userservice/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
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

type userRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepo{db}
}

// User management
func (r *userRepo) CreateUser(ctx context.Context, email, password, role string) (models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}

	user := models.User{
		Email:          email,
		HashedPassword: string(hashedPassword),
		Role:           role,
	}

	err = r.db.QueryRowContext(ctx,
		"INSERT INTO users (email, hashed_password, role) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at",
		user.Email, user.HashedPassword, user.Role).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	return user, err
}

func (r *userRepo) GetUserByID(ctx context.Context, userID int64) (models.User, error) {
	var user models.User
	err := r.db.QueryRowContext(ctx,
		"SELECT id, email, hashed_password, role, created_at, updated_at FROM users WHERE id = $1",
		userID).Scan(&user.ID, &user.Email, &user.HashedPassword, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return models.User{}, errors.New("user not found")
	}
	return user, err
}

func (r *userRepo) UpdateUser(ctx context.Context, user models.User) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE users SET email = $1, hashed_password = $2, role = $3, updated_at = NOW() WHERE id = $4",
		user.Email, user.HashedPassword, user.Role, user.ID)
	return err
}

func (r *userRepo) DeleteUser(ctx context.Context, userID int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", userID)
	return err
}

// Customer management
func (r *userRepo) CreateCustomer(ctx context.Context, customer models.Customer) (models.Customer, error) {
	err := r.db.QueryRowContext(ctx,
		"INSERT INTO customers (user_id, name, email, address) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at",
		customer.UserID, customer.Name, customer.Email, customer.Address).
		Scan(&customer.ID, &customer.CreatedAt, &customer.UpdatedAt)

	return customer, err
}

func (r *userRepo) GetCustomerByID(ctx context.Context, customerID int64) (models.Customer, error) {
	var customer models.Customer
	err := r.db.QueryRowContext(ctx,
		"SELECT id, user_id, name, email, address, created_at, updated_at FROM customers WHERE id = $1",
		customerID).Scan(&customer.ID, &customer.UserID, &customer.Name, &customer.Email, &customer.Address, &customer.CreatedAt, &customer.UpdatedAt)
	if err == sql.ErrNoRows {
		return models.Customer{}, errors.New("customer not found")
	}
	return customer, err
}

func (r *userRepo) UpdateCustomer(ctx context.Context, customer models.Customer) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE customers SET name = $1, email = $2, address = $3, updated_at = NOW() WHERE id = $4",
		customer.Name, customer.Email, customer.Address, customer.ID)
	return err
}

func (r *userRepo) DeleteCustomer(ctx context.Context, customerID int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM customers WHERE id = $1", customerID)
	return err
}
