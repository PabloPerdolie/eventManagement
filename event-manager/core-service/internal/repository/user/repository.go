package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/event-management/core-service/internal/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// Repository defines user repository interface
type Repository interface {
	Create(ctx context.Context, user model.User) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.User, error)
	GetByEmail(ctx context.Context, email string) (model.User, error)
	GetByUsername(ctx context.Context, username string) (model.User, error)
	Update(ctx context.Context, user model.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]model.User, int, error)
}

type repository struct {
	db *sqlx.DB
}

// NewRepository creates a new user repository
func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

// Create creates a new user in the database
func (r *repository) Create(ctx context.Context, user model.User) (uuid.UUID, error) {
	user.ID = uuid.New()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	query := `
		INSERT INTO users (id, username, email, password_hash, first_name, last_name, created_at, updated_at, is_active, role)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		user.ID,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
		user.CreatedAt,
		user.UpdatedAt,
		user.IsActive,
		user.Role,
	)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user.ID, nil
}

// GetByID retrieves a user by ID
func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (model.User, error) {
	var user model.User

	query := `
		SELECT id, username, email, password_hash, first_name, last_name, created_at, updated_at, is_active, role
		FROM users
		WHERE id = $1
	`

	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, fmt.Errorf("user not found: %w", err)
		}
		return model.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetByEmail retrieves a user by email
func (r *repository) GetByEmail(ctx context.Context, email string) (model.User, error) {
	var user model.User

	query := `
		SELECT id, username, email, password_hash, first_name, last_name, created_at, updated_at, is_active, role
		FROM users
		WHERE email = $1
	`

	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, fmt.Errorf("user not found: %w", err)
		}
		return model.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetByUsername retrieves a user by username
func (r *repository) GetByUsername(ctx context.Context, username string) (model.User, error) {
	var user model.User

	query := `
		SELECT id, username, email, password_hash, first_name, last_name, created_at, updated_at, is_active, role
		FROM users
		WHERE username = $1
	`

	err := r.db.GetContext(ctx, &user, query, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, fmt.Errorf("user not found: %w", err)
		}
		return model.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// Update updates a user in the database
func (r *repository) Update(ctx context.Context, user model.User) error {
	user.UpdatedAt = time.Now()

	query := `
		UPDATE users
		SET username = $1, email = $2, first_name = $3, last_name = $4, updated_at = $5, is_active = $6, role = $7
		WHERE id = $8
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		user.Username,
		user.Email,
		user.FirstName,
		user.LastName,
		user.UpdatedAt,
		user.IsActive,
		user.Role,
		user.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// Delete deletes a user from the database
func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// List retrieves a list of users with pagination
func (r *repository) List(ctx context.Context, limit, offset int) ([]model.User, int, error) {
	var users []model.User
	var total int

	// Count total users
	countQuery := `SELECT COUNT(*) FROM users`
	err := r.db.GetContext(ctx, &total, countQuery)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Retrieve users with pagination
	query := `
		SELECT id, username, email, password_hash, first_name, last_name, created_at, updated_at, is_active, role
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	err = r.db.SelectContext(ctx, &users, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	return users, total, nil
}
