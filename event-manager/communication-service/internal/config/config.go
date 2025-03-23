package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port     string
	Postgres Postgres
	RabbitMQ RabbitMQ
}

type Postgres struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type RabbitMQ struct {
	Host     string
	Port     string
	Username string
	Password string
	Queue    string
}

func New() (*Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
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

	return &Config{
		Port: port,
		RabbitMQ: RabbitMQ{
			Host:     rabbitHost,
			Port:     rabbitPort,
			Username: rabbitUsername,
			Password: rabbitPassword,
			Queue:    rabbitQueue,
		},
	}, nil
}

func (c *RabbitMQ) GetRabbitMQURL() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s/", c.Username, c.Password, c.Host, c.Port)
}

func (p Postgres) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		p.Host, p.Port, p.User, p.Password, p.DBName, p.SSLMode)
}
