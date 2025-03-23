package task

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

// Repository defines task repository interface
type Repository interface {
	Create(ctx context.Context, task model.Task) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.Task, error)
	Update(ctx context.Context, task model.Task) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByEvent(ctx context.Context, eventID uuid.UUID, limit, offset int) ([]model.Task, int, error)
	ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.Task, int, error)
	ListByStatus(ctx context.Context, eventID uuid.UUID, status model.TaskStatus, limit, offset int) ([]model.Task, int, error)
}

type repository struct {
	db *sqlx.DB
}

// NewRepository creates a new task repository
func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

// Create creates a new task in the database
func (r *repository) Create(ctx context.Context, task model.Task) (uuid.UUID, error) {
	task.ID = uuid.New()
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()
	
	// Set default status if not provided
	if task.Status == "" {
		task.Status = model.TaskStatusPending
	}

	query := `
		INSERT INTO tasks (id, event_id, title, description, status, assigned_to, due_date, completed_at, priority, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		task.ID,
		task.EventID,
		task.Title,
		task.Description,
		task.Status,
		task.AssignedTo,
		task.DueDate,
		task.CompletedAt,
		task.Priority,
		task.CreatedBy,
		task.CreatedAt,
		task.UpdatedAt,
	)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create task: %w", err)
	}

	return task.ID, nil
}

// GetByID retrieves a task by ID
func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (model.Task, error) {
	var task model.Task

	query := `
		SELECT id, event_id, title, description, status, assigned_to, due_date, completed_at, priority, created_by, created_at, updated_at
		FROM tasks
		WHERE id = $1
	`

	err := r.db.GetContext(ctx, &task, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Task{}, fmt.Errorf("task not found: %w", err)
		}
		return model.Task{}, fmt.Errorf("failed to get task: %w", err)
	}

	return task, nil
}

// Update updates a task in the database
func (r *repository) Update(ctx context.Context, task model.Task) error {
	task.UpdatedAt = time.Now()

	query := `
		UPDATE tasks
		SET title = $1, description = $2, status = $3, assigned_to = $4, due_date = $5, 
		    completed_at = $6, priority = $7, updated_at = $8
		WHERE id = $9
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		task.Title,
		task.Description,
		task.Status,
		task.AssignedTo,
		task.DueDate,
		task.CompletedAt,
		task.Priority,
		task.UpdatedAt,
		task.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	return nil
}

// Delete deletes a task from the database
func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM tasks WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	return nil
}

// ListByEvent retrieves a list of tasks for a specific event with pagination
func (r *repository) ListByEvent(ctx context.Context, eventID uuid.UUID, limit, offset int) ([]model.Task, int, error) {
	var tasks []model.Task
	var total int

	// Count total tasks for the event
	countQuery := `SELECT COUNT(*) FROM tasks WHERE event_id = $1`
	err := r.db.GetContext(ctx, &total, countQuery, eventID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count event tasks: %w", err)
	}

	// Retrieve tasks with pagination
	query := `
		SELECT id, event_id, title, description, status, assigned_to, due_date, completed_at, priority, created_by, created_at, updated_at
		FROM tasks
		WHERE event_id = $1
		ORDER BY due_date ASC NULLS LAST, priority DESC, created_at DESC
		LIMIT $2 OFFSET $3
	`

	err = r.db.SelectContext(ctx, &tasks, query, eventID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list event tasks: %w", err)
	}

	return tasks, total, nil
}

// ListByUser retrieves a list of tasks assigned to a specific user with pagination
func (r *repository) ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.Task, int, error) {
	var tasks []model.Task
	var total int

	// Count total tasks for the user
	countQuery := `SELECT COUNT(*) FROM tasks WHERE assigned_to = $1`
	err := r.db.GetContext(ctx, &total, countQuery, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count user tasks: %w", err)
	}

	// Retrieve tasks with pagination
	query := `
		SELECT id, event_id, title, description, status, assigned_to, due_date, completed_at, priority, created_by, created_at, updated_at
		FROM tasks
		WHERE assigned_to = $1
		ORDER BY due_date ASC NULLS LAST, priority DESC, created_at DESC
		LIMIT $2 OFFSET $3
	`

	err = r.db.SelectContext(ctx, &tasks, query, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list user tasks: %w", err)
	}

	return tasks, total, nil
}

// ListByStatus retrieves a list of tasks for a specific event and status with pagination
func (r *repository) ListByStatus(ctx context.Context, eventID uuid.UUID, status model.TaskStatus, limit, offset int) ([]model.Task, int, error) {
	var tasks []model.Task
	var total int

	// Count total tasks for the event and status
	countQuery := `SELECT COUNT(*) FROM tasks WHERE event_id = $1 AND status = $2`
	err := r.db.GetContext(ctx, &total, countQuery, eventID, status)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count event tasks by status: %w", err)
	}

	// Retrieve tasks with pagination
	query := `
		SELECT id, event_id, title, description, status, assigned_to, due_date, completed_at, priority, created_by, created_at, updated_at
		FROM tasks
		WHERE event_id = $1 AND status = $2
		ORDER BY due_date ASC NULLS LAST, priority DESC, created_at DESC
		LIMIT $3 OFFSET $4
	`

	err = r.db.SelectContext(ctx, &tasks, query, eventID, status, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list event tasks by status: %w", err)
	}

	return tasks, total, nil
}
