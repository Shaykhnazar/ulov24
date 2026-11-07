package main

import (
	"testing"
	"time"
)

func TestValidateSeats(t *testing.T) {
	for _, n := range []int{0, 9} {
		if err := validateSeats(n); err == nil {
			t.Errorf("validateSeats(%d) must fail", n)
		}
	}
	for _, n := range []int{1, 8} {
		if err := validateSeats(n); err != nil {
			t.Errorf("validateSeats(%d) = %v", n, err)
		}
	}
	if err := validateSeats(9); err.Error() != "joylar soni 1 dan 8 gacha bo'lishi kerak: 9" {
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
		err := validateTripMode(c.mode, c.dep, testNow)
		got := ""
		if err != nil {
			got = err.Error()
		}
		if got != c.wantErr {
			t.Errorf("%s: got %q, want %q", c.name, got, c.wantErr)
		}
	}
}
