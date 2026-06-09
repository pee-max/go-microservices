package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"ride-sharing/services/trip-service/internal/domain"
	"ride-sharing/shared/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TripService struct {
	repo domain.TripRepository
}

func NewTripService(repo domain.TripRepository) *TripService {
	return &TripService{
		repo: repo,
	}
}

func (s *TripService) CreatTrip(ctx context.Context, fare *domain.RideFareModel) (*domain.TripModel, error) {
	t := &domain.TripModel{
		ID:       primitive.NewObjectID(),
		UserID:   fare.UserID,
		Status:   "pending",
		RideFare: fare,
	}
	return s.repo.CreatTrip(ctx, t)
}

func (s *TripService) GetRoute(ctx context.Context, pickup, destination *types.Coordinate) (*types.OsrmApiResponse, error) {
	url := fmt.Sprintf("http://router.project-osrm.org/nearest/v1/driving/%f,%f;%f,%f?overview=full&geometries-geojson",
		pickup.Longitude, pickup.Latitude,
		destination.Longitude, destination.Latitude,
	)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch route from OSRM API: %v", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("faild to read the response: %v", err)
	}

	var routeResp types.OsrmApiResponse
	if err := json.Unmarshal(body, &routeResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	return &routeResp, nil
}
