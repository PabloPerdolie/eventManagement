package event

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/PabloPerdolie/event-manager/core-service/internal/domain"
	"github.com/pkg/errors"
	"time"

	"github.com/PabloPerdolie/event-manager/core-service/internal/model"
	"go.uber.org/zap"
)

type ParticipantRepo interface {
	Create(ctx context.Context, participant model.EventParticipant) (int, error)
	GetById(ctx context.Context, id int) (model.EventParticipant, error)
	GetByEventAndUser(ctx context.Context, eventId, userId int) (model.EventParticipant, error)
	Update(ctx context.Context, participant model.EventParticipant) error
	Delete(ctx context.Context, id int) error
	//DeleteByEventAndUser(ctx context.Context, eventId, userId int) error // no usages
	ListByEvent(ctx context.Context, eventId, limit, offset int) ([]model.EventParticipant, error)
	ListByUser(ctx context.Context, userId, limit, offset int) ([]model.EventParticipant, error)
}

type UserRepo interface {
	GetUserById(ctx context.Context, id int) (*model.User, error)
	ListUsers(ctx context.Context, limit, offset int) ([]model.User, int, error)
}

type Participant struct {
	repo      ParticipantRepo
	userRepo  UserRepo
	notifyPbl NotifyPublisher
	logger    *zap.SugaredLogger
}

func NewParticipantService(repo ParticipantRepo, userRepo UserRepo, notifyPbl NotifyPublisher, logger *zap.SugaredLogger) Participant {
	return Participant{
		repo:      repo,
		userRepo:  userRepo,
		notifyPbl: notifyPbl,
		logger:    logger,
	}
}

func (s Participant) Create(ctx context.Context, eventID int, req domain.EventParticipantCreateRequest) (*domain.EventParticipantResponse, error) {
	user, err := s.userRepo.GetUserById(ctx, req.UserID)
	if err != nil {
		s.logger.Errorw("Failed to get user for participant creation", "error", err, "userID", req.UserID)
		return nil, errors.WithMessage(err, "invalid user")
	}

	_, err = s.repo.GetByEventAndUser(ctx, eventID, req.UserID)
	if err == nil {
		return nil, model.ErrUserAlreadyAnParticipant
	}
	if !errors.Is(err, sql.ErrNoRows) {
		s.logger.Errorw("Failed to check if participant exists", "error", err, "eventID", eventID, "userID", req.UserID)
		return nil, errors.WithMessage(err, "check participant")
	}

	now := time.Now()
	participant := model.EventParticipant{
		EventID:     eventID,
		UserID:      req.UserID,
		Role:        model.RoleParticipant,
		JoinedAt:    &now,
		IsConfirmed: ptr(true), // Default to confirmed // todo
	}

	id, err := s.repo.Create(ctx, participant)
	if err != nil {
		s.logger.Errorw("Failed to create participant", "error", err, "eventID", eventID, "userID", req.UserID)
		return nil, errors.WithMessage(err, "create participant")
	}

	data := map[string]any{
		"event": "participant_added",
		"data": map[string]any{
			"event_name": req.EventTitle,
			"user_email": user.Email,
		},
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, errors.WithMessage(err, "marshal participant notify")
	}

	err = s.notifyPbl.Publish(ctx, bytes)
	if err != nil {
		return nil, errors.WithMessage(err, "publish participant notify")
	}

	return &domain.EventParticipantResponse{
		Id:          id,
		EventID:     eventID,
		User:        userToResponse(user),
		Role:        string(participant.Role),
		JoinedAt:    participant.JoinedAt,
		IsConfirmed: participant.IsConfirmed,
	}, nil
}

func (s Participant) GetById(ctx context.Context, id int) (*domain.EventParticipantResponse, error) {
	participant, err := s.repo.GetById(ctx, id)
	if err != nil {
		s.logger.Errorw("Failed to get participant by ID", "error", err, "id", id)
		return nil, errors.WithMessage(err, "get participant")
	}

	user, err := s.userRepo.GetUserById(ctx, participant.UserID)
	if err != nil {
		s.logger.Warnw("Failed to get participant user details", "error", err, "userID", participant.UserID)
		user = &model.User{UserId: participant.UserID}
	}

	return &domain.EventParticipantResponse{
		Id:          participant.EventParticipantID,
		EventID:     participant.EventID,
		User:        userToResponse(user),
		Role:        string(participant.Role),
		JoinedAt:    participant.JoinedAt,
		IsConfirmed: participant.IsConfirmed,
	}, nil
}

func (s Participant) Delete(ctx context.Context, id int) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Errorw("Failed to delete participant", "error", err, "id", id)
		return errors.WithMessage(err, "failed to delete participant")
	}

	return nil
}

func (s Participant) ListByEvent(ctx context.Context, eventID int, page, size int) (*domain.EventParticipantsResponse, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}
	offset := (page - 1) * size

	participants, err := s.repo.ListByEvent(ctx, eventID, size, offset)
	if err != nil {
		s.logger.Errorw("Failed to list event participants", "error", err, "eventID", eventID, "page", page, "size", size)
		return nil, errors.WithMessage(err, "list participants")
	}

	participantResponses := make([]domain.EventParticipantResponse, len(participants))
	for i, participant := range participants {
		user, err := s.userRepo.GetUserById(ctx, participant.UserID)
		if err != nil {
			s.logger.Warnw("Failed to get participant user details", "error", err, "userID", participant.UserID)
			user = &model.User{UserId: participant.UserID}
		}
		participantResponses[i] = domain.EventParticipantResponse{
			Id:          participant.EventParticipantID,
			EventID:     participant.EventID,
			User:        userToResponse(user),
			Role:        string(participant.Role),
			JoinedAt:    participant.JoinedAt,
			IsConfirmed: participant.IsConfirmed,
		}
	}

	return &domain.EventParticipantsResponse{
		Participants: participantResponses,
		Total:        len(participants),
	}, nil
}

func (s Participant) ListByUser(ctx context.Context, userId int, page, size int) (*domain.EventParticipantsResponse, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}
	offset := (page - 1) * size

	participants, err := s.repo.ListByUser(ctx, userId, size, offset)
	if err != nil {
		s.logger.Errorw("Failed to list user participations", "error", err, "userId", userId, "page", page, "size", size)
		return nil, errors.WithMessage(err, "list participations")
	}

	participantResponses := make([]domain.EventParticipantResponse, len(participants))
	for i, participant := range participants {
		participantResponses[i] = domain.EventParticipantResponse{
			Id:          participant.EventParticipantID,
			EventID:     participant.EventID,
			Role:        string(participant.Role),
			JoinedAt:    participant.JoinedAt,
			IsConfirmed: participant.IsConfirmed,
		}
	}

	return &domain.EventParticipantsResponse{
		Participants: participantResponses,
		Total:        len(participants),
	}, nil
}

func (s Participant) ConfirmParticipation(ctx context.Context, id int) error {
	participant, err := s.repo.GetById(ctx, id)
	if err != nil {
		s.logger.Errorw("Failed to get participant for confirmation", "error", err, "id", id)
		return errors.WithMessage(err, "get participant")
	}

	participant.IsConfirmed = ptr(true)

	if err := s.repo.Update(ctx, participant); err != nil {
		s.logger.Errorw("Failed to confirm participation", "error", err, "id", id)
		return errors.WithMessage(err, "confirm participation")
	}

	return nil
}

func (s Participant) DeclineParticipation(ctx context.Context, id int) error {
	participant, err := s.repo.GetById(ctx, id)
	if err != nil {
		s.logger.Errorw("Failed to get participant for declining", "error", err, "id", id)
		return errors.WithMessage(err, "failed to get participant")
	}

	participant.IsConfirmed = ptr(false)

	if err := s.repo.Update(ctx, participant); err != nil {
		s.logger.Errorw("Failed to decline participation", "error", err, "id", id)
		return errors.WithMessage(err, "failed to decline participation")
	}

	return nil
}

func userToResponse(user *model.User) domain.UserResponse {
	return domain.UserResponse{
		Id:        user.UserId,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		IsActive:  user.IsActive,
		Role:      user.Role,
	}
}

func ptr(b bool) *bool {
	return &b
}
