package task

import (
	"context"
	"database/sql"
	"time"

	"github.com/PabloPerdolie/event-manager/core-service/internal/model"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

//type AssignmentRepository interface {
//	Create(ctx context.Context, assignment model.TaskAssignment) (int, error)
//	GetById(ctx context.Context, id int) (model.TaskAssignment, error)
//	GetByTaskAndUser(ctx context.Context, taskId, userId int) (model.TaskAssignment, error)
//	Update(ctx context.Context, assignment model.TaskAssignment) error
//	Delete(ctx context.Context, id int) error
//	ListByTask(ctx context.Context, taskId int, limit, offset int) ([]model.TaskAssignment, int, error)
//	ListByUser(ctx context.Context, userId int, limit, offset int) ([]model.TaskAssignment, int, error)
//}

type TaskAssignment struct {
	db *sqlx.DB
}

func NewTaskAssignment(db *sqlx.DB) TaskAssignment {
	return TaskAssignment{
		db: db,
	}
}

func (r *TaskAssignment) Create(ctx context.Context, assignment model.TaskAssignment) (int, error) {
	assignment.AssignedAt = time.Now()

	query := `
        INSERT INTO task_assignments (task_id, user_id, assigned_at)
        VALUES ($1, $2, $3)
        RETURNING task_assignment_id
    `

	var assignmentID int
	err := r.db.QueryRowContext(
		ctx,
		query,
		assignment.TaskId,
		assignment.UserId,
		assignment.AssignedAt,
	).Scan(&assignmentID)
	if err != nil {
		return 0, errors.WithMessage(err, "create task assignment")
	}

	return assignmentID, nil
}

func (r *TaskAssignment) GetByTaskAndUser(ctx context.Context, taskId, userId int) (model.TaskAssignment, error) {
	var assignment model.TaskAssignment

	query := `
        SELECT task_assignment_id, task_id, user_id, assigned_at, completed_at
        FROM task_assignments
        WHERE task_id = $1 AND user_id = $2
    `

	err := r.db.GetContext(ctx, &assignment, query, taskId, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.TaskAssignment{}, errors.WithMessage(err, "task assignment not found")
		}
		return model.TaskAssignment{}, errors.WithMessage(err, "get task assignment")
	}

	return assignment, nil
}

func (r *TaskAssignment) Update(ctx context.Context, assignment model.TaskAssignment) error {
	query := `
        UPDATE task_assignments
        SET completed_at = $1
        WHERE task_id = $2
    `

	_, err := r.db.ExecContext(
		ctx,
		query,
		assignment.CompletedAt,
		assignment.TaskId,
	)
	if err != nil {
		return errors.WithMessage(err, "update task assignment")
	}

	return nil
}

func (r *TaskAssignment) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM task_assignments WHERE user_id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return errors.WithMessage(err, "delete task assignment")
	}

	return nil
}

func (r *TaskAssignment) ListByTask(ctx context.Context, taskId, limit, offset int) ([]model.TaskAssignment, error) {
	var assignments []model.TaskAssignment

	query := `
        SELECT task_assignment_id, task_id, user_id, assigned_at, completed_at
        FROM task_assignments
        WHERE task_id = $1
        ORDER BY assigned_at DESC
        LIMIT $2 OFFSET $3
    `

	err := r.db.SelectContext(ctx, &assignments, query, taskId, limit, offset)
	if err != nil {
		return nil, errors.WithMessage(err, "list task assignments")
	}

	return assignments, nil
}

func (r *TaskAssignment) ListByUser(ctx context.Context, userId, limit, offset int) ([]model.TaskAssignment, error) {
	var assignments []model.TaskAssignment

	query := `
        SELECT task_assignment_id, task_id, user_id, assigned_at, completed_at
        FROM task_assignments
        WHERE user_id = $1
        ORDER BY assigned_at DESC
        LIMIT $2 OFFSET $3
    `

	err := r.db.SelectContext(ctx, &assignments, query, userId, limit, offset)
	if err != nil {
		return nil, errors.WithMessage(err, "list user task assignments")
	}

	return assignments, nil
}
