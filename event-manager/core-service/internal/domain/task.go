package domain

import "time"

type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusCompleted  TaskStatus = "completed"
	TaskStatusCancelled  TaskStatus = "cancelled"
)

type TaskCreateRequest struct {
	EventId     int     `json:"event_id" binding:"required"`
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description"`
	ParentId    *int    `json:"parent_id"`
	StoryPoints *int    `json:"story_points"`
	Priority    *string `json:"priority"`
	AssignedTo  *int    `json:"assigned_to"`
}

type TaskUpdateRequest struct {
	Title       *string     `json:"title"`
	Description *string     `json:"description"`
	ParentId    *int        `json:"parent_id"`
	StoryPoints *int        `json:"story_points"`
	Priority    *string     `json:"priority"`
	Status      *TaskStatus `json:"status"`
	AssignedTo  *int        `json:"assigned_to"`
}

type TaskResponse struct {
	Id          int        `json:"id"`
	EventId     int        `json:"event_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	ParentID    *int       `json:"parent_id,omitempty"`
	StoryPoints *int       `json:"story_points,omitempty"`
	Priority    *string    `json:"priority,omitempty"`
	Status      TaskStatus `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	AssignedTo  *int       `json:"assigned_to,omitempty"`
}

type TasksResponse struct {
	Tasks []TaskResponse `json:"tasks"`
	Total int            `json:"total"`
}
