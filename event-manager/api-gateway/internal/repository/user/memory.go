package user

import (
	"context"
	"sync"
	"time"

	"github.com/event-management/api-gateway/internal/domain"
	"github.com/google/uuid"
)

// InMemoryRepository represents an in-memory user repository
type InMemoryRepository struct {
	users map[uuid.UUID]domain.User
	mutex sync.RWMutex
}

// NewInMemoryRepository creates a new in-memory user repository
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		users: make(map[uuid.UUID]domain.User),
	}
}

// CreateUser creates a new user
func (r *InMemoryRepository) CreateUser(ctx context.Context, user domain.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Check if email already exists
	exists, _ := r.CheckEmailExists(ctx, user.Email)
	if exists {
		return domain.ErrEmailAlreadyExists
	}

	r.users[user.ID] = user
	return nil
}

// GetUserByID gets a user by ID
func (r *InMemoryRepository) GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	user, ok := r.users[id]
	if !ok {
		return domain.User{}, domain.ErrUserNotFound
	}

	return user, nil
}

// GetUserByEmail gets a user by email
func (r *InMemoryRepository) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}

	return domain.User{}, domain.ErrUserNotFound
}

// UpdateUser updates a user
func (r *InMemoryRepository) UpdateUser(ctx context.Context, user domain.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.users[user.ID]; !ok {
		return domain.ErrUserNotFound
	}

	// Update user
	user.UpdatedAt = time.Now()
	r.users[user.ID] = user

	return nil
}

// DeleteUser deletes a user
func (r *InMemoryRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.users[id]; !ok {
		return domain.ErrUserNotFound
	}

	delete(r.users, id)
	return nil
}

// ListUsers lists users with pagination
func (r *InMemoryRepository) ListUsers(ctx context.Context, limit, offset int) ([]domain.User, int, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// Get total count
	total := len(r.users)

	// Convert map to slice
	users := make([]domain.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}

	// Apply pagination
	if offset >= len(users) {
		return []domain.User{}, total, nil
	}

	end := offset + limit
	if end > len(users) {
		end = len(users)
	}

	return users[offset:end], total, nil
}

// CheckEmailExists checks if an email exists
func (r *InMemoryRepository) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, user := range r.users {
		if user.Email == email {
			return true, nil
		}
	}

	return false, nil
}
