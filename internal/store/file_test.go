package store

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Shaykhnazar/ulov24/internal/ridematching"
	"github.com/Shaykhnazar/ulov24/internal/rides"
)

func TestLoadMissingFile(t *testing.T) {
	s, err := LoadFile(filepath.Join(t.TempDir(), "yoq.json"))
	if err != nil {
		t.Fatalf("missing file: %v", err)
	}
	if len(s.Routes)+len(s.Requests)+len(s.Matches) != 0 {
		t.Fatalf("missing file must give an empty snapshot, got %+v", s)
	}
}

func TestLoadCorruptFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "d.json")
	if err := os.WriteFile(path, []byte("{bad"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadFile(path); !errors.Is(err, ErrCorruptData) {
		t.Fatalf("err = %v, want ErrCorruptData", err)
	}
}

func TestSaveLoadRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "d.json")
	price := 150000
	dep := time.Date(2026, 3, 20, 8, 0, 0, 0, rides.Tashkent)
	in := Snapshot{
		Routes: []rides.DriverRoute{
			{ID: 1, DriverUserID: 1, FromLocationID: 1, ToLocationID: 2, DepartureAt: dep,
				SeatsTotal: 3, SeatsRemaining: 3, PricePerSeat: &price, TripMode: rides.TripModeToday, Status: rides.DriverRouteOpen},
			{ID: 2, DriverUserID: 1, FromLocationID: 1, ToLocationID: 3, DepartureAt: dep,
				SeatsTotal: 2, SeatsRemaining: 2, TripMode: rides.TripModeToday, Status: rides.DriverRouteOpen},
		},
		Matches: []ridematching.RideMatch{{
			RideMatchState: ridematching.RideMatchState{ID: 1, DriverStatus: ridematching.DualStatusPending,
				PassengerStatus: ridematching.DualStatusAccepted, OverallStatus: ridematching.MatchPartial},
			DriverRouteID: 1, PassengerRequestID: 1, PassengerCount: 2,
		}},
	}
	if err := SaveFile(path, in); err != nil {
		t.Fatal(err)
	}
	out, err := LoadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(out.Routes) != 2 || out.Routes[0].PricePerSeat == nil || *out.Routes[0].PricePerSeat != price ||
		out.Routes[1].PricePerSeat != nil || !out.Routes[0].DepartureAt.Equal(dep) {
		t.Fatalf("routes did not round-trip: %+v", out.Routes)
	}
	if len(out.Matches) != 1 || out.Matches[0] != in.Matches[0] {
		t.Fatalf("matches did not round-trip: %+v", out.Matches)
	}
	if _, err := os.Stat(path + ".tmp"); !os.IsNotExist(err) {
		t.Fatalf("temporary file left behind: %v", err)
	}
}
