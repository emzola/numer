package models

import "time"

type User struct {
	ID             int64
	Email          string
	HashedPassword string
	Role           string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
