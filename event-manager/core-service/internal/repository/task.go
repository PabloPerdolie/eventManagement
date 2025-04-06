package repository

import (
	"context"
	"database/sql"
	"github.com/PabloPerdolie/event-manager/core-service/internal/model"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

//
//type Repository interface {
//	Create(ctx context.Context, task model.Task) (int, error)
//	GetById(ctx context.Context, id int) (model.Task, error)
//	Update(ctx context.Context, task model.Task) error
//	Delete(ctx context.Context, id int) error
//	ListByEvent(ctx context.Context, eventId int, limit, offset int) ([]model.Task, int, error)
//	ListByUser(ctx context.Context, userId int, limit, offset int) ([]model.Task, int, error)
//	ListByStatus(ctx context.Context, eventId int, status string, limit, offset int) ([]model.Task, int, error)
//}

type Task struct {
	db *sqlx.DB
}

func NewTask(db *sqlx.DB) Task {
	return Task{
		db: db,
	}
}

func (r Task) Create(ctx context.Context, task model.Task) (int, error) {
	if task.Status == "" {
		task.Status = "pending"
	}

	query := `
        INSERT INTO tasks (event_id, parent_id, title, description, story_points, priority, status, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING task_id
    `

	var taskID int
	err := r.db.QueryRowContext(
		ctx,
		query,
		task.EventId,
		task.ParentId,
		task.Title,
		task.Description,
		task.StoryPoints,
		task.Priority,
		task.Status,
		task.CreatedAt,
	).Scan(&taskID)
	if err != nil {
		return 0, errors.WithMessage(err, "create task")
	}

	return taskID, nil
}

func (r Task) GetById(ctx context.Context, id int) (model.Task, error) {
	var task model.Task
	query := `
        SELECT task_id, event_id, parent_id, title, description, story_points, priority, status, created_at
        FROM tasks
        WHERE task_id = $1
    `

	err := r.db.GetContext(ctx, &task, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Task{}, errors.WithMessage(err, "task not found")
		}
		return model.Task{}, errors.WithMessage(err, "get task")
	}

	return task, nil
}

func (r Task) Update(ctx context.Context, task model.Task) error {
	query := `
        UPDATE tasks
        SET title = $1, description = $2, story_points = $3, priority = $4, status = $5, parent_id = $6
        WHERE task_id = $7
    `

	_, err := r.db.ExecContext(
		ctx,
		query,
		task.Title,
		task.Description,
		task.StoryPoints,
		task.Priority,
		task.Status,
		task.ParentId,
		task.TaskId,
	)
	if err != nil {
		return errors.WithMessage(err, "update task")
	}

	return nil
}

func (r Task) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM tasks WHERE task_id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return errors.WithMessage(err, "delete task")
	}

	return nil
}

func (r Task) ListByEvent(ctx context.Context, eventId, limit, offset int) ([]model.Task, error) {
	var tasks []model.Task

	query := `
        SELECT task_id, event_id, parent_id, title, description, story_points, priority, status, created_at
        FROM tasks
        WHERE event_id = $1
        ORDER BY created_at DESC
        LIMIT $2 OFFSET $3
    `

	err := r.db.SelectContext(ctx, &tasks, query, eventId, limit, offset)
	if err != nil {
		return nil, errors.WithMessage(err, "list event tasks")
	}

	return tasks, nil
}

func (r Task) ListByUser(ctx context.Context, userId, limit, offset int) ([]model.Task, error) {
	var tasks []model.Task

	query := `
        SELECT DISTINCT t.task_id, t.event_id, t.parent_id, t.title, t.description, t.story_points, t.priority, t.status, t.created_at
        FROM tasks t
        JOIN task_assignments ta ON t.task_id = ta.task_id
        WHERE ta.user_id = $1
        ORDER BY t.created_at DESC
        LIMIT $2 OFFSET $3
    `

	err := r.db.SelectContext(ctx, &tasks, query, userId, limit, offset)
	if err != nil {
		return nil, errors.WithMessage(err, "list user tasks")
	}

	return tasks, nil
}

func (r Task) ListByStatus(ctx context.Context, eventId int, status string, limit, offset int) ([]model.Task, error) {
	var tasks []model.Task

	query := `
        SELECT task_id, event_id, parent_id, title, description, story_points, priority, status, created_at
        FROM tasks
        WHERE event_id = $1 AND status = $2
        ORDER BY created_at DESC
        LIMIT $3 OFFSET $4
    `

	err := r.db.SelectContext(ctx, &tasks, query, eventId, status, limit, offset)
	if err != nil {
		return nil, errors.WithMessage(err, "list event tasks by status")
	}

	return tasks, nil
}
