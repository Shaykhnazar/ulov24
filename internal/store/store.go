package store

import (
	"github.com/Shaykhnazar/ulov24/internal/ridematching"
	"github.com/Shaykhnazar/ulov24/internal/rides"
)

// Kichik interfeyslar: kodning har bir qismi o'ziga kerak bo'lganinigina so'raydi.

// RouteStore yo'nalishlarni saqlaydi. SaveRoute ID 0 bo'lsa yangi raqam beradi,
// aks holda shu raqamli yo'nalishni almashtiradi.
type RouteStore interface {
	SaveRoute(rides.DriverRoute) (rides.DriverRoute, error)
	Routes() ([]rides.DriverRoute, error)
}

// RequestStore yo'lovchi so'rovlarini saqlaydi.
type RequestStore interface {
	SaveRequest(rides.PassengerRequest) (rides.PassengerRequest, error)
	Requests() ([]rides.PassengerRequest, error)
}

// MatchStore mosliklarni saqlaydi.
type MatchStore interface {
	SaveMatch(ridematching.RideMatch) (ridematching.RideMatch, error)
	Matches() ([]ridematching.RideMatch, error)
}

// Store — uchalasi birga. Interfeyslarni ham embedding bilan qo'shish mumkin.
type Store interface {
	RouteStore
	RequestStore
	MatchStore
}
