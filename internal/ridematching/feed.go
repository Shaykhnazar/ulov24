package ridematching

import (
	"cmp"
	"slices"
	"time"

	"github.com/Shaykhnazar/ulov24/internal/rides"
)

// Direction — qayerdan qayerga.
type Direction struct {
	From, To int64
}

// IndexByDirection yo'nalishlarni tomoni bo'yicha guruhlaydi, shunda
// "Toshkent → Samarqand" safarlarini butun ro'yxatni aylanmasdan topamiz.
func IndexByDirection(routes []rides.DriverRoute) map[Direction][]rides.DriverRoute {
	index := make(map[Direction][]rides.DriverRoute)
	for _, r := range routes {
		key := Direction{From: r.FromLocationID, To: r.ToLocationID}
		index[key] = append(index[key], r)
	}
	return index
}

// Feed — yo'lovchiga ko'rsatiladigan lenta: faqat ochiq va eskirmagan
// yo'nalishlar; avval "hozir", keyin "bugun", keyin "rejali", har biri ichida
// erta jo'naydigani oldinda. Kiruvchi slice o'zgarmaydi.
func Feed(routes []rides.DriverRoute, now time.Time) []rides.DriverRoute {
	var feed []rides.DriverRoute
	for _, r := range routes {
		if r.Status == rides.DriverRouteOpen && !r.IsExpired(now) {
			feed = append(feed, r)
		}
	}
	slices.SortFunc(feed, func(a, b rides.DriverRoute) int {
		if c := cmp.Compare(b.TripMode.Priority(), a.TripMode.Priority()); c != 0 {
			return c
		}
		return a.DepartureAt.Compare(b.DepartureAt)
	})
	return feed
}
