package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/PabloPerdolie/event-manager/api-gateway/internal/domain"
	"github.com/PabloPerdolie/event-manager/api-gateway/internal/model"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type PostgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) CreateUser(ctx context.Context, user model.User) (int, error) {
	query := `
		INSERT INTO users (username, password_hash, email, is_active, is_deleted, role, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING user_id
	`

	var id int
	err := r.db.QueryRowContext(ctx, query,
		user.Username,
		user.PasswordHash,
		user.Email,
		user.IsActive,
		user.IsDeleted,
		user.Role,
		time.Now(),
	).Scan(&id)
	if err != nil {
		return 0, errors.WithMessage(err, "create user")
	}

	return id, nil
}

func (r *PostgresRepository) GetUserById(ctx context.Context, id int) (*model.User, error) {
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
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, errors.WithMessage(err, "get user by id")
	}

	return &user, nil
}

func (r *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT user_id, username, password_hash, email, is_active, is_deleted, role, created_at
		FROM users
		WHERE email = $1 AND is_deleted = FALSE
	`

	var user model.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
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
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, errors.WithMessage(err, "get user by email")
	}

	return &user, nil
}

func (r *PostgresRepository) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	query := `
		SELECT user_id, username, password_hash, email, is_active, is_deleted, role, created_at
		FROM users
		WHERE username = $1 AND is_deleted = FALSE
	`

	var user model.User
	err := r.db.QueryRowContext(ctx, query, username).Scan(
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
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, errors.WithMessage(err, "get user by username")
	}

	return &user, nil
}

func (r *PostgresRepository) UpdateUser(ctx context.Context, user model.User) error {
	query := `
		UPDATE users
		SET username = $1,
			password_hash = $2,
			email = $3,
			is_active = $4,
			is_deleted = $5,
			role = $6
		WHERE user_id = $7 AND is_deleted = FALSE
	`

	res, err := r.db.ExecContext(ctx, query,
		user.Username,
		user.PasswordHash,
		user.Email,
		user.IsActive,
		user.IsDeleted,
		user.Role,
		user.UserId,
	)
	if err != nil {
		return errors.WithMessage(err, "update user")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return errors.WithMessage(err, "scan result rows affected")
	}

	if rowsAffected == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (r *PostgresRepository) DeleteUser(ctx context.Context, id int) error {
	query := `
		UPDATE users
		SET is_deleted = TRUE
		WHERE user_id = $1 AND is_deleted = FALSE
	`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return errors.WithMessage(err, "delete user")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return errors.WithMessage(err, "scan result rows affected")
	}
	if rowsAffected == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (r *PostgresRepository) ListUsers(ctx context.Context, limit, offset int) ([]model.User, int, error) {
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
