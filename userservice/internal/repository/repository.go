package repository

import (
	"context"
	"database/sql"

	"github.com/emzola/numer/userservice/pkg/model"
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
