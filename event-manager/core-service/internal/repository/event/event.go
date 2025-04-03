package event

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"time"

	"github.com/PabloPerdolie/event-manager/core-service/internal/model"
	"github.com/jmoiron/sqlx"
)

type Event struct {
	db *sqlx.DB
}

func NewEvent(db *sqlx.DB) Event {
	return Event{
		db: db,
	}
}

func (r *Event) Create(ctx context.Context, event model.Event) (int, error) {
	event.CreatedAt = time.Now()

	query := `
		INSERT INTO events (organizer_id, title, description, start_date, end_date, location, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING event_id
	`

	var eventID int
	err := r.db.QueryRowContext(
		ctx,
		query,
		event.OrganizerID,
		event.Title,
		event.Description,
		event.StartDate,
		event.EndDate,
		event.Location,
		event.Status,
		event.CreatedAt,
	).Scan(&eventID)
	if err != nil {
		return 0, errors.WithMessage(err, "create event")
	}

	return eventID, nil
}

func (r *Event) GetById(ctx context.Context, eventID int) (model.Event, error) {
	var event model.Event

	query := `
		SELECT event_id, organizer_id, title, description, start_date, end_date, location, status, created_at
		FROM events
		WHERE event_id = $1
	`

	err := r.db.GetContext(ctx, &event, query, eventID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Event{}, errors.WithMessage(err, "event not found")
		}
		return model.Event{}, errors.WithMessage(err, "get event")
	}

	return event, nil
}

func (r *Event) Update(ctx context.Context, event model.Event) error {
	query := `
		UPDATE events
		SET title = $1, description = $2, start_date = $3, end_date = $4, location = $5, status = $6
		WHERE event_id = $7
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		event.Title,
		event.Description,
		event.StartDate,
		event.EndDate,
		event.Location,
		event.Status,
		event.EventID,
	)
	if err != nil {
		return errors.WithMessage(err, "update event")
	}

	return nil
}

func (r *Event) Delete(ctx context.Context, eventID int) error {
	query := `DELETE FROM events WHERE event_id = $1`

	_, err := r.db.ExecContext(ctx, query, eventID)
	if err != nil {
		return errors.WithMessage(err, "delete event")
	}

	return nil
}

func (r *Event) List(ctx context.Context, limit, offset int) ([]model.Event, int, error) {
	var events []model.Event
	var total int

	countQuery := `SELECT COUNT(*) FROM events`
	err := r.db.GetContext(ctx, &total, countQuery)
	if err != nil {
		return nil, 0, errors.WithMessage(err, "count events")
	}

	query := `
		SELECT event_id, organizer_id, title, description, start_date, end_date, location, status, created_at
		FROM events
		ORDER BY start_date DESC
		LIMIT $1 OFFSET $2
	`

	err = r.db.SelectContext(ctx, &events, query, limit, offset)
	if err != nil {
		return nil, 0, errors.WithMessage(err, "list events")
	}

	return events, total, nil
}

func (r *Event) ListByOrganizer(ctx context.Context, organizerID, limit, offset int) ([]model.Event, int, error) {
	var events []model.Event
	var total int

	countQuery := `SELECT COUNT(*) FROM events WHERE organizer_id = $1`
	err := r.db.GetContext(ctx, &total, countQuery, organizerID)
	if err != nil {
		return nil, 0, errors.WithMessage(err, "count organizer events")
	}

	query := `
		SELECT event_id, organizer_id, title, description, start_date, end_date, location, status, created_at
		FROM events
		WHERE organizer_id = $1
		ORDER BY start_date DESC
		LIMIT $2 OFFSET $3
	`

	err = r.db.SelectContext(ctx, &events, query, organizerID, limit, offset)
	if err != nil {
		return nil, 0, errors.WithMessage(err, "list organizer events")
	}

	return events, total, nil
}

func (r *Event) ListByParticipant(ctx context.Context, participantID, limit, offset int) ([]model.Event, int, error) {
	var events []model.Event
	var total int

	countQuery := `
		SELECT COUNT(*)
		FROM events e
		JOIN event_participants ep ON e.event_id = ep.event_id
		WHERE ep.user_id = $1
	`
	err := r.db.GetContext(ctx, &total, countQuery, participantID)
	if err != nil {
		return nil, 0, errors.WithMessage(err, "count participant events")
	}

	query := `
		SELECT e.event_id, e.organizer_id, e.title, e.description, e.start_date, e.end_date, e.location, e.status, e.created_at
		FROM events e
		JOIN event_participants ep ON e.event_id = ep.event_id
		WHERE ep.user_id = $1
		ORDER BY e.start_date DESC
		LIMIT $2 OFFSET $3
	`

	err = r.db.SelectContext(ctx, &events, query, participantID, limit, offset)
	if err != nil {
		return nil, 0, errors.WithMessage(err, "list participant events")
	}

	return events, total, nil
}
