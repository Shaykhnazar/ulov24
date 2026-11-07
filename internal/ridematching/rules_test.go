package ridematching

import (
	"errors"
	"testing"
	"time"

	"github.com/Shaykhnazar/ulov24/internal/rides"
)

func TestDepartureInWindow(t *testing.T) {
	from := testNow
	to := testNow.Add(time.Hour)
	cases := []struct {
		dep  time.Time
		want bool
	}{
		{from, true},
		{to, true},
		{from.Add(-time.Minute), false},
		{to.Add(time.Minute), false},
	}
	for _, c := range cases {
		if got := departureInWindow(c.dep, from, to); got != c.want {
			t.Errorf("departureInWindow(%v) = %v, want %v", c.dep.Format("15:04"), got, c.want)
		}
	}
}

func TestCompatible(t *testing.T) {
	routes, requests := sampleData(testNow)
	fit := routes[0]    // Toshkent → Samarqand, 08:00, 3 seats
	want := requests[0] // Toshkent → Samarqand, 07:00–09:00, 2 people

	if err := Compatible(fit, want, testNow); err != nil {
		t.Fatalf("route 1 and request 1: %v", err)
	}

	cases := []struct {
		name    string
		change  func(r *rides.DriverRoute, q *rides.PassengerRequest)
		wantErr error
	}{
		{"other direction", func(r *rides.DriverRoute, q *rides.PassengerRequest) { r.ToLocationID = rides.CityBukhara }, ErrDirectionMismatch},
		{"outside window", func(r *rides.DriverRoute, q *rides.PassengerRequest) { r.DepartureAt = testNow.Add(4 * time.Hour) }, ErrOutsideWindow},
		{"expired", func(r *rides.DriverRoute, q *rides.PassengerRequest) {
			r.DepartureAt = testNow.Add(-time.Minute)
			q.DesiredDepartureFrom = testNow.Add(-time.Hour)
		}, ErrRouteExpired},
		{"not open", func(r *rides.DriverRoute, q *rides.PassengerRequest) { r.Status = rides.DriverRouteCancelled }, ErrRouteNotOpen},
		{"not enough seats", func(r *rides.DriverRoute, q *rides.PassengerRequest) { r.SeatsRemaining = 1 }, ErrRideCapacityFull},
		{"same user", func(r *rides.DriverRoute, q *rides.PassengerRequest) { q.PassengerUserID = r.DriverUserID }, ErrSameUser},
	}
	for _, c := range cases {
		r, q := fit, want
		c.change(&r, &q)
		if err := Compatible(r, q, testNow); !errors.Is(err, c.wantErr) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.wantErr)
		}
	}
}

func TestCapacityErrorIsShared(t *testing.T) {
	if !errors.Is(ErrRideCapacityFull, rides.ErrRideCapacityFull) {
		t.Fatal("ridematching and rides must share one ErrRideCapacityFull")
	}
}
