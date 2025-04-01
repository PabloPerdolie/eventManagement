package domain

import (
	"errors"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserNotActive      = errors.New("user is not active")
	ErrUserSameEmail      = errors.New("new email is the same as current")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")

	ErrInvalidToken      = errors.New("invalid token")
	ErrTokenExpired      = errors.New("token expired")
	ErrInvalidResetToken = errors.New("invalid or expired reset token")

	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")

	ErrDuplicateKey = errors.New("duplicate key")
)
