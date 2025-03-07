package domain

import (
	"errors"
)

// Common errors
var (
	// user errors
	ErrUserNotFound       = errors.New("user not found")
	ErrUserNotActive      = errors.New("user is not active")
	ErrUserSameEmail      = errors.New("new email is the same as current")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")

	// Token errors
	ErrInvalidToken      = errors.New("invalid token")
	ErrTokenExpired      = errors.New("token expired")
	ErrInvalidResetToken = errors.New("invalid or expired reset token")

	// Auth errors
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")

	// Repository errors
	ErrDuplicateKey = errors.New("duplicate key")
)
