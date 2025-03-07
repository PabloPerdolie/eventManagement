package user

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/event-management/api-gateway/internal/domain"
	"go.uber.org/zap"
)

// Repository определяет методы для работы с пользователями
type Repository interface {
	GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	ListUsers(ctx context.Context, limit, offset int) ([]domain.User, int, error)
	UpdateUser(ctx context.Context, user domain.User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	CheckEmailExists(ctx context.Context, email string) (bool, error)
}

// Service предоставляет методы для управления пользователями
type Service struct {
	repo   Repository
	logger *zap.SugaredLogger
}

// New создает новый сервис пользователей
func New(repo Repository, logger *zap.SugaredLogger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

// GetUserByID получает пользователя по ID
func (s *Service) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.UserResponse, error) {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &domain.UserResponse{
		ID:              user.ID,
		Email:           user.Email,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		Role:            user.Role,
		IsEmailVerified: user.IsEmailVerified,
		CreatedAt:       user.CreatedAt,
	}, nil
}

// ListUsers получает список пользователей с пагинацией
func (s *Service) ListUsers(ctx context.Context, limit, offset int) ([]domain.UserResponse, int, error) {
	users, total, err := s.repo.ListUsers(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	userResponses := make([]domain.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = domain.UserResponse{
			ID:              user.ID,
			Email:           user.Email,
			FirstName:       user.FirstName,
			LastName:        user.LastName,
			Role:            user.Role,
			IsEmailVerified: user.IsEmailVerified,
			CreatedAt:       user.CreatedAt,
		}
	}

	return userResponses, total, nil
}

// UpdateUser обновляет профиль пользователя
func (s *Service) UpdateUser(ctx context.Context, userID uuid.UUID, req domain.UserUpdateRequest) (*domain.UserResponse, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Обновляем email, если указан
	if req.Email != nil && *req.Email != "" && *req.Email != user.Email {
		// Проверяем, не занят ли email
		existingUser, err := s.repo.GetUserByEmail(ctx, *req.Email)
		if err == nil && existingUser.ID != userID {
			return nil, domain.ErrEmailAlreadyExists
		}

		user.Email = *req.Email
	}

	// Обновляем имя, если указано
	if req.FirstName != nil && *req.FirstName != "" {
		user.FirstName = *req.FirstName
	}

	// Обновляем фамилию, если указана
	if req.LastName != nil && *req.LastName != "" {
		user.LastName = *req.LastName
	}

	// Обновляем время последнего обновления
	user.UpdatedAt = time.Now()

	// Сохраняем обновления
	if err := s.repo.UpdateUser(ctx, user); err != nil {
		return nil, err
	}

	return &domain.UserResponse{
		ID:              user.ID,
		Email:           user.Email,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		Role:            user.Role,
		IsEmailVerified: user.IsEmailVerified,
		CreatedAt:       user.CreatedAt,
	}, nil
}

// DeactivateUser устанавливает статус is_active пользователя в false (мягкое удаление)
func (s *Service) DeactivateUser(ctx context.Context, userID uuid.UUID) error {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	user.IsActive = false
	user.UpdatedAt = time.Now()

	return s.repo.UpdateUser(ctx, user)
}

// HardDeleteUser полностью удаляет пользователя из системы
func (s *Service) HardDeleteUser(ctx context.Context, userID uuid.UUID) error {
	return s.repo.DeleteUser(ctx, userID)
}
