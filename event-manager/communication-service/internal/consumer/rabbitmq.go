package consumer

import (
	"github.com/PabloPerdolie/event-manager/communication-service/internal/config"
	"github.com/pkg/errors"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type Handler interface {
	ProcessMessage(msg amqp.Delivery)
}

type RabbitMQConsumer struct {
	handler    Handler
	config     *config.Config
	logger     *zap.SugaredLogger
	connection *amqp.Connection
	channel    *amqp.Channel
}

func New(handler Handler, config *config.Config, logger *zap.SugaredLogger) *RabbitMQConsumer {
	return &RabbitMQConsumer{
		handler: handler,
		config:  config,
		logger:  logger,
	}
}

func (c *RabbitMQConsumer) Start() error {
	var err error
	c.connection, err = amqp.Dial(c.config.RabbitMQ.GetRabbitMQURL())
	if err != nil {
		return errors.WithMessage(err, "connect to RabbitMQ")
	}

	c.channel, err = c.connection.Channel()
	if err != nil {
		c.connection.Close()
		return errors.WithMessage(err, "open a channel")
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
		return errors.WithMessage(err, "declare a queue")
	}

	err = c.channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		c.cleanup()
		return errors.WithMessage(err, "set QoS")
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
		return errors.WithMessage(err, "register a consumer")
	}

	c.logger.Infof("Connected to RabbitMQ and consuming from queue: %s", c.config.RabbitMQ.Queue)

	for msg := range msgs {
		c.handler.ProcessMessage(msg)
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
