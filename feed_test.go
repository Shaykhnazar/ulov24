package main

import (
	"slices"
	"testing"
)

func routeIDs(routes []DriverRoute) []int64 {
	var ids []int64
	for _, r := range routes {
		ids = append(ids, r.ID)
	}
	return ids
}

func TestFeedOrder(t *testing.T) {
	routes, _ := sampleData(testNow)
	got := routeIDs(Feed(routes, testNow))
	if want := []int64{2, 1, 3}; !slices.Equal(got, want) {
		t.Fatalf("Feed ids = %v, want %v", got, want)
	}
}

func TestFeedDoesNotReorderInput(t *testing.T) {
	routes, _ := sampleData(testNow)
	Feed(routes, testNow)
	if got := routeIDs(routes); !slices.Equal(got, []int64{1, 2, 3, 4}) {
		t.Fatalf("input reordered: %v", got)
	}
}

func TestIndexByDirection(t *testing.T) {
	routes, _ := sampleData(testNow)
	index := IndexByDirection(routes)
	got := routeIDs(index[Direction{From: CityTashkent, To: CitySamarkand}])
	if want := []int64{1, 2}; !slices.Equal(got, want) {
		t.Fatalf("Toshkent → Samarqand = %v, want %v", got, want)
	}
}
