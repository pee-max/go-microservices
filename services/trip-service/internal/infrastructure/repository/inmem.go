package repository

import (
	"context"
	"fmt"
	"ride-sharing/services/trip-service/internal/domain"
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
