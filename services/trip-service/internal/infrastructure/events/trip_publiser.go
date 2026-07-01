package events

import (
	"context"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/messaging"
)

type TripEventPublish struct {
	rabbitmq *messaging.RabbitMQ
}

func NewTripEventPublisher(rabbitmq *messaging.RabbitMQ) *TripEventPublish {
	return &TripEventPublish{
		rabbitmq: rabbitmq,
	}
}

func (p *TripEventPublish) PublishMessage(ctx context.Context) error {
	return p.rabbitmq.PublishMessage(ctx, contracts.TripEventCreated, "Trip has been created")
}
