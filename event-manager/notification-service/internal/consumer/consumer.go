package consumer

import (
	"encoding/json"
	"github.com/PabloPerdolie/event-manager/notification-service/internal/model"
	"github.com/streadway/amqp"
)

func (c *RabbitMQConsumer) processMessage(msg amqp.Delivery) {
	c.logger.Infof("Received a message: %s", msg.Body)

	var notification model.NotificationMessage
	if err := json.Unmarshal(msg.Body, &notification); err != nil {
		c.logger.Errorf("Failed to unmarshal message: %v", err)
		msg.Reject(false)
		return
	}

	if err := c.service.Notification.ProcessNotification(&notification); err != nil {
		c.logger.Errorf("Failed to process notification: %v", err)
		if err == model.ErrUnsupportedEventType || err == model.ErrInvalidNotificationData {
			msg.Reject(false)
		} else {
			msg.Reject(true)
		}
		return
	}

	msg.Ack(false)
	c.logger.Info("Message processed successfully")
}
