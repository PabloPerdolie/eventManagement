package repository

import (
	"context"
	"github.com/PabloPerdolie/event-manager/communication-service/internal/model"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Comment struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) Comment {
	return Comment{
		db: db,
	}
}

func (r Comment) Insert(ctx context.Context, comment model.Comment) (int, error) {
	query := `
		INSERT INTO comments(event_id, sender_id, content, task_id) 
		VALUES ($1, $2, $3, $4) RETURNING comment_id`

	row := r.db.QueryRowContext(ctx, query, comment.EventId, comment.SenderId, comment.Content, comment.TaskId)
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

func (r Comment) GetByEventId(ctx context.Context, eventId int) ([]model.Comment, error) {
	query := `
		SELECT comment_id, event_id, sender_id, task_id, content, created_at, is_deleted, is_read
		FROM comments 
		WHERE event_id = $1 AND is_deleted = false
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, eventId)
	if err != nil {
		return nil, errors.WithMessage(err, "get comments by event id")
	}
	defer rows.Close()

	var comments []model.Comment
	for rows.Next() {
		var comment model.Comment
		err := rows.Scan(
			&comment.CommentId,
			&comment.EventId,
			&comment.SenderId,
			&comment.TaskId,
			&comment.Content,
			&comment.CreatedAt,
			&comment.IsDeleted,
			&comment.IsRead,
		)
		if err != nil {
			return nil, errors.WithMessage(err, "scan comment")
		}
		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.WithMessage(err, "iterate comments")
	}

	return comments, nil
}

func (r Comment) Delete(ctx context.Context, commentId int) error {
	query := `
		UPDATE comments 
		SET is_deleted = true 
		WHERE comment_id = $1`

	_, err := r.db.ExecContext(ctx, query, commentId)
	if err != nil {
		return errors.WithMessage(err, "delete comment")
	}

	return nil
}

func (r Comment) MarkAsRead(ctx context.Context, commentId int) error {
	query := `
		UPDATE comments 
		SET is_read = true 
		WHERE comment_id = $1`

	_, err := r.db.ExecContext(ctx, query, commentId)
	if err != nil {
		return errors.WithMessage(err, "mark comment as read")
	}

	return nil
}
