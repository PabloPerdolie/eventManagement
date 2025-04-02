package config

import (
	"github.com/pkg/errors"
	"os"
)

type Config struct {
	Port             string
	DatabaseURL      string
	RabbitMQURL      string
	CommentQueueName string
}

func New() (*Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, errors.New("DATABASE_URL environment variable is required")
	}

	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	if rabbitMQURL == "" {
		rabbitMQURL = "amqp://guest:guest@localhost:5672/"
	}

	commentQueueName := os.Getenv("COMMENT_QUEUE_NAME")
	if commentQueueName == "" {
		commentQueueName = "comments"
	}

	return &Config{
		Port:             port,
		DatabaseURL:      dbURL,
		RabbitMQURL:      rabbitMQURL,
		CommentQueueName: commentQueueName,
	}, nil
}
