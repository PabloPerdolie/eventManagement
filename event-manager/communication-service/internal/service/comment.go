package service

import (
	"context"
	"github.com/PabloPerdolie/event-manager/communication-service/internal/domain"
	"github.com/PabloPerdolie/event-manager/communication-service/internal/model"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type CommentRepo interface {
	Insert(ctx context.Context, comment model.Comment) (int, error)
	GetByEventId(ctx context.Context, eventId int) ([]model.Comment, error)
	Delete(ctx context.Context, commentId int) error
	MarkAsRead(ctx context.Context, commentId int) error
}

type Comment struct {
	commentRepo CommentRepo
	logger      *zap.SugaredLogger
}

func NewComment(commentRepo CommentRepo, logger *zap.SugaredLogger) Comment {
	return Comment{
		commentRepo: commentRepo,
		logger:      logger,
	}
}

func (s Comment) CreateComment(ctx context.Context, comment domain.CreateCommentMessage) (int, error) {
	commentModel := model.Comment{
		EventId:  comment.EventId,
		SenderId: comment.SenderId,
		Content:  comment.Content,
		TaskId:   comment.TaskId,
	}

	id, err := s.commentRepo.Insert(ctx, commentModel)
	if err != nil {
		return 0, errors.WithMessage(err, "insert comment")
	}
	return id, nil
}

func (s Comment) GetCommentsByEventId(ctx context.Context, eventId int) ([]model.Comment, error) {
	comments, err := s.commentRepo.GetByEventId(ctx, eventId)
	if err != nil {
		return nil, errors.WithMessage(err, "get comments by event id")
	}
	return comments, nil
}

func (s Comment) DeleteComment(ctx context.Context, id int) error {
	err := s.commentRepo.Delete(ctx, id)
	if err != nil {
		return errors.WithMessage(err, "delete comment")
	}
	return nil
}

func (s Comment) MarkCommentAsRead(ctx context.Context, id int) error {
	err := s.commentRepo.MarkAsRead(ctx, id)
	if err != nil {
		return errors.WithMessage(err, "mark comment as read")
	}
	return nil
}
