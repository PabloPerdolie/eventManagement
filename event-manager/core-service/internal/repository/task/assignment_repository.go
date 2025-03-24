package task

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/PabloPerdolie/event-manager/core-service/internal/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// AssignmentRepository defines task assignment repository interface
type AssignmentRepository interface {
	Create(ctx context.Context, assignment model.TaskAssignment) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.TaskAssignment, error)
	GetByTaskAndUser(ctx context.Context, taskID, userID uuid.UUID) (model.TaskAssignment, error)
	Update(ctx context.Context, assignment model.TaskAssignment) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByTask(ctx context.Context, taskID uuid.UUID, limit, offset int) ([]model.TaskAssignment, int, error)
	ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.TaskAssignment, int, error)
}

type assignmentRepository struct {
	db *sqlx.DB
}

// NewAssignmentRepository creates a new task assignment repository
func NewAssignmentRepository(db *sqlx.DB) AssignmentRepository {
	return &assignmentRepository{db: db}
}

// Create creates a new task assignment in the database
func (r *assignmentRepository) Create(ctx context.Context, assignment model.TaskAssignment) (uuid.UUID, error) {
	assignment.ID = uuid.New()
	assignment.AssignedAt = time.Now()

	query := `
		INSERT INTO task_assignments (id, task_id, user_id, assigned_at, completed_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		assignment.ID,
		assignment.TaskID,
		assignment.UserID,
		assignment.AssignedAt,
		assignment.CompletedAt,
	)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create task assignment: %w", err)
	}

	return assignment.ID, nil
}

// GetByID retrieves a task assignment by ID
func (r *assignmentRepository) GetByID(ctx context.Context, id uuid.UUID) (model.TaskAssignment, error) {
	var assignment model.TaskAssignment

	query := `
		SELECT id, task_id, user_id, assigned_at, completed_at
		FROM task_assignments
		WHERE id = $1
	`

	err := r.db.GetContext(ctx, &assignment, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.TaskAssignment{}, fmt.Errorf("task assignment not found: %w", err)
		}
		return model.TaskAssignment{}, fmt.Errorf("failed to get task assignment: %w", err)
	}

	return assignment, nil
}

// GetByTaskAndUser retrieves a task assignment by task ID and user ID
func (r *assignmentRepository) GetByTaskAndUser(ctx context.Context, taskID, userID uuid.UUID) (model.TaskAssignment, error) {
	var assignment model.TaskAssignment

	query := `
		SELECT id, task_id, user_id, assigned_at, completed_at
		FROM task_assignments
		WHERE task_id = $1 AND user_id = $2
	`

	err := r.db.GetContext(ctx, &assignment, query, taskID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.TaskAssignment{}, fmt.Errorf("task assignment not found: %w", err)
		}
		return model.TaskAssignment{}, fmt.Errorf("failed to get task assignment: %w", err)
	}

	return assignment, nil
}

// Update updates a task assignment in the database
func (r *assignmentRepository) Update(ctx context.Context, assignment model.TaskAssignment) error {
	query := `
		UPDATE task_assignments
		SET completed_at = $1
		WHERE id = $2
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		assignment.CompletedAt,
		assignment.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update task assignment: %w", err)
	}

	return nil
}

// Delete deletes a task assignment from the database
func (r *assignmentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM task_assignments WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task assignment: %w", err)
	}

	return nil
}

// ListByTask retrieves a list of task assignments for a specific task with pagination
func (r *assignmentRepository) ListByTask(ctx context.Context, taskID uuid.UUID, limit, offset int) ([]model.TaskAssignment, int, error) {
	var assignments []model.TaskAssignment
	var total int

	// Count total assignments for the task
	countQuery := `SELECT COUNT(*) FROM task_assignments WHERE task_id = $1`
	err := r.db.GetContext(ctx, &total, countQuery, taskID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count task assignments: %w", err)
	}

	// Retrieve assignments with pagination
	query := `
		SELECT id, task_id, user_id, assigned_at, completed_at
		FROM task_assignments
		WHERE task_id = $1
		ORDER BY assigned_at DESC
		LIMIT $2 OFFSET $3
	`

	err = r.db.SelectContext(ctx, &assignments, query, taskID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list task assignments: %w", err)
	}

	return assignments, total, nil
}

// ListByUser retrieves a list of task assignments for a specific user with pagination
func (r *assignmentRepository) ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.TaskAssignment, int, error) {
	var assignments []model.TaskAssignment
	var total int

	// Count total assignments for the user
	countQuery := `SELECT COUNT(*) FROM task_assignments WHERE user_id = $1`
	err := r.db.GetContext(ctx, &total, countQuery, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count user task assignments: %w", err)
	}

	// Retrieve assignments with pagination
	query := `
		SELECT id, task_id, user_id, assigned_at, completed_at
		FROM task_assignments
		WHERE user_id = $1
		ORDER BY assigned_at DESC
		LIMIT $2 OFFSET $3
	`

	err = r.db.SelectContext(ctx, &assignments, query, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list user task assignments: %w", err)
	}

	return assignments, total, nil
}
