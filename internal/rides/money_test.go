package rides

import "testing"

func TestFormatSum(t *testing.T) {
	cases := []struct {
		amount int
		want   string
	}{
		{0, "0 so'm"},
		{999, "999 so'm"},
		{1000, "1 000 so'm"},
		{150000, "150 000 so'm"},
		{1234567, "1 234 567 so'm"},
	}
	for _, c := range cases {
		if got := FormatSum(c.amount); got != c.want {
			t.Errorf("FormatSum(%d) = %q, want %q", c.amount, got, c.want)
		}
	}
}

func TestTripModePriority(t *testing.T) {
	if !(TripModeNow.Priority() > TripModeToday.Priority() &&
		TripModeToday.Priority() > TripModePlanned.Priority() &&
		TripModePlanned.Priority() > TripMode("x").Priority()) {
		t.Fatal("priority must be now > today > planned > unknown")
	}
}
