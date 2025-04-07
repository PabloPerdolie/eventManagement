package event

import (
	"context"
	"github.com/PabloPerdolie/event-manager/core-service/internal/domain"
	"github.com/PabloPerdolie/event-manager/core-service/internal/model"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type EventRepo interface {
	Create(ctx context.Context, event model.Event) (int, error)
	GetById(ctx context.Context, eventID int) (model.Event, error)
	Update(ctx context.Context, event model.Event) error
	Delete(ctx context.Context, eventID int) error
	List(ctx context.Context, limit, offset int) ([]model.Event, error)
	ListByOrganizer(ctx context.Context, organizerID, limit, offset int) ([]model.Event, error)
	ListByParticipant(ctx context.Context, participantID, limit, offset int) ([]model.Event, error)
}

type EventParticipantRepo interface {
	Create(ctx context.Context, participant model.EventParticipant) (int, error)
}

type NotifyPublisher interface {
	Publish(ctx context.Context, data []byte) error
}

type Service struct {
	eventRepo       EventRepo
	participantRepo EventParticipantRepo
	notifyPbl       NotifyPublisher
	logger          *zap.SugaredLogger
}

func NewService(eventRepo EventRepo, participantRepo EventParticipantRepo, notifyPbl NotifyPublisher, logger *zap.SugaredLogger) Service {
	return Service{
		eventRepo:       eventRepo,
		participantRepo: participantRepo,
		notifyPbl:       notifyPbl,
		logger:          logger,
	}
}

func (s Service) Create(ctx context.Context, userId int, req domain.EventCreateRequest) (int, error) {
	event := model.Event{
		Title:       req.Title,
		Description: req.Description,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Location:    &req.Location,
		OrganizerID: userId,
	}

	id, err := s.eventRepo.Create(ctx, event)
	if err != nil {
		s.logger.Errorw("Failed to create event", "error", err, "userId", userId)
		return 0, errors.WithMessage(err, "create event")
	}

	participant := model.EventParticipant{
		EventID:     id,
		UserID:      userId,
		IsConfirmed: ptr(true),
		Role:        model.RoleOrganizer,
	}

	_, err = s.participantRepo.Create(ctx, participant)
	if err != nil {
		s.logger.Errorw("Failed to add creator as participant", "error", err, "eventId", id, "userId", userId)
		// We don't return an error here, as the event was already created successfully // todo
	}

	return id, nil
}

func (s Service) GetById(ctx context.Context, id int) (*domain.EventResponse, error) {
	event, err := s.eventRepo.GetById(ctx, id)
	if err != nil {
		s.logger.Errorw("Failed to get event by Id", "error", err, "id", id)
		return nil, errors.WithMessage(err, "get event")
	}

	return &domain.EventResponse{
		Id:          event.EventId,
		Title:       event.Title,
		Description: event.Description,
		StartDate:   event.StartDate,
		EndDate:     event.EndDate,
		Location:    *event.Location,
		CreatedBy:   event.OrganizerID,
		CreatedAt:   event.CreatedAt,
	}, nil
}

func (s Service) Update(ctx context.Context, id int, req domain.EventUpdateRequest) error {
	event, err := s.eventRepo.GetById(ctx, id)
	if err != nil {
		s.logger.Errorw("Failed to get event for update", "error", err, "id", id)
		return errors.WithMessage(err, "get event")
	}

	if req.Title != nil {
		event.Title = *req.Title
	}

	if req.Description != nil {
		event.Description = *req.Description
	}

	if req.StartDate != nil {
		event.StartDate = *req.StartDate
	}

	if req.EndDate != nil {
		event.EndDate = *req.EndDate
	}

	if req.Location != nil {
		event.Location = req.Location
	}

	if err := s.eventRepo.Update(ctx, event); err != nil {
		s.logger.Errorw("Failed to update event", "error", err, "id", id)
		return errors.WithMessage(err, "update event")
	}

	return nil
}

func (s Service) Delete(ctx context.Context, id int) error {
	_, err := s.eventRepo.GetById(ctx, id)
	if err != nil {
		s.logger.Errorw("Failed to get event for deletion", "error", err, "id", id)
		return errors.WithMessage(err, "get event")
	}

	if err := s.eventRepo.Delete(ctx, id); err != nil {
		s.logger.Errorw("Failed to delete event", "error", err, "id", id)
		return errors.WithMessage(err, "delete event")
	}

	return nil
}

func (s Service) List(ctx context.Context, page, size int) (*domain.EventsResponse, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}

	offset := (page - 1) * size
	events, err := s.eventRepo.List(ctx, size, offset)
	if err != nil {
		s.logger.Errorw("Failed to list events", "error", err, "page", page, "size", size)
		return nil, errors.WithMessage(err, "list events")
	}

	return convertToEventsResponse(events), nil
}

func (s Service) ListByOrganizer(ctx context.Context, organizerId int, page, size int) (*domain.EventsResponse, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}

	offset := (page - 1) * size
	events, err := s.eventRepo.ListByOrganizer(ctx, organizerId, size, offset)
	if err != nil {
		s.logger.Errorw("Failed to list organizer events", "error", err, "organizerId", organizerId, "page", page, "size", size)
		return nil, errors.WithMessage(err, "list organizer events")
	}

	return convertToEventsResponse(events), nil
}

func (s Service) ListByParticipant(ctx context.Context, participantId int, page, size int) (*domain.EventsResponse, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}

	offset := (page - 1) * size
	events, err := s.eventRepo.ListByParticipant(ctx, participantId, size, offset)
	if err != nil {
		s.logger.Errorw("Failed to list participant events", "error", err, "participantId", participantId, "page", page, "size", size)
		return nil, errors.WithMessage(err, "list participant events")
	}

	return convertToEventsResponse(events), nil
}

func convertToEventsResponse(events []model.Event) *domain.EventsResponse {
	eventResponses := make([]domain.EventResponse, len(events))

	for i, event := range events {
		eventResponses[i] = domain.EventResponse{
			Id:          event.EventId,
			Title:       event.Title,
			Description: event.Description,
			StartDate:   event.StartDate,
			EndDate:     event.EndDate,
			Location:    *event.Location,
			CreatedBy:   event.OrganizerID,
			CreatedAt:   event.CreatedAt,
		}
	}

	return &domain.EventsResponse{
		Events: eventResponses,
		Total:  len(events),
	}
}
