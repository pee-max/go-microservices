package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ride-sharing/shared/env"
	"ride-sharing/shared/messaging"
)

var (
	httpAddr = env.GetString("HTTP_ADDR", ":8081")
)

func main() {
	uri := env.GetString("RABBITMQ_URI", "amqp://guest:guest@rabbitmq:5672/")
	log.Println("Starting API Gateway")

	rabbitmq, err := messaging.NewRabbitMQ(uri)
	if err != nil {
		log.Fatal(err)
	}
	defer rabbitmq.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("POST /trip/preview", enableCros(handleTripPreview))
	mux.HandleFunc("POST /trip/start", enableCros(handleTripStart))
	mux.HandleFunc("/ws/drivers", func(w http.ResponseWriter, r *http.Request) { handleDriversWebSocket(w, r, rabbitmq) })
	mux.HandleFunc("/ws/riders", func(w http.ResponseWriter, r *http.Request) { handleRidersWebSocket(w, r, rabbitmq) })

	server := &http.Server{
		Addr:    httpAddr,
		Handler: mux,
	}

	serverError := make(chan error, 1)

	go func() {
		log.Printf("server starting on %s", httpAddr)
		serverError <- server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)

	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverError:
		log.Printf("Error starting the server: %v", err)

	case sig := <-shutdown:
		log.Printf("Server is shutting down due to %v", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Could not stop the server gracefully: %v", err)
			server.Close()
		}
	}

}
