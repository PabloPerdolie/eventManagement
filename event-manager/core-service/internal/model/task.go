package model

import (
	"time"
)

type Task struct {
	TaskId      int       `db:"task_id"`
	EventId     int       `db:"event_id"`
	ParentId    *int      `db:"parent_id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	StoryPoints *int      `db:"story_points"`
	Priority    *string   `db:"priority"`
	Status      string    `db:"status"`
	CreatedAt   time.Time `db:"created_at"`
}

type TaskAssignment struct {
	TaskAssignmentID int        `db:"task_assignment_id"`
	TaskId           int        `db:"task_id"`
	UserId           int        `db:"user_id"`
	AssignedAt       time.Time  `db:"assigned_at"`
	CompletedAt      *time.Time `db:"completed_at"`
}
