package model

import (
	"github.com/event-management/api-gateway/internal/domain"
	"time"

	"github.com/google/uuid"
)


type User struct {
	ID              uuid.UUID `json:"id" db:"id"`
	Email           string    `json:"email" db:"email"`
	PasswordHash    string    `json:"-" db:"password_hash"`
	FirstName       string    `json:"first_name" db:"first_name"`
	LastName        string    `json:"last_name" db:"last_name"`
	Role            domain.UserRole  `json:"role" db:"role"`
	IsActive        bool      `json:"is_active" db:"is_active"`
	IsEmailVerified bool      `json:"is_email_verified" db:"is_email_verified"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
	LastLoginAt     *time.Time `json:"last_login_at,omitempty" db:"last_login_at"`
}
