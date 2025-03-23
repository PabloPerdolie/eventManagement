package comment

import (
	"context"
	"encoding/json"
	"github.com/event-management/api-gateway/internal/config"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"time"
)

// CreateCommentMessage is the message structure for creating a comment
type CreateCommentMessage struct {
	EventId  int    `json:"event_id" binding:"required"`
	SenderId int    `json:"sender_id" binding:"required"`
	Content  string `json:"content" binding:"required"`
}

// Service handles comments operations
type Service struct {
	rabbitMQConn *amqp.Connection
	rabbitMQChan *amqp.Channel
	queueName    string
	logger       *zap.SugaredLogger
}

// New creates a new comment service
func New(cfg *config.Config, logger *zap.SugaredLogger) (*Service, error) {
	conn, err := amqp.Dial(cfg.RabbitMQURL)
	if err != nil {
		return nil, errors.Wrap(err, "connect to RabbitMQ")
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, errors.Wrap(err, "open a channel")
	}

	// Ensure the queue exists
	_, err = ch.QueueDeclare(
		cfg.CommentQueueName, // name
		true,                // durable
		false,               // delete when unused
		false,               // exclusive
		false,               // no-wait
		nil,                 // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, errors.Wrap(err, "declare a queue")
	}

	return &Service{
		rabbitMQConn: conn,
		rabbitMQChan: ch,
		queueName:    cfg.CommentQueueName,
		logger:       logger,
	}, nil
}

// Close releases resources used by the service
func (s *Service) Close() {
	if s.rabbitMQChan != nil {
		s.rabbitMQChan.Close()
	}
	if s.rabbitMQConn != nil {
		s.rabbitMQConn.Close()
	}
}

// CreateComment publishes a message to create a comment
func (s *Service) CreateComment(ctx context.Context, message CreateCommentMessage) error {
	body, err := json.Marshal(message)
	if err != nil {
		return errors.Wrap(err, "marshal message")
	}

	// Create a context with timeout for the publishing
	publishCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err = s.rabbitMQChan.PublishWithContext(
		publishCtx,
		"",           // exchange
		s.queueName,  // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // Make message persistent
		},
	)
	if err != nil {
		return errors.Wrap(err, "publish message")
	}

	s.logger.Infow("Comment creation message published",
		"event_id", message.EventId,
		"sender_id", message.SenderId,
	)

	return nil
}
