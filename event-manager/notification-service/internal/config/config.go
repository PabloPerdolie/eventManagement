package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port     string
	RabbitMQ RabbitMQConfig
	SMTP     SMTPConfig
}

type RabbitMQConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Queue    string
}

type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Sender   string
}

func New() (*Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	rabbitHost := os.Getenv("RABBITMQ_HOST")
	if rabbitHost == "" {
		rabbitHost = "localhost"
	}

	rabbitPort := os.Getenv("RABBITMQ_PORT")
	if rabbitPort == "" {
		rabbitPort = "5672"
	}

	rabbitUsername := os.Getenv("RABBITMQ_USERNAME")
	if rabbitUsername == "" {
		rabbitUsername = "guest"
	}

	rabbitPassword := os.Getenv("RABBITMQ_PASSWORD")
	if rabbitPassword == "" {
		rabbitPassword = "guest"
	}

	rabbitQueue := os.Getenv("RABBITMQ_QUEUE")
	if rabbitQueue == "" {
		rabbitQueue = "notifications"
	}

	// SMTP configuration
	smtpHost := os.Getenv("SMTP_HOST")
	if smtpHost == "" {
		smtpHost = "smtp.example.com"
	}

	smtpPort := os.Getenv("SMTP_PORT")
	if smtpPort == "" {
		smtpPort = "587"
	}

	smtpUsername := os.Getenv("SMTP_USERNAME")
	if smtpUsername == "" {
		smtpUsername = "notifications@system.com"
	}

	smtpPassword := os.Getenv("SMTP_PASSWORD")
	if smtpPassword == "" {
		smtpPassword = "password"
	}

	smtpSender := os.Getenv("SMTP_SENDER")
	if smtpSender == "" {
		smtpSender = "notifications@system.com"
	}

	return &Config{
		Port: port,
		RabbitMQ: RabbitMQConfig{
			Host:     rabbitHost,
			Port:     rabbitPort,
			Username: rabbitUsername,
			Password: rabbitPassword,
			Queue:    rabbitQueue,
		},
		SMTP: SMTPConfig{
			Host:     smtpHost,
			Port:     smtpPort,
			Username: smtpUsername,
			Password: smtpPassword,
			Sender:   smtpSender,
		},
	}, nil
}

func (c *RabbitMQConfig) GetRabbitMQURL() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s/", c.Username, c.Password, c.Host, c.Port)
}

func (c *SMTPConfig) GetSMTPAddress() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}
