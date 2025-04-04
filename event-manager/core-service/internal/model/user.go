package model

import (
	"time"
)

type User struct {
	UserId       int       `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Role         string    `json:"role" db:"role"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	IsDeleted    bool      `json:"is_deleted" db:"is_deleted"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}
