package main

import (
	"github.com/Shaykhnazar/ulov24/internal/ridematching"
	"github.com/Shaykhnazar/ulov24/internal/rides"
	"github.com/Shaykhnazar/ulov24/internal/store"
)

// Yangi raqam — mavjudlarining eng kattasidan bitta ko'p.

func nextRouteID(s store.Snapshot) int64 {
	var last int64
	for _, r := range s.Routes {
		last = max(last, r.ID)
	}
	return last + 1
}

func nextRequestID(s store.Snapshot) int64 {
	var last int64
	for _, q := range s.Requests {
		last = max(last, q.ID)
	}
	return last + 1
}

func nextMatchID(s store.Snapshot) int64 {
	var last int64
	for _, m := range s.Matches {
		last = max(last, m.ID)
	}
	return last + 1
}

// find* — slice ichidagi elementga pointer: o'zgarish Snapshot'da qoladi.

func findRoute(s store.Snapshot, id int64) *rides.DriverRoute {
	for i := range s.Routes {
		if s.Routes[i].ID == id {
			return &s.Routes[i]
		}
	}
	return nil
}

func findRequest(s store.Snapshot, id int64) *rides.PassengerRequest {
	for i := range s.Requests {
		if s.Requests[i].ID == id {
			return &s.Requests[i]
		}
	}
	return nil
}

func findMatch(s store.Snapshot, id int64) *ridematching.RideMatch {
	for i := range s.Matches {
		if s.Matches[i].ID == id {
			return &s.Matches[i]
		}
	}
	return nil
}

// hasActiveMatch: shu juftlik uchun kutilayotgan yoki qabul qilingan moslik bor.
func hasActiveMatch(s store.Snapshot, routeID, requestID int64) bool {
	for _, m := range s.Matches {
		if m.DriverRouteID != routeID || m.PassengerRequestID != requestID {
			continue
		}
		switch m.OverallStatus {
		case ridematching.MatchPending, ridematching.MatchPartial, ridematching.MatchAccepted:
			return true
		}
	}
	return false
}
