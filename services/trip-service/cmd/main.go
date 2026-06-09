package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	h "ride-sharing/services/trip-service/internal/infrastructure/http"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"
	"syscall"
	"time"
)

func main() {

	inmemRepo := repository.NewInmemRepository()
	svc := service.NewTripService(inmemRepo)
	mux := http.NewServeMux()

	httphandler := h.HttpHandler{Service: svc}

	mux.HandleFunc("POST /preview", httphandler.HandleTripPreview)

	server := &http.Server{
		Addr:    ":8083",
		Handler: mux,
	}

	serverError := make(chan error, 1)

	go func() {
		log.Printf("Trip server starting on: %v", server.Addr)
		serverError <- server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)

	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverError:
		log.Printf("Error starting trip servier: %v", err)

	case sig := <-shutdown:
		log.Printf("Trip service shutdown due to: %v", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Trip server could not shutdown greatfully: %v", err)
			server.Close()
		}

	}

}
