package repository

import (
	"github.com/event-management/core-service/internal/repository/event"
	"github.com/event-management/core-service/internal/repository/expense"
	"github.com/event-management/core-service/internal/repository/task"
	"github.com/event-management/core-service/internal/repository/user"
	"github.com/jmoiron/sqlx"
)

// Repository is the main repository structure containing all repositories
type Repository struct {
	User             user.Repository
	Event            event.Repository
	EventParticipant event.ParticipantRepository
	Task             task.Repository
	TaskAssignment   task.AssignmentRepository
	Expense          expense.Repository
	ExpenseShare     expense.ShareRepository
}

// New creates a new repository
func New(db *sqlx.DB) *Repository {
	return &Repository{
		User:             user.NewRepository(db),
		Event:            event.NewRepository(db),
		EventParticipant: event.NewParticipantRepository(db),
		Task:             task.NewRepository(db),
		TaskAssignment:   task.NewAssignmentRepository(db),
		Expense:          expense.NewRepository(db),
		ExpenseShare:     expense.NewShareRepository(db),
	}
}
