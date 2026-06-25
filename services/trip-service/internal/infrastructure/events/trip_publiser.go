package events

import (
	"context"
	"ride-sharing/shared/message"
)

type TripEventPublish struct {
	rabbitmq *message.RabbitMQ
}

func NewTripEventPublisher(rabbitmq *message.RabbitMQ) *TripEventPublish {
	return &TripEventPublish{
		rabbitmq: rabbitmq,
	}
}

func (p *TripEventPublish) PublishMessage(ctx context.Context) error {
	return p.rabbitmq.PublishMessage(ctx, "hello", "hello world")
}
