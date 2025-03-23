package model

import (
	"time"

	"github.com/google/uuid"
)

// UserRole represents the role of a user
type UserRole string

const (
	UserRoleOrganizer   UserRole = "organizer"
	UserRoleParticipant UserRole = "participant"
)

// User represents a user entity
type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	FirstName    string    `json:"first_name" db:"first_name"`
	LastName     string    `json:"last_name" db:"last_name"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	Role         UserRole  `json:"role" db:"role"`
}

// UserCreateRequest represents the input for creating a new user
type UserCreateRequest struct {
	Username  string   `json:"username" binding:"required"`
	Email     string   `json:"email" binding:"required,email"`
	Password  string   `json:"password" binding:"required,min=8"`
	FirstName string   `json:"first_name" binding:"required"`
	LastName  string   `json:"last_name" binding:"required"`
	Role      UserRole `json:"role" binding:"required"`
}

// UserUpdateRequest represents the input for updating a user
type UserUpdateRequest struct {
	Username  *string   `json:"username"`
	Email     *string   `json:"email" binding:"omitempty,email"`
	FirstName *string   `json:"first_name"`
	LastName  *string   `json:"last_name"`
	IsActive  *bool     `json:"is_active"`
	Role      *UserRole `json:"role"`
}

// UserResponse represents the output for user data
type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IsActive  bool      `json:"is_active"`
	Role      UserRole  `json:"role"`
}

// UsersResponse represents a list of users
type UsersResponse struct {
	Users []UserResponse `json:"users"`
	Total int            `json:"total"`
}
