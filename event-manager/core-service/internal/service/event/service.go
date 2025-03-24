package event

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/PabloPerdolie/event-manager/core-service/internal/model"
	"github.com/PabloPerdolie/event-manager/core-service/internal/repository/event"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Service provides event-related operations
type Service interface {
	Create(ctx context.Context, userID uuid.UUID, req model.EventCreateRequest) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.EventResponse, error)
	Update(ctx context.Context, id uuid.UUID, req model.EventUpdateRequest) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, page, size int) (model.EventsResponse, error)
	ListByOrganizer(ctx context.Context, organizerID uuid.UUID, page, size int) (model.EventsResponse, error)
	ListByParticipant(ctx context.Context, participantID uuid.UUID, page, size int) (model.EventsResponse, error)
}

type service struct {
	eventRepo       event.Repository
	participantRepo event.ParticipantRepository
	logger          *zap.SugaredLogger
}

// NewService creates a new event service
func NewService(eventRepo event.Repository, participantRepo event.ParticipantRepository, logger *zap.SugaredLogger) Service {
	return &service{
		eventRepo:       eventRepo,
		participantRepo: participantRepo,
		logger:          logger,
	}
}

// Create creates a new event
func (s *service) Create(ctx context.Context, userID uuid.UUID, req model.EventCreateRequest) (uuid.UUID, error) {
	event := model.Event{
		Name:        req.Name,
		Description: req.Description,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Location:    req.Location,
		CreatedBy:   userID,
	}

	id, err := s.eventRepo.Create(ctx, event)
	if err != nil {
		s.logger.Errorw("Failed to create event", "error", err, "userId", userID)
		return uuid.Nil, fmt.Errorf("failed to create event: %w", err)
	}

	// Automatically add event creator as a participant
	participant := model.EventParticipant{
		EventID:     id,
		UserID:      userID,
		IsConfirmed: true,
	}

	_, err = s.participantRepo.Create(ctx, participant)
	if err != nil {
		s.logger.Errorw("Failed to add creator as participant", "error", err, "eventId", id, "userId", userID)
		// We don't return an error here, as the event was already created successfully
	}

	return id, nil
}

// GetByID retrieves an event by ID
func (s *service) GetByID(ctx context.Context, id uuid.UUID) (model.EventResponse, error) {
	event, err := s.eventRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Errorw("Failed to get event by ID", "error", err, "id", id)
		return model.EventResponse{}, fmt.Errorf("failed to get event: %w", err)
	}

	return model.EventResponse{
		ID:          event.ID,
		Name:        event.Name,
		Description: event.Description,
		StartDate:   event.StartDate,
		EndDate:     event.EndDate,
		Location:    event.Location,
		CreatedBy:   event.CreatedBy,
		CreatedAt:   event.CreatedAt,
		UpdatedAt:   event.UpdatedAt,
	}, nil
}

// Update updates an event
func (s *service) Update(ctx context.Context, id uuid.UUID, req model.EventUpdateRequest) error {
	event, err := s.eventRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Errorw("Failed to get event for update", "error", err, "id", id)
		return fmt.Errorf("failed to get event: %w", err)
	}

	// Update fields if provided
	if req.Name != nil {
		event.Name = *req.Name
	}

	if req.Description != nil {
		event.Description = *req.Description
	}

	if req.StartDate != nil {
		event.StartDate = *req.StartDate
	}

	if req.EndDate != nil {
		event.EndDate = req.EndDate
	}

	if req.Location != nil {
		event.Location = *req.Location
	}

	if err := s.eventRepo.Update(ctx, event); err != nil {
		s.logger.Errorw("Failed to update event", "error", err, "id", id)
		return fmt.Errorf("failed to update event: %w", err)
	}

	return nil
}

// Delete deletes an event
func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	// First check if the event exists
	_, err := s.eventRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Errorw("Failed to get event for deletion", "error", err, "id", id)
		return fmt.Errorf("failed to get event: %w", err)
	}

	if err := s.eventRepo.Delete(ctx, id); err != nil {
		s.logger.Errorw("Failed to delete event", "error", err, "id", id)
		return fmt.Errorf("failed to delete event: %w", err)
	}

	return nil
}

// List retrieves a list of events with pagination
func (s *service) List(ctx context.Context, page, size int) (model.EventsResponse, error) {
	// Set default pagination values if not provided
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}

	offset := (page - 1) * size
	events, total, err := s.eventRepo.List(ctx, size, offset)
	if err != nil {
		s.logger.Errorw("Failed to list events", "error", err, "page", page, "size", size)
		return model.EventsResponse{}, fmt.Errorf("failed to list events: %w", err)
	}

	// Convert to response objects
	eventResponses := make([]model.EventResponse, len(events))
	for i, event := range events {
		eventResponses[i] = model.EventResponse{
			ID:          event.ID,
			Name:        event.Name,
			Description: event.Description,
			StartDate:   event.StartDate,
			EndDate:     event.EndDate,
			Location:    event.Location,
			CreatedBy:   event.CreatedBy,
			CreatedAt:   event.CreatedAt,
			UpdatedAt:   event.UpdatedAt,
		}
	}

	return model.EventsResponse{
		Events: eventResponses,
		Total:  total,
	}, nil
}

// ListByOrganizer retrieves a list of events for a specific organizer with pagination
func (s *service) ListByOrganizer(ctx context.Context, organizerID uuid.UUID, page, size int) (model.EventsResponse, error) {
	// Set default pagination values if not provided
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}

	offset := (page - 1) * size
	events, total, err := s.eventRepo.ListByOrganizer(ctx, organizerID, size, offset)
	if err != nil {
		s.logger.Errorw("Failed to list organizer events", "error", err, "organizerId", organizerID, "page", page, "size", size)
		return model.EventsResponse{}, fmt.Errorf("failed to list organizer events: %w", err)
	}

	// Convert to response objects
	eventResponses := make([]model.EventResponse, len(events))
	for i, event := range events {
		eventResponses[i] = model.EventResponse{
			ID:          event.ID,
			Name:        event.Name,
			Description: event.Description,
			StartDate:   event.StartDate,
			EndDate:     event.EndDate,
			Location:    event.Location,
			CreatedBy:   event.CreatedBy,
			CreatedAt:   event.CreatedAt,
			UpdatedAt:   event.UpdatedAt,
		}
	}

	return model.EventsResponse{
		Events: eventResponses,
		Total:  total,
	}, nil
}

// ListByParticipant retrieves a list of events for a specific participant with pagination
func (s *service) ListByParticipant(ctx context.Context, participantID uuid.UUID, page, size int) (model.EventsResponse, error) {
	// Set default pagination values if not provided
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}

	offset := (page - 1) * size
	events, total, err := s.eventRepo.ListByParticipant(ctx, participantID, size, offset)
	if err != nil {
		s.logger.Errorw("Failed to list participant events", "error", err, "participantId", participantID, "page", page, "size", size)
		return model.EventsResponse{}, fmt.Errorf("failed to list participant events: %w", err)
	}

	// Convert to response objects
	eventResponses := make([]model.EventResponse, len(events))
	for i, event := range events {
		eventResponses[i] = model.EventResponse{
			ID:          event.ID,
			Name:        event.Name,
			Description: event.Description,
			StartDate:   event.StartDate,
			EndDate:     event.EndDate,
			Location:    event.Location,
			CreatedBy:   event.CreatedBy,
			CreatedAt:   event.CreatedAt,
			UpdatedAt:   event.UpdatedAt,
		}
	}

	return model.EventsResponse{
		Events: eventResponses,
		Total:  total,
	}, nil
}
