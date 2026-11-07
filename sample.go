package main

import "time"

// sampleData — namuna ma'lumotlar: to'rtta yo'nalish va to'rtta so'rov.
// Vaqtlar hozirgi vaqtdan hisoblanadi, shuning uchun dastur har kuni ishlaydi.
// 7-modulda bular faylga ko'chadi.
func sampleData(now time.Time) ([]DriverRoute, []PassengerRequest) {
	routes := []DriverRoute{
		{ID: 1, DriverUserID: 101, FromLocationID: CityTashkent, ToLocationID: CitySamarkand,
			DepartureAt: now.Add(2 * time.Hour), SeatsTotal: 3, SeatsRemaining: 3,
			PricePerSeat: 150000, TripMode: TripModeToday, Status: DriverRouteOpen},
		{ID: 2, DriverUserID: 102, FromLocationID: CityTashkent, ToLocationID: CitySamarkand,
			DepartureAt: now.Add(30 * time.Minute), SeatsTotal: 1, SeatsRemaining: 1,
			PricePerSeat: 140000, TripMode: TripModeNow, Status: DriverRouteOpen},
		{ID: 3, DriverUserID: 103, FromLocationID: CityTashkent, ToLocationID: CityBukhara,
			DepartureAt: now.Add(26 * time.Hour), SeatsTotal: 4, SeatsRemaining: 4,
			PricePerSeat: 200000, TripMode: TripModePlanned, Status: DriverRouteOpen},
		{ID: 4, DriverUserID: 104, FromLocationID: CityFergana, ToLocationID: CityTashkent,
			DepartureAt: now.Add(-1 * time.Hour), SeatsTotal: 2, SeatsRemaining: 2,
			PricePerSeat: 120000, TripMode: TripModeToday, Status: DriverRouteOpen},
	}
	requests := []PassengerRequest{
		{ID: 1, PassengerUserID: 201, FromLocationID: CityTashkent, ToLocationID: CitySamarkand,
			DesiredDepartureFrom: now.Add(1 * time.Hour), DesiredDepartureTo: now.Add(3 * time.Hour),
			PassengerCount: 2, TripMode: TripModeToday, Status: PassengerRequestOpen},
		{ID: 2, PassengerUserID: 202, FromLocationID: CityTashkent, ToLocationID: CitySamarkand,
			DesiredDepartureFrom: now, DesiredDepartureTo: now.Add(1 * time.Hour),
			PassengerCount: 1, TripMode: TripModeNow, Status: PassengerRequestOpen},
		{ID: 3, PassengerUserID: 203, FromLocationID: CityTashkent, ToLocationID: CityBukhara,
			DesiredDepartureFrom: now.Add(24 * time.Hour), DesiredDepartureTo: now.Add(30 * time.Hour),
			PassengerCount: 5, TripMode: TripModePlanned, Status: PassengerRequestOpen},
		{ID: 4, PassengerUserID: 204, FromLocationID: CityAndijan, ToLocationID: CityTashkent,
			DesiredDepartureFrom: now, DesiredDepartureTo: now.Add(4 * time.Hour),
			PassengerCount: 1, TripMode: TripModeToday, Status: PassengerRequestOpen},
	}
	return routes, requests
}
