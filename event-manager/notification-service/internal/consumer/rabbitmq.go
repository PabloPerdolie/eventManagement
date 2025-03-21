package consumer

import (
	"fmt"

	"github.com/PabloPerdolie/event-manager/notification-service/internal/config"
	"github.com/PabloPerdolie/event-manager/notification-service/internal/service"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type RabbitMQConsumer struct {
	service    *service.Service
	config     *config.Config
	logger     *zap.SugaredLogger
	connection *amqp.Connection
	channel    *amqp.Channel
}

func New(svc *service.Service, cfg *config.Config, logger *zap.SugaredLogger) *RabbitMQConsumer {
	return &RabbitMQConsumer{
		service: svc,
		config:  cfg,
		logger:  logger,
	}
}

func (c *RabbitMQConsumer) Start() error {
	var err error
	c.connection, err = amqp.Dial(c.config.RabbitMQ.GetRabbitMQURL())
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	c.channel, err = c.connection.Channel()
	if err != nil {
		c.connection.Close()
		return fmt.Errorf("failed to open a channel: %w", err)
	}

	queue, err := c.channel.QueueDeclare(
		c.config.RabbitMQ.Queue, // name
		true,                    // durable
		false,                   // delete when unused
		false,                   // exclusive
		false,                   // no-wait
		nil,                     // arguments
	)
	if err != nil {
		c.cleanup()
		return fmt.Errorf("failed to declare a queue: %w", err)
	}

	err = c.channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		c.cleanup()
		return fmt.Errorf("failed to set QoS: %w", err)
	}

	msgs, err := c.channel.Consume(
		queue.Name, // queue
		"",         // consumer
		false,      // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		c.cleanup()
		return fmt.Errorf("failed to register a consumer: %w", err)
	}

	c.logger.Infof("Connected to RabbitMQ and consuming from queue: %s", c.config.RabbitMQ.Queue)

	for msg := range msgs {
		c.processMessage(msg)
	}

	return nil
}

func (c *RabbitMQConsumer) Stop() {
	c.cleanup()
}

func (c *RabbitMQConsumer) cleanup() {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.connection != nil {
		c.connection.Close()
	}
}
