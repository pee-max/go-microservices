package main

import (
	"context"
	"log"
	"ride-sharing/shared/message"

	amqp "github.com/rabbitmq/amqp091-go"
)

type tripConsumer struct {
	rabbitmq *message.RabbitMQ
}

func NewTripConsumer(rabbitmq *message.RabbitMQ) *tripConsumer {
	return &tripConsumer{
		rabbitmq: rabbitmq,
	}
}

func (c *tripConsumer) Listen() error {
	return c.rabbitmq.ConsumeMessage("hello", func(ctx context.Context, msg amqp.Delivery) error {
		log.Printf("driver received message: %v", msg)
		return nil
	})
}
