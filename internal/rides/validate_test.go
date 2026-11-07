package rides

import (
	"testing"
	"time"
)

// testNow: 20 March 2026, 06:00 in Tashkent.
var testNow = time.Date(2026, 3, 20, 6, 0, 0, 0, Tashkent)

func TestValidateSeats(t *testing.T) {
	for _, n := range []int{0, 9} {
		if err := ValidateSeats(n); err == nil {
			t.Errorf("ValidateSeats(%d) must fail", n)
		}
	}
	for _, n := range []int{1, 8} {
		if err := ValidateSeats(n); err != nil {
			t.Errorf("ValidateSeats(%d) = %v", n, err)
		}
	}
	if err := ValidateSeats(9); err.Error() != "joylar soni 1 dan 8 gacha bo'lishi kerak: 9" {
		t.Errorf("message = %q", err)
	}
}

func TestValidateTripMode(t *testing.T) {
	sameDay2359 := time.Date(2026, 3, 20, 23, 59, 0, 0, Tashkent)
	nextDay0000 := time.Date(2026, 3, 21, 0, 0, 0, 0, Tashkent)
	cases := []struct {
		name    string
		mode    TripMode
		dep     time.Time
		wantErr string
	}{
		{"now +3h", TripModeNow, testNow.Add(3 * time.Hour), ""},
		{"now +3h01m", TripModeNow, testNow.Add(3*time.Hour + time.Minute), `"hozir" rejimida jo'nash 3 soat ichida bo'lishi kerak`},
		{"today 23:59", TripModeToday, sameDay2359, ""},
		{"today next day", TripModeToday, nextDay0000, `"bugun" rejimida jo'nash bugun bo'lishi kerak`},
		{"planned -1m", TripModePlanned, testNow.Add(-time.Minute), `"rejali" rejimida jo'nash kelajakda bo'lishi kerak`},
		{"unknown", TripMode("x"), testNow.Add(time.Hour), "noma'lum rejim: x"},
	}
	for _, c := range cases {
		err := ValidateTripMode(c.mode, c.dep, testNow)
		got := ""
		if err != nil {
			got = err.Error()
		}
		if got != c.wantErr {
			t.Errorf("%s: got %q, want %q", c.name, got, c.wantErr)
		}
	}
}

func TestParseLocalTime(t *testing.T) {
	got, err := ParseLocalTime("2026-03-20 08:00")
	if err != nil || !got.Equal(testNow.Add(2*time.Hour)) {
		t.Fatalf("got %v, %v", got, err)
	}
	if _, err := ParseLocalTime("20.03 08:00"); err == nil {
		t.Fatal("20.03 08:00 must not parse")
	}
}

func TestReserve(t *testing.T) {
	r := DriverRoute{SeatsTotal: 2, SeatsRemaining: 2, Status: DriverRouteOpen}
	if err := r.Reserve(3); err != ErrRideCapacityFull || r.SeatsRemaining != 2 {
		t.Fatalf("over-reserve: err=%v seats=%d", err, r.SeatsRemaining)
	}
	if err := r.Reserve(2); err != nil || r.SeatsRemaining != 0 || r.Status != DriverRouteFull {
		t.Fatalf("reserve all: err=%v seats=%d status=%s", err, r.SeatsRemaining, r.Status)
	}
}
