package model

import (
	"time"

	"github.com/google/uuid"
)

// TaskAssignment represents an assignment of a task to a user
type TaskAssignment struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	TaskID      uuid.UUID  `json:"task_id" db:"task_id"`
	UserID      uuid.UUID  `json:"user_id" db:"user_id"`
	AssignedAt  time.Time  `json:"assigned_at" db:"assigned_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty" db:"completed_at"`
}

// TaskAssignmentCreateRequest represents the input for creating a new task assignment
type TaskAssignmentCreateRequest struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
}

// TaskAssignmentUpdateRequest represents the input for updating a task assignment
type TaskAssignmentUpdateRequest struct {
	CompletedAt *time.Time `json:"completed_at"`
}

// TaskAssignmentResponse represents the output for task assignment data
type TaskAssignmentResponse struct {
	ID          uuid.UUID    `json:"id"`
	Task        TaskResponse `json:"task"`
	User        UserResponse `json:"user"`
	AssignedAt  time.Time    `json:"assigned_at"`
	CompletedAt *time.Time   `json:"completed_at,omitempty"`
}

// TaskAssignmentsResponse represents a list of task assignments
type TaskAssignmentsResponse struct {
	Assignments []TaskAssignmentResponse `json:"assignments"`
	Total       int                      `json:"total"`
}

// TaskResponse represents the output for task data
type TaskResponse struct {
	ID          uuid.UUID  `json:"id"`
	EventID     uuid.UUID  `json:"event_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	Priority    int        `json:"priority"`
	CreatedBy   uuid.UUID  `json:"created_by"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
