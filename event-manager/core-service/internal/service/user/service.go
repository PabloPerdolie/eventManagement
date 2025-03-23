package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/event-management/core-service/internal/model"
	"github.com/event-management/core-service/internal/repository/user"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// Service provides user-related operations
type Service interface {
	Create(ctx context.Context, req model.UserCreateRequest) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.UserResponse, error)
	GetByEmail(ctx context.Context, email string) (model.UserResponse, error)
	GetByUsername(ctx context.Context, username string) (model.UserResponse, error)
	Update(ctx context.Context, id uuid.UUID, req model.UserUpdateRequest) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, page, size int) (model.UsersResponse, error)
}

type service struct {
	repo   user.Repository
	logger *zap.SugaredLogger
}

// NewService creates a new user service
func NewService(repo user.Repository, logger *zap.SugaredLogger) Service {
	return &service{
		repo:   repo,
		logger: logger,
	}
}

// Create creates a new user
func (s *service) Create(ctx context.Context, req model.UserCreateRequest) (uuid.UUID, error) {
	// Check if email already exists
	_, err := s.repo.GetByEmail(ctx, req.Email)
	if err == nil {
		return uuid.Nil, errors.New("email already registered")
	}

	// Check if username already exists
	_, err = s.repo.GetByUsername(ctx, req.Username)
	if err == nil {
		return uuid.Nil, errors.New("username already taken")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Errorw("Failed to hash password", "error", err)
		return uuid.Nil, errors.New("failed to process password")
	}

	user := model.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		IsActive:     true,
		Role:         req.Role,
	}

	id, err := s.repo.Create(ctx, user)
	if err != nil {
		s.logger.Errorw("Failed to create user", "error", err)
		return uuid.Nil, fmt.Errorf("failed to create user: %w", err)
	}

	return id, nil
}

// GetByID retrieves a user by ID
func (s *service) GetByID(ctx context.Context, id uuid.UUID) (model.UserResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Errorw("Failed to get user by ID", "error", err, "id", id)
		return model.UserResponse{}, fmt.Errorf("failed to get user: %w", err)
	}

	return model.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		IsActive:  user.IsActive,
		Role:      user.Role,
	}, nil
}

// GetByEmail retrieves a user by email
func (s *service) GetByEmail(ctx context.Context, email string) (model.UserResponse, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		s.logger.Errorw("Failed to get user by email", "error", err, "email", email)
		return model.UserResponse{}, fmt.Errorf("failed to get user: %w", err)
	}

	return model.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		IsActive:  user.IsActive,
		Role:      user.Role,
	}, nil
}

// GetByUsername retrieves a user by username
func (s *service) GetByUsername(ctx context.Context, username string) (model.UserResponse, error) {
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		s.logger.Errorw("Failed to get user by username", "error", err, "username", username)
		return model.UserResponse{}, fmt.Errorf("failed to get user: %w", err)
	}

	return model.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		IsActive:  user.IsActive,
		Role:      user.Role,
	}, nil
}

// Update updates a user
func (s *service) Update(ctx context.Context, id uuid.UUID, req model.UserUpdateRequest) error {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Errorw("Failed to get user for update", "error", err, "id", id)
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Update fields if provided
	if req.Username != nil {
		// Check if username is already taken by another user
		existingUser, err := s.repo.GetByUsername(ctx, *req.Username)
		if err == nil && existingUser.ID != id {
			return errors.New("username already taken")
		}
		user.Username = *req.Username
	}

	if req.Email != nil {
		// Check if email is already registered by another user
		existingUser, err := s.repo.GetByEmail(ctx, *req.Email)
		if err == nil && existingUser.ID != id {
			return errors.New("email already registered")
		}
		user.Email = *req.Email
	}

	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}

	if req.LastName != nil {
		user.LastName = *req.LastName
	}

	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if req.Role != nil {
		user.Role = *req.Role
	}

	if err := s.repo.Update(ctx, user); err != nil {
		s.logger.Errorw("Failed to update user", "error", err, "id", id)
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// Delete deletes a user
func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Errorw("Failed to delete user", "error", err, "id", id)
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// List retrieves a list of users with pagination
func (s *service) List(ctx context.Context, page, size int) (model.UsersResponse, error) {
	// Set default pagination values if not provided
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}

	offset := (page - 1) * size
	users, total, err := s.repo.List(ctx, size, offset)
	if err != nil {
		s.logger.Errorw("Failed to list users", "error", err, "page", page, "size", size)
		return model.UsersResponse{}, fmt.Errorf("failed to list users: %w", err)
	}

	// Convert to response objects
	userResponses := make([]model.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = model.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			IsActive:  user.IsActive,
			Role:      user.Role,
		}
	}

	return model.UsersResponse{
		Users: userResponses,
		Total: total,
	}, nil
}
