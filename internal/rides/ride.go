package rides

import (
	"errors"
	"time"
)

// ErrRideCapacityFull — yo'nalishda so'ralgancha bo'sh joy yo'q.
// ridematching paketi ham aynan shu qiymatni ishlatadi.
var ErrRideCapacityFull = errors.New("ride route does not have enough remaining seats")

// DriverRouteStatus — haydovchi yo'nalishining holati.
type DriverRouteStatus string

const (
	DriverRouteOpen      DriverRouteStatus = "open"
	DriverRouteFull      DriverRouteStatus = "full"
	DriverRouteCancelled DriverRouteStatus = "cancelled"
	DriverRouteCompleted DriverRouteStatus = "completed"
	DriverRouteExpired   DriverRouteStatus = "expired"
)

// PassengerRequestStatus — yo'lovchi so'rovining holati.
type PassengerRequestStatus string

const (
	PassengerRequestOpen      PassengerRequestStatus = "open"
	PassengerRequestMatched   PassengerRequestStatus = "matched"
	PassengerRequestAccepted  PassengerRequestStatus = "accepted"
	PassengerRequestCancelled PassengerRequestStatus = "cancelled"
	PassengerRequestCompleted PassengerRequestStatus = "completed"
	PassengerRequestExpired   PassengerRequestStatus = "expired"
)

// DriverRoute — haydovchi e'lon qilgan safar. JSON nomlari manapochta'dagi
// Ulov24 bilan bir xil.
type DriverRoute struct {
	ID             int64             `json:"id"`
	DriverUserID   int64             `json:"driver_user_id"`
	FromLocationID int64             `json:"from_location_id"`
	ToLocationID   int64             `json:"to_location_id"`
	DepartureAt    time.Time         `json:"departure_at"`
	SeatsTotal     int               `json:"seats_total"`
	SeatsRemaining int               `json:"seats_remaining"`
	PricePerSeat   *int              `json:"price_per_seat,omitempty"` // nil — narx kelishiladi
	TripMode       TripMode          `json:"trip_mode"`
	Status         DriverRouteStatus `json:"status"`
}

// PassengerRequest — yo'lovchining so'rovi.
type PassengerRequest struct {
	ID                   int64                  `json:"id"`
	PassengerUserID      int64                  `json:"passenger_user_id"`
	FromLocationID       int64                  `json:"from_location_id"`
	ToLocationID         int64                  `json:"to_location_id"`
	DesiredDepartureFrom time.Time              `json:"desired_departure_from"`
	DesiredDepartureTo   time.Time              `json:"desired_departure_to"`
	PassengerCount       int                    `json:"passenger_count"`
	TripMode             TripMode               `json:"trip_mode"`
	Status               PassengerRequestStatus `json:"status"`
}

// IsExpired: jo'nash vaqti o'tib ketgan.
func (r DriverRoute) IsExpired(now time.Time) bool {
	return now.After(r.DepartureAt)
}

// Reserve yo'nalishdan count ta joyni band qiladi. Joy yetmasa, yo'nalish
// o'zgarmaydi va ErrRideCapacityFull qaytadi.
func (r *DriverRoute) Reserve(count int) error {
	if count > r.SeatsRemaining {
		return ErrRideCapacityFull
	}
	r.SeatsRemaining -= count
	if r.SeatsRemaining == 0 {
		r.Status = DriverRouteFull
	}
	return nil
}
