package messaging

import pb "ride-sharing/shared/proto/trip"

const (
	FindAvailableDriversQueue = "find_available_drivers"
	DriverCmdTripRequestQueue = "trip_cmd_request"
)

type TripEventData struct {
	Trip *pb.Trip `json:"trip"`
}
