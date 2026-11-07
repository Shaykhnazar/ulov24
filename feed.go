package main

import (
	"cmp"
	"fmt"
	"slices"
	"time"
)

// Direction — qayerdan qayerga.
type Direction struct {
	From, To int64
}

// IndexByDirection yo'nalishlarni tomoni bo'yicha guruhlaydi, shunda
// "Toshkent → Samarqand" safarlarini butun ro'yxatni aylanmasdan topamiz.
func IndexByDirection(routes []DriverRoute) map[Direction][]DriverRoute {
	index := make(map[Direction][]DriverRoute)
	for _, r := range routes {
		key := Direction{From: r.FromLocationID, To: r.ToLocationID}
		index[key] = append(index[key], r)
	}
	return index
}

// Feed — yo'lovchiga ko'rsatiladigan lenta: faqat ochiq va eskirmagan
// yo'nalishlar; avval "hozir", keyin "bugun", keyin "rejali", har biri ichida
// erta jo'naydigani oldinda. Kiruvchi slice o'zgarmaydi.
func Feed(routes []DriverRoute, now time.Time) []DriverRoute {
	var feed []DriverRoute
	for _, r := range routes {
		if r.Status == DriverRouteOpen && !isRouteExpired(r, now) {
			feed = append(feed, r)
		}
	}
	slices.SortFunc(feed, func(a, b DriverRoute) int {
		if c := cmp.Compare(b.TripMode.Priority(), a.TripMode.Priority()); c != 0 {
			return c
		}
		return a.DepartureAt.Compare(b.DepartureAt)
	})
	return feed
}

// routeLine — bitta yo'nalishning lentadagi qatori.
func routeLine(r DriverRoute, label string) string {
	return fmt.Sprintf("#%d %s → %s · %s · %d joy · %s · %s",
		r.ID, CityNames[r.FromLocationID], CityNames[r.ToLocationID],
		r.DepartureAt.In(Tashkent).Format("02.01 15:04"),
		r.SeatsRemaining, FormatSum(r.PricePerSeat), label)
}
