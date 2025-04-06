package repository

import (
	"context"
	"database/sql"
	"github.com/PabloPerdolie/event-manager/core-service/internal/model"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type User struct {
	db *sqlx.DB
}

func NewUser(db *sqlx.DB) User {
	return User{
		db: db,
	}
}

func (r User) GetUserById(ctx context.Context, id int) (*model.User, error) {
	query := `
		SELECT user_id, username, password_hash, email, is_active, is_deleted, role, created_at
		FROM users
		WHERE user_id = $1 AND is_deleted = FALSE
	`
	var user model.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.UserId,
		&user.Username,
		&user.PasswordHash,
		&user.Email,
		&user.IsActive,
		&user.IsDeleted,
		&user.Role,
		&user.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, model.ErrUserNotFound
	}
	if err != nil {
		return nil, errors.WithMessage(err, "get user by id")
	}

	return &user, nil
}

func (r User) ListUsers(ctx context.Context, limit, offset int) ([]model.User, int, error) {
	countQuery := `SELECT COUNT(*) FROM users WHERE is_deleted = FALSE`
	var total int
	err := r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, errors.WithMessage(err, "count users")
	}

	query := `
		SELECT user_id, username, password_hash, email, is_active, is_deleted, role, created_at
		FROM users
		WHERE is_deleted = FALSE
		ORDER BY user_id
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, errors.WithMessage(err, "list users")
	}
	defer rows.Close()

	users := make([]model.User, 0)
	for rows.Next() {
		var user model.User
		err := rows.Scan(
			&user.UserId,
			&user.Username,
			&user.PasswordHash,
			&user.Email,
			&user.IsActive,
			&user.IsDeleted,
			&user.Role,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, 0, errors.WithMessage(err, "scan users")
		}
		users = append(users, user)
	}

	if rows.Err() != nil {
		return nil, 0, errors.WithMessage(err, "rows err")
	}

	return users, total, nil
}
