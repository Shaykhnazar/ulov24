package ridematching

import (
	"time"

	"github.com/Shaykhnazar/ulov24/internal/rides"
)

// testNow: 20 March 2026, 06:00 in Tashkent.
var testNow = time.Date(2026, 3, 20, 6, 0, 0, 0, rides.Tashkent)

func price(n int) *int { return &n }

// sampleData — the four routes and four requests the course has used since
// module 3, kept as a test fixture.
func sampleData(now time.Time) ([]rides.DriverRoute, []rides.PassengerRequest) {
	routes := []rides.DriverRoute{
		{ID: 1, DriverUserID: 101, FromLocationID: rides.CityTashkent, ToLocationID: rides.CitySamarkand,
			DepartureAt: now.Add(2 * time.Hour), SeatsTotal: 3, SeatsRemaining: 3,
			PricePerSeat: price(150000), TripMode: rides.TripModeToday, Status: rides.DriverRouteOpen},
		{ID: 2, DriverUserID: 102, FromLocationID: rides.CityTashkent, ToLocationID: rides.CitySamarkand,
			DepartureAt: now.Add(30 * time.Minute), SeatsTotal: 1, SeatsRemaining: 1,
			PricePerSeat: price(140000), TripMode: rides.TripModeNow, Status: rides.DriverRouteOpen},
		{ID: 3, DriverUserID: 103, FromLocationID: rides.CityTashkent, ToLocationID: rides.CityBukhara,
			DepartureAt: now.Add(26 * time.Hour), SeatsTotal: 4, SeatsRemaining: 4,
			PricePerSeat: price(200000), TripMode: rides.TripModePlanned, Status: rides.DriverRouteOpen},
		{ID: 4, DriverUserID: 104, FromLocationID: rides.CityFergana, ToLocationID: rides.CityTashkent,
			DepartureAt: now.Add(-1 * time.Hour), SeatsTotal: 2, SeatsRemaining: 2,
			PricePerSeat: price(120000), TripMode: rides.TripModeToday, Status: rides.DriverRouteOpen},
	}
	requests := []rides.PassengerRequest{
		{ID: 1, PassengerUserID: 201, FromLocationID: rides.CityTashkent, ToLocationID: rides.CitySamarkand,
			DesiredDepartureFrom: now.Add(1 * time.Hour), DesiredDepartureTo: now.Add(3 * time.Hour),
			PassengerCount: 2, TripMode: rides.TripModeToday, Status: rides.PassengerRequestOpen},
		{ID: 2, PassengerUserID: 202, FromLocationID: rides.CityTashkent, ToLocationID: rides.CitySamarkand,
			DesiredDepartureFrom: now, DesiredDepartureTo: now.Add(1 * time.Hour),
			PassengerCount: 1, TripMode: rides.TripModeNow, Status: rides.PassengerRequestOpen},
		{ID: 3, PassengerUserID: 203, FromLocationID: rides.CityTashkent, ToLocationID: rides.CityBukhara,
			DesiredDepartureFrom: now.Add(24 * time.Hour), DesiredDepartureTo: now.Add(30 * time.Hour),
			PassengerCount: 5, TripMode: rides.TripModePlanned, Status: rides.PassengerRequestOpen},
		{ID: 4, PassengerUserID: 204, FromLocationID: rides.CityAndijan, ToLocationID: rides.CityTashkent,
			DesiredDepartureFrom: now, DesiredDepartureTo: now.Add(4 * time.Hour),
			PassengerCount: 1, TripMode: rides.TripModeToday, Status: rides.PassengerRequestOpen},
	}
	return routes, requests
}
