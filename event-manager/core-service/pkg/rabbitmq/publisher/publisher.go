package rabbitmq

import (
	"context"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQPublisher struct {
	queueName  string
	connection *amqp.Connection
	channel    *amqp.Channel
}

func New(queueName, rabbitUrl string) (*RabbitMQPublisher, error) {
	connection, err := amqp.Dial(rabbitUrl)
	if err != nil {
		return nil, errors.WithMessage(err, "connect to RabbitMQ")
	}

	channel, err := connection.Channel()
	if err != nil {
		connection.Close()
		return nil, errors.WithMessage(err, "open a channel")
	}

	return &RabbitMQPublisher{
		connection: connection,
		channel:    channel,
		queueName:  queueName,
	}, nil
}

func (c *RabbitMQPublisher) Publish(ctx context.Context, body []byte) error {
	return c.channel.PublishWithContext(
		ctx,
		"",
		c.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)
}

func (c *RabbitMQPublisher) Stop() {
	c.cleanup()
}

func (c *RabbitMQPublisher) cleanup() {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.connection != nil {
		c.connection.Close()
	}
}
