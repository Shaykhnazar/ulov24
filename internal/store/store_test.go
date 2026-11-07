package store

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Shaykhnazar/ulov24/internal/rides"
)

func testRoute(seats int) rides.DriverRoute {
	return rides.DriverRoute{
		DriverUserID: 1, FromLocationID: 1, ToLocationID: 2,
		DepartureAt: time.Date(2026, 3, 20, 8, 0, 0, 0, rides.Tashkent),
		SeatsTotal:  seats, SeatsRemaining: seats,
		TripMode: rides.TripModeToday, Status: rides.DriverRouteOpen,
	}
}

func TestStoreContract(t *testing.T) {
	stores := map[string]func(t *testing.T) Store{
		"memory": func(t *testing.T) Store { return NewMemoryStore() },
		"file": func(t *testing.T) Store {
			s, err := OpenFileStore(filepath.Join(t.TempDir(), "d.json"))
			if err != nil {
				t.Fatal(err)
			}
			return s
		},
	}
	for name, open := range stores {
		t.Run(name, func(t *testing.T) {
			s := open(t)
			first, err := s.SaveRoute(testRoute(3))
			if err != nil || first.ID != 1 {
				t.Fatalf("first SaveRoute: id=%d err=%v", first.ID, err)
			}
			second, err := s.SaveRoute(testRoute(2))
			if err != nil || second.ID != 2 {
				t.Fatalf("second SaveRoute: id=%d err=%v", second.ID, err)
			}

			first.SeatsRemaining = 1
			if _, err := s.SaveRoute(first); err != nil {
				t.Fatal(err)
			}
			routes, err := s.Routes()
			if err != nil || len(routes) != 2 || routes[0].SeatsRemaining != 1 {
				t.Fatalf("after re-save: %+v, err=%v", routes, err)
			}

			routes[0].SeatsRemaining = 99 // the caller's copy, not the store
			again, _ := s.Routes()
			if again[0].SeatsRemaining != 1 {
				t.Fatal("Routes() must return a copy")
			}
		})
	}
}

func TestFileStorePersists(t *testing.T) {
	path := filepath.Join(t.TempDir(), "d.json")
	s, err := OpenFileStore(path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.SaveRoute(testRoute(3)); err != nil {
		t.Fatal(err)
	}
	reopened, err := OpenFileStore(path)
	if err != nil {
		t.Fatal(err)
	}
	routes, _ := reopened.Routes()
	if len(routes) != 1 || routes[0].ID != 1 {
		t.Fatalf("reopened store: %+v", routes)
	}
}

func TestLoggingStore(t *testing.T) {
	var log bytes.Buffer
	s := LoggingStore{Store: NewMemoryStore(), Log: &log}
	if _, err := s.SaveRoute(testRoute(3)); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(log.String(), "store: SaveRoute #1") {
		t.Fatalf("log = %q", log.String())
	}
}
