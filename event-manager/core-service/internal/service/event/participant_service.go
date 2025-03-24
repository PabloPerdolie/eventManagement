package event

import (
	"context"
	"errors"
	"fmt"

	"github.com/PabloPerdolie/event-manager/core-service/internal/model"
	"github.com/PabloPerdolie/event-manager/core-service/internal/repository/event"
	"github.com/PabloPerdolie/event-manager/core-service/internal/repository/user"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// ParticipantService provides event participant-related operations
type ParticipantService interface {
	Create(ctx context.Context, req model.EventParticipantCreateRequest) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.EventParticipantResponse, error)
	Update(ctx context.Context, id uuid.UUID, req model.EventParticipantUpdateRequest) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByEvent(ctx context.Context, eventID uuid.UUID, page, size int) (model.EventParticipantsResponse, error)
	ListByUser(ctx context.Context, userID uuid.UUID, page, size int) (model.EventParticipantsResponse, error)
	ConfirmParticipation(ctx context.Context, id uuid.UUID) error
	DeclineParticipation(ctx context.Context, id uuid.UUID) error
}

type participantService struct {
	repo     event.ParticipantRepository
	userRepo user.Repository
	logger   *zap.SugaredLogger
}

// NewParticipantService creates a new event participant service
func NewParticipantService(repo event.ParticipantRepository, userRepo user.Repository, logger *zap.SugaredLogger) ParticipantService {
	return &participantService{
		repo:     repo,
		userRepo: userRepo,
		logger:   logger,
	}
}

// Create creates a new event participant
func (s *participantService) Create(ctx context.Context, req model.EventParticipantCreateRequest) (uuid.UUID, error) {
	// Verify that user exists
	_, err := s.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		s.logger.Errorw("Failed to get user for participant creation", "error", err, "userId", req.UserID)
		return uuid.Nil, fmt.Errorf("invalid user: %w", err)
	}

	// Check if participant already exists
	exists, err := s.repo.Exists(ctx, req.EventID, req.UserID)
	if err != nil {
		s.logger.Errorw("Failed to check if participant exists", "error", err, "eventId", req.EventID, "userId", req.UserID)
		return uuid.Nil, fmt.Errorf("failed to check participant: %w", err)
	}

	if exists {
		return uuid.Nil, errors.New("user is already a participant of this event")
	}

	participant := model.EventParticipant{
		EventID:     req.EventID,
		UserID:      req.UserID,
		IsConfirmed: req.IsConfirmed,
	}

	id, err := s.repo.Create(ctx, participant)
	if err != nil {
		s.logger.Errorw("Failed to create participant", "error", err, "eventId", req.EventID, "userId", req.UserID)
		return uuid.Nil, fmt.Errorf("failed to create participant: %w", err)
	}

	return id, nil
}

// GetByID retrieves an event participant by ID
func (s *participantService) GetByID(ctx context.Context, id uuid.UUID) (model.EventParticipantResponse, error) {
	participant, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Errorw("Failed to get participant by ID", "error", err, "id", id)
		return model.EventParticipantResponse{}, fmt.Errorf("failed to get participant: %w", err)
	}

	// Get user info
	user, err := s.userRepo.GetByID(ctx, participant.UserID)
	if err != nil {
		s.logger.Warnw("Failed to get participant user details", "error", err, "userId", participant.UserID)
		// Continue even if we can't get user details
	}

	return model.EventParticipantResponse{
		ID:          participant.ID,
		EventID:     participant.EventID,
		UserID:      participant.UserID,
		Username:    user.Username,
		FullName:    fmt.Sprintf("%s %s", user.FirstName, user.LastName),
		IsConfirmed: participant.IsConfirmed,
		JoinedAt:    participant.JoinedAt,
	}, nil
}

// Update updates an event participant
func (s *participantService) Update(ctx context.Context, id uuid.UUID, req model.EventParticipantUpdateRequest) error {
	participant, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Errorw("Failed to get participant for update", "error", err, "id", id)
		return fmt.Errorf("failed to get participant: %w", err)
	}

	// Update fields if provided
	if req.IsConfirmed != nil {
		participant.IsConfirmed = *req.IsConfirmed
	}

	if err := s.repo.Update(ctx, participant); err != nil {
		s.logger.Errorw("Failed to update participant", "error", err, "id", id)
		return fmt.Errorf("failed to update participant: %w", err)
	}

	return nil
}

// Delete deletes an event participant
func (s *participantService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Errorw("Failed to delete participant", "error", err, "id", id)
		return fmt.Errorf("failed to delete participant: %w", err)
	}

	return nil
}

// ListByEvent retrieves a list of participants for a specific event with pagination
func (s *participantService) ListByEvent(ctx context.Context, eventID uuid.UUID, page, size int) (model.EventParticipantsResponse, error) {
	// Set default pagination values if not provided
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}

	offset := (page - 1) * size
	participants, total, err := s.repo.ListByEvent(ctx, eventID, size, offset)
	if err != nil {
		s.logger.Errorw("Failed to list event participants", "error", err, "eventId", eventID, "page", page, "size", size)
		return model.EventParticipantsResponse{}, fmt.Errorf("failed to list participants: %w", err)
	}

	// Convert to response objects
	participantResponses := make([]model.EventParticipantResponse, len(participants))
	for i, participant := range participants {
		// Get user info
		user, err := s.userRepo.GetByID(ctx, participant.UserID)
		if err != nil {
			s.logger.Warnw("Failed to get participant user details", "error", err, "userId", participant.UserID)
			// Continue with minimal user info if we can't get full details
			participantResponses[i] = model.EventParticipantResponse{
				ID:          participant.ID,
				EventID:     participant.EventID,
				UserID:      participant.UserID,
				IsConfirmed: participant.IsConfirmed,
				JoinedAt:    participant.JoinedAt,
			}
			continue
		}

		participantResponses[i] = model.EventParticipantResponse{
			ID:          participant.ID,
			EventID:     participant.EventID,
			UserID:      participant.UserID,
			Username:    user.Username,
			FullName:    fmt.Sprintf("%s %s", user.FirstName, user.LastName),
			IsConfirmed: participant.IsConfirmed,
			JoinedAt:    participant.JoinedAt,
		}
	}

	return model.EventParticipantsResponse{
		Participants: participantResponses,
		Total:        total,
	}, nil
}

// ListByUser retrieves a list of event participations for a specific user with pagination
func (s *participantService) ListByUser(ctx context.Context, userID uuid.UUID, page, size int) (model.EventParticipantsResponse, error) {
	// Set default pagination values if not provided
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}

	offset := (page - 1) * size
	participants, total, err := s.repo.ListByUser(ctx, userID, size, offset)
	if err != nil {
		s.logger.Errorw("Failed to list user participations", "error", err, "userId", userID, "page", page, "size", size)
		return model.EventParticipantsResponse{}, fmt.Errorf("failed to list participations: %w", err)
	}

	// Get user info
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		s.logger.Warnw("Failed to get user details for participations", "error", err, "userId", userID)
		// Continue with minimal info
	}

	// Convert to response objects
	participantResponses := make([]model.EventParticipantResponse, len(participants))
	for i, participant := range participants {
		participantResponses[i] = model.EventParticipantResponse{
			ID:          participant.ID,
			EventID:     participant.EventID,
			UserID:      participant.UserID,
			Username:    user.Username,
			FullName:    fmt.Sprintf("%s %s", user.FirstName, user.LastName),
			IsConfirmed: participant.IsConfirmed,
			JoinedAt:    participant.JoinedAt,
		}
	}

	return model.EventParticipantsResponse{
		Participants: participantResponses,
		Total:        total,
	}, nil
}

// ConfirmParticipation confirms a user's participation in an event
func (s *participantService) ConfirmParticipation(ctx context.Context, id uuid.UUID) error {
	participant, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Errorw("Failed to get participant for confirmation", "error", err, "id", id)
		return fmt.Errorf("failed to get participant: %w", err)
	}

	participant.IsConfirmed = true

	if err := s.repo.Update(ctx, participant); err != nil {
		s.logger.Errorw("Failed to confirm participation", "error", err, "id", id)
		return fmt.Errorf("failed to confirm participation: %w", err)
	}

	return nil
}

// DeclineParticipation declines a user's participation in an event
func (s *participantService) DeclineParticipation(ctx context.Context, id uuid.UUID) error {
	participant, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Errorw("Failed to get participant for declining", "error", err, "id", id)
		return fmt.Errorf("failed to get participant: %w", err)
	}

	participant.IsConfirmed = false

	if err := s.repo.Update(ctx, participant); err != nil {
		s.logger.Errorw("Failed to decline participation", "error", err, "id", id)
		return fmt.Errorf("failed to decline participation: %w", err)
	}

	return nil
}
