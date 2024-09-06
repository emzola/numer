package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/emzola/numer/userservice/pkg/model"
)

var (
	ErrNotFound = errors.New("not found")
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user model.User) error {
	query := `INSERT INTO users (id, name, email, role, active) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query, user.ID, user.Name, user.Email, user.Role, user.Active)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	query := `SELECT id, name, email, role, active FROM users WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)
	var user model.User
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.Active); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}
