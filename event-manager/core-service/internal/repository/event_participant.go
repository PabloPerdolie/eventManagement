package repository

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"

	"github.com/PabloPerdolie/event-manager/core-service/internal/model"
	"github.com/jmoiron/sqlx"
)

type Participant struct {
	db *sqlx.DB
}

func NewParticipant(db *sqlx.DB) Participant {
	return Participant{
		db: db,
	}
}

func (r Participant) Create(ctx context.Context, participant model.EventParticipant) (int, error) {
	query := `
		INSERT INTO event_participant (event_id, user_id, role, joined_at, is_confirmed)
		VALUES ($1, $2, $3, $4, $5) RETURNING event_participant_id
	`

	var id int
	err := r.db.QueryRowContext(
		ctx,
		query,
		participant.EventID,
		participant.UserID,
		participant.Role,
		participant.JoinedAt,
		participant.IsConfirmed,
	).Scan(&id)
	if err != nil {
		return 0, errors.WithMessage(err, "create event participant")
	}

	return id, nil
}

func (r Participant) GetById(ctx context.Context, id int) (model.EventParticipant, error) {
	var participant model.EventParticipant
	query := `SELECT * FROM event_participant WHERE event_participant_id = $1`
	err := r.db.GetContext(ctx, &participant, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.EventParticipant{}, errors.WithMessage(err, "event participant not found")
		}
		return model.EventParticipant{}, errors.WithMessage(err, "get event participant")
	}
	return participant, nil
}

func (r Participant) GetByEventAndUser(ctx context.Context, eventId, userId int) (model.EventParticipant, error) {
	var participant model.EventParticipant
	query := `SELECT * FROM event_participant WHERE event_id = $1 AND user_id = $2`
	err := r.db.GetContext(ctx, &participant, query, eventId, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.EventParticipant{}, errors.WithMessage(err, "event participant not found")
		}
		return model.EventParticipant{}, errors.WithMessage(err, "get event participant")
	}
	return participant, nil
}

func (r Participant) Update(ctx context.Context, participant model.EventParticipant) error {
	query := `UPDATE event_participant SET role = $1, is_confirmed = $2 WHERE event_participant_id = $3`
	_, err := r.db.ExecContext(ctx, query, participant.Role, participant.IsConfirmed, participant.EventParticipantID)
	if err != nil {
		return errors.WithMessage(err, "update event participant")
	}
	return nil
}

func (r Participant) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM event_participant WHERE event_participant_id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return errors.WithMessage(err, "delete event participant")
	}
	return nil
}

func (r Participant) DeleteByEventAndUser(ctx context.Context, eventId, userId int) error {
	query := `DELETE FROM event_participant WHERE event_id = $1 AND user_id = $2`
	_, err := r.db.ExecContext(ctx, query, eventId, userId)
	if err != nil {
		return errors.WithMessage(err, "delete event participant")
	}
	return nil
}

func (r Participant) ListByEvent(ctx context.Context, eventId, limit, offset int) ([]model.EventParticipant, error) {
	var participants []model.EventParticipant

	query := `SELECT * FROM event_participant WHERE event_id = $1 ORDER BY joined_at DESC LIMIT $2 OFFSET $3`
	if err := r.db.SelectContext(ctx, &participants, query, eventId, limit, offset); err != nil {
		return nil, errors.WithMessage(err, "list event participants")
	}

	return participants, nil
}

func (r Participant) ListByUser(ctx context.Context, userId, limit, offset int) ([]model.EventParticipant, error) {
	var participants []model.EventParticipant

	query := `SELECT * FROM event_participant WHERE user_id = $1 ORDER BY joined_at DESC LIMIT $2 OFFSET $3`
	if err := r.db.SelectContext(ctx, &participants, query, userId, limit, offset); err != nil {
		return nil, errors.WithMessage(err, "list user events")
	}

	return participants, nil
}
