package repository

import (
	"context"
	"database/sql"
	"github.com/PabloPerdolie/event-manager/communication-service/internal/model"
	"github.com/pkg/errors"
)

type Comment struct {
	db *sql.DB
}

func New(db *sql.DB) Comment {
	return Comment{
		db: db,
	}
}

func (r Comment) Insert(ctx context.Context, comment model.Comment) (int, error) {
	query := `
		INSERT INTO comments(event_id, sender_id, content) 
		VALUES ($1, $2, $3) RETURNING comment_id`

	row := r.db.QueryRowContext(ctx, query, comment.EventId, comment.SenderId, comment.Content)
	if row.Err() != nil {
		return 0, errors.WithMessage(row.Err(), "insert comment")
	}

	var id int
	err := row.Scan(&id)
	if err != nil {
		return 0, errors.WithMessage(err, "scan value")
	}

	return id, nil
}
