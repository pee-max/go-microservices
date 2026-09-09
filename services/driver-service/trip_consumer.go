package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand/v2"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/messaging"

	amqp "github.com/rabbitmq/amqp091-go"
)

type tripConsumer struct {
	rabbitmq *messaging.RabbitMQ
	service  *Service
}

func NewTripConsumer(rabbitmq *messaging.RabbitMQ, service *Service) *tripConsumer {
	return &tripConsumer{
		rabbitmq: rabbitmq,
		service:  service,
	}
}

func (c *tripConsumer) Listen() error {
	return c.rabbitmq.ConsumeMessage(messaging.FindAvailableDriversQueue, func(ctx context.Context, msg amqp.Delivery) error {
		var tripEvent contracts.AmqpMessage
		if err := json.Unmarshal(msg.Body, &tripEvent); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			return err
		}
		var payload messaging.TripEventData
		if err := json.Unmarshal(tripEvent.Data, &payload); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			return err
		}
		log.Printf("driver received message: %+v", payload)

		switch msg.RoutingKey {
		case contracts.TripEventCreated, contracts.TripEventDriverNotInterested:
			c.handleFindAndNotifyDrivers(ctx, payload)
		}

		return nil
	})
}

func (c *tripConsumer) handleFindAndNotifyDrivers(ctx context.Context, payload messaging.TripEventData) error {
	availableDriversID := c.service.FindAvailableDriversID(payload.Trip.GetSelectedFare().PackageSlug)

	if len(availableDriversID) == 0 {
		if err := c.rabbitmq.PublishMessage(ctx, contracts.TripEventNoDriversFound, contracts.AmqpMessage{
			OwnerID: payload.Trip.UserID,
		}); err != nil {
			log.Printf("Failed to publish message to exchange: %v", err)
			return err
		}
		return nil
	}

	randomIndex := rand.IntN(len(availableDriversID))

	suitableDriver := availableDriversID[randomIndex]
	marshaledEvet, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	if err := c.rabbitmq.PublishMessage(ctx, contracts.DriverCmdTripRequest, contracts.AmqpMessage{
		OwnerID: suitableDriver,
		Data:    marshaledEvet,
	}); err != nil {
		log.Printf("Failed to publish message to exchange: %v", err)
		return err
	}

	return nil
}
