package handler

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type Comment interface {
	ProcessCreateCommentMessage(msg amqp.Delivery)
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
	h.logger.Infow("received message", "routing_key", msg.RoutingKey)
	h.commentService.ProcessCreateCommentMessage(msg)
}
