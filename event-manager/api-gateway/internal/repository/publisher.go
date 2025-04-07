package repository

import (
	"context"
	rabbitmq "github.com/PabloPerdolie/event-manager/api-gateway/pkg/rabbitmq/publisher"

	"github.com/pkg/errors"
)

type Publisher struct {
	pbl *rabbitmq.RabbitMQPublisher
}

func NewPublisher(pbl *rabbitmq.RabbitMQPublisher) Publisher {
	return Publisher{
		pbl: pbl,
	}
}

func (p Publisher) Publish(ctx context.Context, data []byte) error {
	err := p.pbl.Publish(ctx, data)
	if err != nil {
		return errors.WithMessage(err, "publish data")
	}

	return nil
}
