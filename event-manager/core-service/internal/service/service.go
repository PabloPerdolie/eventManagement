package service

import (
	"context"
	"github.com/PabloPerdolie/event-manager/core-service/internal/domain"
	"github.com/PabloPerdolie/event-manager/core-service/internal/model"
	"github.com/pkg/errors"
)

type ParticipantService interface {
	ListByEvent(ctx context.Context, eventID int, page, size int) (*domain.EventParticipantsResponse, error)
}

type EventService interface {
	GetById(ctx context.Context, id int) (*domain.EventResponse, error)
}

type TaskService interface {
	ListByEvent(ctx context.Context, eventId int, page, size int) (*domain.TasksResponse, error)
}

type CommentsRepo interface {
	GetByEventId(_ context.Context, eventId int) (*model.CommunicationServiceResponse, error)
}

type ExpenseService interface {
	ListExpensesByEvent(ctx context.Context, eventId int, page, size int) (domain.ExpensesResponse, error)
	GetEventBalanceReport(ctx context.Context, eventId int) (domain.BalanceReportResponse, error)
}

type Service struct {
	taskService        TaskService
	participantService ParticipantService
	eventService       EventService
	commentsRepo       CommentsRepo
	expenseService     ExpenseService
}

func NewService(taskService TaskService, participantService ParticipantService, eventService EventService, commentsRepo CommentsRepo, expenseService ExpenseService) Service {
	return Service{
		taskService:        taskService,
		participantService: participantService,
		eventService:       eventService,
		commentsRepo:       commentsRepo,
		expenseService:     expenseService,
	}
}

func (s Service) GetEventSummary(ctx context.Context, eventId int) (*domain.EventData, error) {
	event, err := s.eventService.GetById(ctx, eventId)
	if err != nil {
		return nil, errors.WithMessage(err, "get event by id")
	}

	tasks, err := s.taskService.ListByEvent(ctx, eventId, 1, 100) // todo delete pagination
	if err != nil {
		return nil, errors.WithMessage(err, "get tasks by event id")
	}

	participants, err := s.participantService.ListByEvent(ctx, eventId, 1, 100) // todo delete pagination
	if err != nil {
		return nil, errors.WithMessage(err, "get participants by event id")
	}

	comments, err := s.commentsRepo.GetByEventId(ctx, eventId)
	if err != nil {
		return nil, errors.WithMessage(err, "get comments by event id")
	}

	expenses, err := s.expenseService.ListExpensesByEvent(ctx, eventId, 1, 100)
	if err != nil {
		return nil, errors.WithMessage(err, "get expenses by event id")
	}

	balanceReport, err := s.expenseService.GetEventBalanceReport(ctx, eventId)
	if err != nil {
		return nil, errors.WithMessage(err, "get balance report by event id")
	}

	return &domain.EventData{
		EventParticipants: *participants,
		EventData:         *event,
		Tasks:             *tasks,
		Comments:          *comments,
		Expenses:          expenses,
		BalanceReport:     balanceReport,
	}, nil
}
