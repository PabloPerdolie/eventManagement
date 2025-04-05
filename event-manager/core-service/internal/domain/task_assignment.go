package domain

import "time"

type TaskAssignmentCreateRequest struct {
	UserId int `json:"user_id" binding:"required"`
	TaskId int `json:"task_id" binging:"required"`
}

type TaskAssignmentUpdateRequest struct {
	CompletedAt *time.Time `json:"completed_at"`
}

type TaskAssignmentResponse struct {
	Id          int          `json:"id"`
	Task        TaskResponse `json:"task"`
	User        UserResponse `json:"user"`
	AssignedAt  time.Time    `json:"assigned_at"`
	CompletedAt *time.Time   `json:"completed_at,omitempty"`
}

type TaskAssignmentsResponse struct {
	Assignments []TaskAssignmentResponse `json:"assignments"`
	Total       int                      `json:"total"`
}
