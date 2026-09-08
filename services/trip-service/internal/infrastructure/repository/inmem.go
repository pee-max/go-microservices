package repository

import (
	"context"
	"fmt"
	"ride-sharing/services/trip-service/internal/domain"
	pbd "ride-sharing/shared/proto/driver"
	pb "ride-sharing/shared/proto/trip"
)

type inmemRopository struct {
	trips     map[string]*domain.TripModel
	rideFares map[string]*domain.RideFareModel
}

func NewInmemRepository() *inmemRopository {
	return &inmemRopository{
		trips:     make(map[string]*domain.TripModel),
		rideFares: make(map[string]*domain.RideFareModel),
	}
}

func (r *inmemRopository) CreatTrip(ctx context.Context, trip *domain.TripModel) (*domain.TripModel, error) {
	r.trips[trip.ID.Hex()] = trip
	return trip, nil
}

func (r *inmemRopository) SaveRideFare(ctx context.Context, f *domain.RideFareModel) error {
	r.rideFares[f.ID.Hex()] = f
	return nil
}

func (r *inmemRopository) GetRideFareByID(ctx context.Context, id string) (*domain.RideFareModel, error) {
	rideFare, ok := r.rideFares[id]
	if ok {
		return rideFare, nil
	}
	return nil, fmt.Errorf("fare does not exist with the id: %v", id)
}

func (r *inmemRopository) GetTripByID(ctx context.Context, id string) (*domain.TripModel, error) {
	trip, ok := r.trips[id]
	if ok {
		return trip, nil
	}
	return nil, nil
}

func (r *inmemRopository) UpdateTrip(ctx context.Context, id string, status string, diver *pbd.Driver) error {
	trip, ok := r.trips[id]
	if !ok {
		return fmt.Errorf("Trip not found whit the ID: %s", id)
	}

	trip.Status = status
	if diver != nil {
		trip.Driver = &pb.TripDriver{
			Id:             diver.Id,
			Name:           diver.Name,
			ProfilePicture: diver.ProfilePicture,
			CarPlate:       diver.CarPlate,
		}
	}
	r.trips[id] = trip
	return nil
}
