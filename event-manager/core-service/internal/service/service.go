package service

import (
	"context"
	"github.com/PabloPerdolie/event-manager/core-service/internal/domain"
	"github.com/PabloPerdolie/event-manager/core-service/internal/model"
	"github.com/pkg/errors"
)

type ParticipantService interface {
	Create(ctx context.Context, eventID int, req domain.EventParticipantCreateRequest) (*domain.EventParticipantResponse, error)
	GetById(ctx context.Context, id int) (*domain.EventParticipantResponse, error)
	Delete(ctx context.Context, id int) error
	ListByEvent(ctx context.Context, eventID int, page, size int) (*domain.EventParticipantsResponse, error)
	ListByUser(ctx context.Context, userId int, page, size int) (*domain.EventParticipantsResponse, error)
	ConfirmParticipation(ctx context.Context, id int) error
	DeclineParticipation(ctx context.Context, id int) error
}

type EventService interface {
	Create(ctx context.Context, userId int, req domain.EventCreateRequest) (int, error)
	GetById(ctx context.Context, id int) (*domain.EventResponse, error)
	Update(ctx context.Context, id int, req domain.EventUpdateRequest) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, page, size int) (*domain.EventsResponse, error)
	ListByOrganizer(ctx context.Context, organizerId int, page, size int) (*domain.EventsResponse, error)
	ListByParticipant(ctx context.Context, participantId int, page, size int) (*domain.EventsResponse, error)
}

type TaskService interface {
	Create(ctx context.Context, req domain.TaskCreateRequest) (*domain.TaskResponse, error)
	Update(ctx context.Context, id int, req domain.TaskUpdateRequest) error
	Delete(ctx context.Context, id int) error
	ListByEvent(ctx context.Context, eventId int, page, size int) (*domain.TasksResponse, error)
	ListByUser(ctx context.Context, userId int, page, size int) (*domain.TasksResponse, error)
	UpdateStatus(ctx context.Context, id int, status domain.TaskStatus) error
}

type CommentsRepo interface {
	GetByEventId(_ context.Context, eventId int) (*model.CommunicationServiceResponse, error)
}

type Service struct {
	taskService        TaskService
	participantService ParticipantService
	eventService       EventService
	commentsRepo       CommentsRepo
}

func NewService(taskService TaskService, participantService ParticipantService, eventService EventService, commentsRepo CommentsRepo) Service {
	return Service{
		taskService:        taskService,
		participantService: participantService,
		eventService:       eventService,
		commentsRepo:       commentsRepo,
	}
}

func (s Service) GetEventSummary(ctx context.Context, eventId int) (*domain.EventData, error) {
	event, err := s.eventService.GetById(ctx, eventId)
	if err != nil {
		return nil, errors.WithMessage(err, "get event by id")
	}

	tasks, err := s.taskService.ListByEvent(ctx, eventId, 1, 10) // todo delete pagination
	if err != nil {
		return nil, errors.WithMessage(err, "get tasks by event id")
	}

	participants, err := s.participantService.ListByEvent(ctx, eventId, 1, 10) // todo delete pagination
	if err != nil {
		return nil, errors.WithMessage(err, "get participants by event id")
	}

	comments, err := s.commentsRepo.GetByEventId(ctx, eventId)
	if err != nil {
		return nil, errors.WithMessage(err, "get comments by event id")
	}

	return &domain.EventData{
		EventParticipants: *participants,
		EventData:         *event,
		Tasks:             *tasks,
		Comments:          *comments,
	}, nil
}
