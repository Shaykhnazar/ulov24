package main

import "github.com/Shaykhnazar/ulov24/internal/ridematching"

// findByID items ichidan id'si mos elementga pointer qaytaradi, topilmasa nil.
// Generic: yo'nalish, so'rov va moslik uchun bitta funksiya.
func findByID[T any](items []T, id int64, getID func(T) int64) *T {
	for i := range items {
		if getID(items[i]) == id {
			return &items[i]
		}
	}
	return nil
}

// hasActiveMatch: shu juftlik uchun kutilayotgan yoki qabul qilingan moslik bor.
func hasActiveMatch(matches []ridematching.RideMatch, routeID, requestID int64) bool {
	for _, m := range matches {
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
