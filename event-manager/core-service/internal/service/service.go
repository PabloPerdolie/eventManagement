package service

import (
	"github.com/PabloPerdolie/event-manager/core-service/internal/repository"
	"github.com/PabloPerdolie/event-manager/core-service/internal/service/event"
	"github.com/PabloPerdolie/event-manager/core-service/internal/service/expense"
	"github.com/PabloPerdolie/event-manager/core-service/internal/service/task"
	"github.com/PabloPerdolie/event-manager/core-service/internal/service/user"
	"go.uber.org/zap"
)

// Service contains all the services for the application
type Service struct {
	User             user.Service
	Event            event.Service
	EventParticipant event.ParticipantService
	Task             task.Service
	TaskAssignment   task.AssignmentService
	Expense          expense.Service
	ExpenseShare     expense.ShareService
	logger           *zap.SugaredLogger
}

// New creates a new service
func New(repo *repository.Repository, logger *zap.SugaredLogger) *Service {
	return &Service{
		User:             user.NewService(repo.User, logger),
		Event:            event.NewService(repo.Event, repo.EventParticipant, logger),
		EventParticipant: event.NewParticipantService(repo.EventParticipant, repo.User, logger),
		Task:             task.NewService(repo.Task, repo.TaskAssignment, logger),
		TaskAssignment:   task.NewAssignmentService(repo.TaskAssignment, repo.Task, repo.User, logger),
		Expense:          expense.NewService(repo.Expense, repo.ExpenseShare, logger),
		ExpenseShare:     expense.NewShareService(repo.ExpenseShare, repo.Expense, repo.User, logger),
		logger:           logger,
	}
}
