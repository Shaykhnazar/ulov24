package main

import (
	"testing"
	"time"
)

// testNow is the same moment as fixedNow in main_test.go.
var testNow = time.Date(2026, 3, 20, 6, 0, 0, 0, Tashkent)

func TestDirectionFilter(t *testing.T) {
	routes, _ := sampleData(testNow)
	keep := directionFilter(CityTashkent, CitySamarkand)
	want := map[int64]bool{1: true, 2: true, 3: false, 4: false}
	for _, r := range routes {
		if keep(r) != want[r.ID] {
			t.Errorf("route %d: got %v, want %v", r.ID, keep(r), want[r.ID])
		}
	}
}

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

func TestNewIDGenerator(t *testing.T) {
	a := newIDGenerator()
	b := newIDGenerator()
	if a() != 1 || a() != 2 || a() != 3 {
		t.Fatal("generator a must count 1, 2, 3")
	}
	if b() != 1 {
		t.Fatal("generator b must start at 1 on its own")
	}
}

func TestCompatibleSample(t *testing.T) {
	routes, requests := sampleData(testNow)
	if !compatible(routes[0], requests[0], testNow) {
		t.Error("route 1 must fit request 1")
	}
	if compatible(routes[2], requests[2], testNow) {
		t.Error("route 3 has 4 seats, request 3 needs 5")
	}
	if routes[0].IsExpired(testNow) || !routes[3].IsExpired(testNow) {
		t.Error("only route 4 is expired")
	}
}
