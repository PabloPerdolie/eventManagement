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

// Repository defines event repository interface
type Repository interface {
	Create(ctx context.Context, event model.Event) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.Event, error)
	Update(ctx context.Context, event model.Event) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]model.Event, int, error)
	ListByOrganizer(ctx context.Context, organizerID uuid.UUID, limit, offset int) ([]model.Event, int, error)
	ListByParticipant(ctx context.Context, participantID uuid.UUID, limit, offset int) ([]model.Event, int, error)
}

type repository struct {
	db *sqlx.DB
}

// NewRepository creates a new event repository
func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

// Create creates a new event in the database
func (r *repository) Create(ctx context.Context, event model.Event) (uuid.UUID, error) {
	event.ID = uuid.New()
	event.CreatedAt = time.Now()
	event.UpdatedAt = time.Now()

	query := `
		INSERT INTO events (id, name, description, start_date, end_date, location, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		event.ID,
		event.Name,
		event.Description,
		event.StartDate,
		event.EndDate,
		event.Location,
		event.CreatedBy,
		event.CreatedAt,
		event.UpdatedAt,
	)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create event: %w", err)
	}

	return event.ID, nil
}

// GetByID retrieves an event by ID
func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (model.Event, error) {
	var event model.Event

	query := `
		SELECT id, name, description, start_date, end_date, location, created_by, created_at, updated_at
		FROM events
		WHERE id = $1
	`

	err := r.db.GetContext(ctx, &event, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Event{}, fmt.Errorf("event not found: %w", err)
		}
		return model.Event{}, fmt.Errorf("failed to get event: %w", err)
	}

	return event, nil
}

// Update updates an event in the database
func (r *repository) Update(ctx context.Context, event model.Event) error {
	event.UpdatedAt = time.Now()

	query := `
		UPDATE events
		SET name = $1, description = $2, start_date = $3, end_date = $4, location = $5, updated_at = $6
		WHERE id = $7
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		event.Name,
		event.Description,
		event.StartDate,
		event.EndDate,
		event.Location,
		event.UpdatedAt,
		event.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update event: %w", err)
	}

	return nil
}

// Delete deletes an event from the database
func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM events WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}

	return nil
}

// List retrieves a list of events with pagination
func (r *repository) List(ctx context.Context, limit, offset int) ([]model.Event, int, error) {
	var events []model.Event
	var total int

	// Count total events
	countQuery := `SELECT COUNT(*) FROM events`
	err := r.db.GetContext(ctx, &total, countQuery)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count events: %w", err)
	}

	// Retrieve events with pagination
	query := `
		SELECT id, name, description, start_date, end_date, location, created_by, created_at, updated_at
		FROM events
		ORDER BY start_date DESC
		LIMIT $1 OFFSET $2
	`

	err = r.db.SelectContext(ctx, &events, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list events: %w", err)
	}

	return events, total, nil
}

// ListByOrganizer retrieves a list of events for a specific organizer with pagination
func (r *repository) ListByOrganizer(ctx context.Context, organizerID uuid.UUID, limit, offset int) ([]model.Event, int, error) {
	var events []model.Event
	var total int

	// Count total events by organizer
	countQuery := `SELECT COUNT(*) FROM events WHERE created_by = $1`
	err := r.db.GetContext(ctx, &total, countQuery, organizerID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count organizer events: %w", err)
	}

	// Retrieve events with pagination
	query := `
		SELECT id, name, description, start_date, end_date, location, created_by, created_at, updated_at
		FROM events
		WHERE created_by = $1
		ORDER BY start_date DESC
		LIMIT $2 OFFSET $3
	`

	err = r.db.SelectContext(ctx, &events, query, organizerID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list organizer events: %w", err)
	}

	return events, total, nil
}

// ListByParticipant retrieves a list of events for a specific participant with pagination
func (r *repository) ListByParticipant(ctx context.Context, participantID uuid.UUID, limit, offset int) ([]model.Event, int, error) {
	var events []model.Event
	var total int

	// Count total events by participant
	countQuery := `
		SELECT COUNT(*) FROM events e
		JOIN event_participants ep ON e.id = ep.event_id
		WHERE ep.user_id = $1
	`
	err := r.db.GetContext(ctx, &total, countQuery, participantID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count participant events: %w", err)
	}

	// Retrieve events with pagination
	query := `
		SELECT e.id, e.name, e.description, e.start_date, e.end_date, e.location, e.created_by, e.created_at, e.updated_at
		FROM events e
		JOIN event_participants ep ON e.id = ep.event_id
		WHERE ep.user_id = $1
		ORDER BY e.start_date DESC
		LIMIT $2 OFFSET $3
	`

	err = r.db.SelectContext(ctx, &events, query, participantID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list participant events: %w", err)
	}

	return events, total, nil
}
