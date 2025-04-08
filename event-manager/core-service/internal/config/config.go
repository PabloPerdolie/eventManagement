package config

import (
	"github.com/pkg/errors"
	"os"
)

type Config struct {
	Port                  string
	DatabaseURL           string
	RabbitMQURL           string
	CommentServiceUrl     string
	NotificationQueueName string
}

func New() (*Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, errors.New("DATABASE_URL environment variable is required")
	}

	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	if rabbitMQURL == "" {
		rabbitMQURL = "amqp://guest:guest@localhost:5672/"
	}

	notificationsQueueName := os.Getenv("NOTIFICATION_QUEUE_NAME")
	if notificationsQueueName == "" {
		notificationsQueueName = "notifications"
	}

	commentServiceUrl := os.Getenv("COMMUNICATION_SERVICE_URL")
	if commentServiceUrl == "" {
		commentServiceUrl = "http://communication-service:8083"
	}

	return &Config{
		Port:                  port,
		DatabaseURL:           dbURL,
		RabbitMQURL:           rabbitMQURL,
		CommentServiceUrl:     commentServiceUrl,
		NotificationQueueName: notificationsQueueName,
	}, nil
}
