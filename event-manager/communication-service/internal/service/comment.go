package service

import (
	"context"
	"encoding/json"
	"github.com/PabloPerdolie/event-manager/communication-service/internal/model"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type CommentRepo interface {
	Insert(ctx context.Context, comment model.Comment) (int, error)
	GetById(ctx context.Context, commentId int) (model.Comment, error)
	GetByEventId(ctx context.Context, eventId int) ([]model.Comment, error)
	Delete(ctx context.Context, commentId int) error
	MarkAsRead(ctx context.Context, commentId int) error
}

type CreateCommentMessage struct {
	EventId  int    `json:"event_id"`
	SenderId int    `json:"sender_id"`
	Content  string `json:"content"`
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

func (s Comment) CreateComment(ctx context.Context, comment model.Comment) (int, error) {
	id, err := s.commentRepo.Insert(ctx, comment)
	if err != nil {
		return 0, errors.WithMessage(err, "insert comment")
	}
	return id, nil
}

func (s Comment) GetCommentById(ctx context.Context, id int) (model.Comment, error) {
	comment, err := s.commentRepo.GetById(ctx, id)
	if err != nil {
		return model.Comment{}, errors.WithMessage(err, "get comment by id")
	}
	return comment, nil
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

func (s Comment) ProcessCreateCommentMessage(msg amqp.Delivery) {
	var message CreateCommentMessage
	err := json.Unmarshal(msg.Body, &message)
	if err != nil {
		s.logger.Errorw("failed to unmarshal message", "error", err)
		msg.Reject(false)
		return
	}

	comment := model.Comment{
		EventId:  message.EventId,
		SenderId: message.SenderId,
		Content:  message.Content,
	}

	ctx := context.Background()
	id, err := s.CreateComment(ctx, comment)
	if err != nil {
		s.logger.Errorw("failed to create comment", "error", err)
		msg.Reject(false)
		return
	}

	s.logger.Infow("comment created", "comment_id", id)
	msg.Ack(false)
}
