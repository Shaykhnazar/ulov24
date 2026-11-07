package store

import (
	"fmt"
	"io"

	"github.com/Shaykhnazar/ulov24/internal/ridematching"
	"github.com/Shaykhnazar/ulov24/internal/rides"
)

// LoggingStore istalgan Store'ni o'raydi va har bir saqlashni Log'ga yozadi.
// O'qish metodlari ichki Store'dan embedding orqali o'zgarishsiz keladi.
type LoggingStore struct {
	Store
	Log io.Writer
}

func (s LoggingStore) SaveRoute(r rides.DriverRoute) (rides.DriverRoute, error) {
	r, err := s.Store.SaveRoute(r)
	if err == nil {
		fmt.Fprintf(s.Log, "store: SaveRoute #%d\n", r.ID)
	}
	return r, err
}

func (s LoggingStore) SaveRequest(q rides.PassengerRequest) (rides.PassengerRequest, error) {
	q, err := s.Store.SaveRequest(q)
	if err == nil {
		fmt.Fprintf(s.Log, "store: SaveRequest #%d\n", q.ID)
	}
	return q, err
}

func (s LoggingStore) SaveMatch(m ridematching.RideMatch) (ridematching.RideMatch, error) {
	m, err := s.Store.SaveMatch(m)
	if err == nil {
		fmt.Fprintf(s.Log, "store: SaveMatch #%d\n", m.ID)
	}
	return m, err
}
