package event

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/PabloPerdolie/event-manager/core-service/internal/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// ParticipantRepository defines event participant repository interface
type ParticipantRepository interface {
	Create(ctx context.Context, participant model.EventParticipant) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.EventParticipant, error)
	GetByEventAndUser(ctx context.Context, eventID, userID uuid.UUID) (model.EventParticipant, error)
	Update(ctx context.Context, participant model.EventParticipant) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByEventAndUser(ctx context.Context, eventID, userID uuid.UUID) error
	ListByEvent(ctx context.Context, eventID uuid.UUID, limit, offset int) ([]model.EventParticipant, int, error)
	ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.EventParticipant, int, error)
}

type participantRepository struct {
	db *sqlx.DB
}

// NewParticipantRepository creates a new event participant repository
func NewParticipantRepository(db *sqlx.DB) ParticipantRepository {
	return &participantRepository{db: db}
}

// Create creates a new event participant in the database
func (r *participantRepository) Create(ctx context.Context, participant model.EventParticipant) (uuid.UUID, error) {
	participant.ID = uuid.New()
	participant.JoinedAt = time.Now()

	query := `
		INSERT INTO event_participants (id, event_id, user_id, joined_at, is_confirmed)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		participant.ID,
		participant.EventID,
		participant.UserID,
		participant.JoinedAt,
		participant.IsConfirmed,
	)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create event participant: %w", err)
	}

	return participant.ID, nil
}

// GetByID retrieves an event participant by ID
func (r *participantRepository) GetByID(ctx context.Context, id uuid.UUID) (model.EventParticipant, error) {
	var participant model.EventParticipant

	query := `
		SELECT id, event_id, user_id, joined_at, is_confirmed
		FROM event_participants
		WHERE id = $1
	`

	err := r.db.GetContext(ctx, &participant, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.EventParticipant{}, fmt.Errorf("event participant not found: %w", err)
		}
		return model.EventParticipant{}, fmt.Errorf("failed to get event participant: %w", err)
	}

	return participant, nil
}

// GetByEventAndUser retrieves an event participant by event ID and user ID
func (r *participantRepository) GetByEventAndUser(ctx context.Context, eventID, userID uuid.UUID) (model.EventParticipant, error) {
	var participant model.EventParticipant

	query := `
		SELECT id, event_id, user_id, joined_at, is_confirmed
		FROM event_participants
		WHERE event_id = $1 AND user_id = $2
	`

	err := r.db.GetContext(ctx, &participant, query, eventID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.EventParticipant{}, fmt.Errorf("event participant not found: %w", err)
		}
		return model.EventParticipant{}, fmt.Errorf("failed to get event participant: %w", err)
	}

	return participant, nil
}

// Update updates an event participant in the database
func (r *participantRepository) Update(ctx context.Context, participant model.EventParticipant) error {
	query := `
		UPDATE event_participants
		SET is_confirmed = $1
		WHERE id = $2
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		participant.IsConfirmed,
		participant.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update event participant: %w", err)
	}

	return nil
}

// Delete deletes an event participant from the database
func (r *participantRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM event_participants WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete event participant: %w", err)
	}

	return nil
}

// DeleteByEventAndUser deletes an event participant by event ID and user ID
func (r *participantRepository) DeleteByEventAndUser(ctx context.Context, eventID, userID uuid.UUID) error {
	query := `DELETE FROM event_participants WHERE event_id = $1 AND user_id = $2`

	_, err := r.db.ExecContext(ctx, query, eventID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete event participant: %w", err)
	}

	return nil
}

// ListByEvent retrieves a list of event participants for a specific event with pagination
func (r *participantRepository) ListByEvent(ctx context.Context, eventID uuid.UUID, limit, offset int) ([]model.EventParticipant, int, error) {
	var participants []model.EventParticipant
	var total int

	// Count total participants for the event
	countQuery := `SELECT COUNT(*) FROM event_participants WHERE event_id = $1`
	err := r.db.GetContext(ctx, &total, countQuery, eventID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count event participants: %w", err)
	}

	// Retrieve participants with pagination
	query := `
		SELECT id, event_id, user_id, joined_at, is_confirmed
		FROM event_participants
		WHERE event_id = $1
		ORDER BY joined_at DESC
		LIMIT $2 OFFSET $3
	`

	err = r.db.SelectContext(ctx, &participants, query, eventID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list event participants: %w", err)
	}

	return participants, total, nil
}

// ListByUser retrieves a list of event participants for a specific user with pagination
func (r *participantRepository) ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.EventParticipant, int, error) {
	var participants []model.EventParticipant
	var total int

	// Count total events for the user
	countQuery := `SELECT COUNT(*) FROM event_participants WHERE user_id = $1`
	err := r.db.GetContext(ctx, &total, countQuery, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count user events: %w", err)
	}

	// Retrieve events with pagination
	query := `
		SELECT id, event_id, user_id, joined_at, is_confirmed
		FROM event_participants
		WHERE user_id = $1
		ORDER BY joined_at DESC
		LIMIT $2 OFFSET $3
	`

	err = r.db.SelectContext(ctx, &participants, query, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list user events: %w", err)
	}

	return participants, total, nil
}
