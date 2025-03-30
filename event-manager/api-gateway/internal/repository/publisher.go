package repository

import (
	"context"

	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	rabbitMQConn *amqp.Connection
	rabbitMQChan *amqp.Channel
	queueName    string
}

func New(rabbitMQConn *amqp.Connection, rabbitMQChan *amqp.Channel, queueName string) Publisher {
	return Publisher{
		rabbitMQConn: rabbitMQConn,
		rabbitMQChan: rabbitMQChan,
		queueName:    queueName,
	}
}

func (p *Publisher) Publish(ctx context.Context, data []byte) error {
	err := p.rabbitMQChan.PublishWithContext(
		ctx,
		"",
		p.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         data,
			DeliveryMode: amqp.Persistent,
		},
	)
	if err != nil {
		return errors.WithMessage(err, "publish data")
	}

	return nil
}
