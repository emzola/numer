package models

import "time"

type User struct {
	ID             int64     `db:"id"`
	Email          string    `db:"email"`
	HashedPassword string    `db:"hashed_password"`
	Role           string    `db:"role"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}
