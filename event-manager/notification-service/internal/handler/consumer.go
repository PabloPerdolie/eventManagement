package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/PabloPerdolie/event-manager/notification-service/internal/domain"
	"github.com/PabloPerdolie/event-manager/notification-service/internal/model"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type NotificationService interface {
	ProcessNotification(ctx context.Context, message domain.NotificationMessage) error
}

type NotifyHandler struct {
	logger  *zap.SugaredLogger
	service NotificationService
}

func NewNotifyHandler(logger *zap.SugaredLogger, service NotificationService) NotifyHandler {
	return NotifyHandler{
		logger:  logger,
		service: service,
	}
}

func (c NotifyHandler) ProcessMessage(msg amqp.Delivery) {
	c.logger.Infof("Received a message: %s", msg.Body)

	ctx := context.Background()

	var notification domain.NotificationMessage
	if err := json.Unmarshal(msg.Body, &notification); err != nil {
		c.logger.Errorf("Failed to unmarshal message: %v", err)
		msg.Reject(false)
		return
	}

	err := c.service.ProcessNotification(ctx, notification)
	switch {
	case errors.Is(err, model.ErrUnsupportedEventType):
		msg.Reject(false)
	case errors.Is(err, model.ErrInvalidNotificationData):
		msg.Reject(false)
	case err != nil:
		msg.Reject(true)
	default:
		msg.Ack(false)
	}
}
