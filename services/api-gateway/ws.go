package main

import (
	"log"
	"net/http"
	"ride-sharing/services/api-gateway/grpc_clients"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/messaging"
	"ride-sharing/shared/proto/driver"
)

var (
	connManager = messaging.NewConnectionManager()
)

func handleDriversWebSocket(w http.ResponseWriter, r *http.Request, rb *messaging.RabbitMQ) {
	conn, err := connManager.Upgrade(w, r)
	if err != nil {

	}
	defer conn.Close()

	userId := r.URL.Query().Get("userID")
	if userId == "" {
		log.Println("No user id provided")
		return
	}

	packageSlug := r.URL.Query().Get("packageSlug")
	if packageSlug == "" {
		log.Println("No package Slug provided")
		return
	}

	connManager.Add(userId, conn)

	ctx := r.Context()

	driverServiceClient, err := grpc_clients.NewDriverServiceClient()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		connManager.Remove(userId)

		driverServiceClient.Client.UnregisterDriver(ctx, &driver.RegisterDriverRequest{
			DriverID:    userId,
			PackageSlug: packageSlug,
		})

		driverServiceClient.Close()

		log.Println("Driver unregistered: ", userId)
	}()

	driverData, err := driverServiceClient.Client.RegisterDriver(ctx, &driver.RegisterDriverRequest{
		DriverID:    userId,
		PackageSlug: packageSlug,
	})
	if err != nil {
		log.Printf("Error registering driver: %v", err)
		return
	}

	if err := connManager.SendMessage(userId, contracts.WSMessage{
		Type: contracts.DriverCmdRegister,
		Data: driverData.Driver,
	}); err != nil {
		log.Printf("Error sending message: %v", err)
		return
	}

	//Initialize queue consumers
	queues := []string{
		messaging.DriverCmdTripRequestQueue,
	}

	for _, q := range queues {
		consumer := messaging.NewQueueConsumer(rb, connManager, q)
		if err := consumer.Start(); err != nil {
			log.Printf("failed to start consumer for queue: %s err: %v", q, err)
		}
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		log.Printf("received message: %s", message)
	}
}

func handleRidersWebSocket(w http.ResponseWriter, r *http.Request, rb *messaging.RabbitMQ) {
	conn, err := connManager.Upgrade(w, r)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	defer conn.Close()

	userId := r.URL.Query().Get("userID")
	if userId == "" {
		log.Println("No user id provided")
		return
	}

	connManager.Add(userId, conn)
	defer connManager.Remove(userId)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		log.Printf("received message: %s", message)
	}

}
