package store

import (
	"slices"

	"github.com/Shaykhnazar/ulov24/internal/ridematching"
	"github.com/Shaykhnazar/ulov24/internal/rides"
)

// MemoryStore hamma narsani xotirada saqlaydi. Dastur yopilsa, ma'lumot yo'qoladi.
type MemoryStore struct {
	data Snapshot
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{}
}

// upsert id 0 bo'lsa yangi raqam beradi va qo'shadi, aks holda almashtiradi.
// Generic: yo'nalish, so'rov va moslik uchun bitta funksiya.
func upsert[T any](items []T, item T, getID func(T) int64, setID func(*T, int64)) ([]T, T) {
	if getID(item) == 0 {
		var last int64
		for _, it := range items {
			last = max(last, getID(it))
		}
		setID(&item, last+1)
		return append(items, item), item
	}
	for i := range items {
		if getID(items[i]) == getID(item) {
			items[i] = item
			return items, item
		}
	}
	return append(items, item), item
}

func (s *MemoryStore) SaveRoute(r rides.DriverRoute) (rides.DriverRoute, error) {
	s.data.Routes, r = upsert(s.data.Routes, r,
		func(r rides.DriverRoute) int64 { return r.ID },
		func(r *rides.DriverRoute, id int64) { r.ID = id })
	return r, nil
}

func (s *MemoryStore) Routes() ([]rides.DriverRoute, error) {
	return slices.Clone(s.data.Routes), nil
}

func (s *MemoryStore) SaveRequest(q rides.PassengerRequest) (rides.PassengerRequest, error) {
	s.data.Requests, q = upsert(s.data.Requests, q,
		func(q rides.PassengerRequest) int64 { return q.ID },
		func(q *rides.PassengerRequest, id int64) { q.ID = id })
	return q, nil
}

func (s *MemoryStore) Requests() ([]rides.PassengerRequest, error) {
	return slices.Clone(s.data.Requests), nil
}

func (s *MemoryStore) SaveMatch(m ridematching.RideMatch) (ridematching.RideMatch, error) {
	s.data.Matches, m = upsert(s.data.Matches, m,
		func(m ridematching.RideMatch) int64 { return m.ID },
		func(m *ridematching.RideMatch, id int64) { m.ID = id })
	return m, nil
}

func (s *MemoryStore) Matches() ([]ridematching.RideMatch, error) {
	return slices.Clone(s.data.Matches), nil
}
