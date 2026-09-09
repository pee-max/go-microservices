package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"ride-sharing/services/trip-service/internal/domain"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/messaging"

	pbd "ride-sharing/shared/proto/driver"

	amqp "github.com/rabbitmq/amqp091-go"
)

type driverConsumer struct {
	rabbitmq *messaging.RabbitMQ
	service  domain.TripService
}

func NewDriverConsumer(rabbitmq *messaging.RabbitMQ, service domain.TripService) *driverConsumer {
	return &driverConsumer{
		rabbitmq: rabbitmq,
		service:  service,
	}
}

func (c *driverConsumer) Listen() error {
	return c.rabbitmq.ConsumeMessage(messaging.DriverTripResponseQueue, func(ctx context.Context, msg amqp.Delivery) error {
		var message contracts.AmqpMessage
		if err := json.Unmarshal(msg.Body, &message); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			return err
		}
		var payload messaging.DriverTripResponseData
		if err := json.Unmarshal(message.Data, &payload); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			return err
		}
		log.Printf("driver response received message: %+v", payload)

		switch msg.RoutingKey {
		case contracts.DriverCmdTripAccept:
			if err := c.handleTripAccept(ctx, payload.TripID, payload.Driver); err != nil {
				log.Printf("Failed to handle the trip accept: %v", err)
				return err
			}
		case contracts.DriverCmdTripDecline:
			if err := c.handleTripDeclined(ctx, payload.TripID, payload.RiderID); err != nil {
				log.Printf("Failed to handle the trip decline: %v", err)
				return err
			}
		}
		log.Printf("Unknow trip event: %v", payload)

		return nil
	})
}

func (c *driverConsumer) handleTripAccept(ctx context.Context, tripID string, driver *pbd.Driver) error {
	trip, err := c.service.GetTripByID(ctx, tripID)
	if err != nil {
		return err
	}
	if trip == nil {
		return fmt.Errorf("Trip was not found: %s", tripID)
	}

	if err := c.service.UpdateTrip(ctx, tripID, "accept", driver); err != nil {
		log.Printf("Failed to update trip: %v", err)
		return err
	}

	trip, err = c.service.GetTripByID(ctx, tripID)
	if err != nil {
		return err
	}

	marshaledTrip, err := json.Marshal(trip)
	if err != nil {
		return err
	}

	if err := c.rabbitmq.PublishMessage(ctx, contracts.TripEventDriverAssigned, contracts.AmqpMessage{
		OwnerID: trip.UserID,
		Data:    marshaledTrip,
	}); err != nil {
		return err
	}

	//todo: payment

	return nil
}

func (c *driverConsumer) handleTripDeclined(ctx context.Context, tripID string, riderID string) error {
	trip, err := c.service.GetTripByID(ctx, tripID)
	if err != nil {
		return err
	}
	if trip == nil {
		log.Printf("Trip not found with the ID: %s", tripID)
		return nil
	}
	newPayload := messaging.TripEventData{
		Trip: trip.Toproto(),
	}

	marshaledPayload, err := json.Marshal(newPayload)
	if err != nil {
		return err
	}

	if err := c.rabbitmq.PublishMessage(ctx, contracts.TripEventDriverNotInterested, contracts.AmqpMessage{
		OwnerID: riderID,
		Data:    marshaledPayload,
	}); err != nil {
		return err
	}

	return nil
}
