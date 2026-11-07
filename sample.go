package main

import "time"

// sampleData — namuna ma'lumotlar: to'rtta yo'nalish va to'rtta so'rov.
// Vaqtlar hozirgi vaqtdan hisoblanadi, shuning uchun dastur har kuni ishlaydi.
// 7-modulda bular faylga ko'chadi.
func sampleData(now time.Time) ([]DriverRoute, []PassengerRequest) {
	routeID := newIDGenerator()
	requestID := newIDGenerator()
	routes := []DriverRoute{
		{ID: routeID(), DriverUserID: 101, FromLocationID: CityTashkent, ToLocationID: CitySamarkand,
			DepartureAt: now.Add(2 * time.Hour), SeatsTotal: 3, SeatsRemaining: 3,
			PricePerSeat: 150000, TripMode: TripModeToday, Status: DriverRouteOpen},
		{ID: routeID(), DriverUserID: 102, FromLocationID: CityTashkent, ToLocationID: CitySamarkand,
			DepartureAt: now.Add(30 * time.Minute), SeatsTotal: 1, SeatsRemaining: 1,
			PricePerSeat: 140000, TripMode: TripModeNow, Status: DriverRouteOpen},
		{ID: routeID(), DriverUserID: 103, FromLocationID: CityTashkent, ToLocationID: CityBukhara,
			DepartureAt: now.Add(26 * time.Hour), SeatsTotal: 4, SeatsRemaining: 4,
			PricePerSeat: 200000, TripMode: TripModePlanned, Status: DriverRouteOpen},
		{ID: routeID(), DriverUserID: 104, FromLocationID: CityFergana, ToLocationID: CityTashkent,
			DepartureAt: now.Add(-1 * time.Hour), SeatsTotal: 2, SeatsRemaining: 2,
			PricePerSeat: 120000, TripMode: TripModeToday, Status: DriverRouteOpen},
	}
	requests := []PassengerRequest{
		{ID: requestID(), PassengerUserID: 201, FromLocationID: CityTashkent, ToLocationID: CitySamarkand,
			DesiredDepartureFrom: now.Add(1 * time.Hour), DesiredDepartureTo: now.Add(3 * time.Hour),
			PassengerCount: 2, TripMode: TripModeToday, Status: PassengerRequestOpen},
		{ID: requestID(), PassengerUserID: 202, FromLocationID: CityTashkent, ToLocationID: CitySamarkand,
			DesiredDepartureFrom: now, DesiredDepartureTo: now.Add(1 * time.Hour),
			PassengerCount: 1, TripMode: TripModeNow, Status: PassengerRequestOpen},
		{ID: requestID(), PassengerUserID: 203, FromLocationID: CityTashkent, ToLocationID: CityBukhara,
			DesiredDepartureFrom: now.Add(24 * time.Hour), DesiredDepartureTo: now.Add(30 * time.Hour),
			PassengerCount: 5, TripMode: TripModePlanned, Status: PassengerRequestOpen},
		{ID: requestID(), PassengerUserID: 204, FromLocationID: CityAndijan, ToLocationID: CityTashkent,
			DesiredDepartureFrom: now, DesiredDepartureTo: now.Add(4 * time.Hour),
			PassengerCount: 1, TripMode: TripModeToday, Status: PassengerRequestOpen},
	}
	return routes, requests
}
