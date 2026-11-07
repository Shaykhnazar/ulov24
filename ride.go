package main

import "time"

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

// DriverRoute — haydovchi e'lon qilgan safar: qayerdan, qayerga, qachon,
// nechta bo'sh joy bilan. struct 6-modulda batafsil o'rganiladi.
type DriverRoute struct {
	ID             int64
	DriverUserID   int64
	FromLocationID int64
	ToLocationID   int64
	DepartureAt    time.Time
	SeatsTotal     int
	SeatsRemaining int
	PricePerSeat   int
	TripMode       TripMode
	Status         DriverRouteStatus
}

// PassengerRequest — yo'lovchining so'rovi: qayerdan, qayerga, qaysi vaqt
// oralig'ida, necha kishi.
type PassengerRequest struct {
	ID                   int64
	PassengerUserID      int64
	FromLocationID       int64
	ToLocationID         int64
	DesiredDepartureFrom time.Time
	DesiredDepartureTo   time.Time
	PassengerCount       int
	TripMode             TripMode
	Status               PassengerRequestStatus
}
