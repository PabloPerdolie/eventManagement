package handler

import (
	"context"
	"encoding/json"
	"github.com/PabloPerdolie/event-manager/communication-service/internal/domain"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type Comment interface {
	CreateComment(ctx context.Context, comment domain.CreateCommentMessage) (int, error)
}

type CommentConsumer struct {
	commentService Comment
	logger         *zap.SugaredLogger
}

func NewCommentHandler(commentService Comment, logger *zap.SugaredLogger) CommentConsumer {
	return CommentConsumer{
		commentService: commentService,
		logger:         logger,
	}
}

func (h CommentConsumer) ProcessMessage(msg amqp.Delivery) {
	var message domain.CreateCommentMessage
	err := json.Unmarshal(msg.Body, &message)
	if err != nil {
		h.logger.Errorw("failed to unmarshal message", "error", err)
		msg.Reject(false)
		return
	}

	ctx := context.Background()
	id, err := h.commentService.CreateComment(ctx, message)
	if err != nil {
		h.logger.Errorw("failed to create comment", "error", err)
		msg.Reject(false)
		return
	}

	h.logger.Infow("comment created", "comment_id", id)
	msg.Ack(false)
}
