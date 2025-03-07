package model

import (
	"time"

	"github.com/google/uuid"
)

// TaskStatus represents the status of a task
type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "pending"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusCompleted  TaskStatus = "completed"
	TaskStatusCancelled  TaskStatus = "cancelled"
)

// Task represents a task for an event
type Task struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	EventID      uuid.UUID  `json:"event_id" db:"event_id"`
	Title        string     `json:"title" db:"title"`
	Description  string     `json:"description" db:"description"`
	Status       TaskStatus `json:"status" db:"status"`
	AssignedTo   *uuid.UUID `json:"assigned_to,omitempty" db:"assigned_to"`
	DueDate      *time.Time `json:"due_date,omitempty" db:"due_date"`
	CompletedAt  *time.Time `json:"completed_at,omitempty" db:"completed_at"`
	Priority     int        `json:"priority" db:"priority"` // 1-5, where 5 is highest
	CreatedBy    uuid.UUID  `json:"created_by" db:"created_by"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

// TaskCreateRequest represents the input for creating a new task
type TaskCreateRequest struct {
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	AssignedTo  *uuid.UUID `json:"assigned_to"`
	DueDate     *time.Time `json:"due_date"`
	Priority    int        `json:"priority"`
}

// TaskUpdateRequest represents the input for updating a task
type TaskUpdateRequest struct {
	Title       *string     `json:"title"`
	Description *string     `json:"description"`
	Status      *TaskStatus `json:"status"`
	AssignedTo  *uuid.UUID  `json:"assigned_to"`
	DueDate     *time.Time  `json:"due_date"`
	Priority    *int        `json:"priority"`
}

// TasksResponse represents a list of tasks
type TasksResponse struct {
	Tasks []Task `json:"tasks"`
	Total int    `json:"total"`
}
