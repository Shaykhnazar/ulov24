package main

import "time"

// sameDirection: yo'nalish va so'rov bir xil shahardan bir xil shaharga.
func sameDirection(r DriverRoute, q PassengerRequest) bool {
	return r.FromLocationID == q.FromLocationID && r.ToLocationID == q.ToLocationID
}

// departureInWindow: jo'nash vaqti [from, to] oralig'ida (chegaralar ham kiradi).
func departureInWindow(dep, from, to time.Time) bool {
	return !dep.Before(from) && !dep.After(to)
}

// isRouteExpired: jo'nash vaqti o'tib ketgan.
func isRouteExpired(r DriverRoute, now time.Time) bool {
	return now.After(r.DepartureAt)
}

// compatible: yo'nalish shu so'rovga to'g'ri keladimi.
func compatible(r DriverRoute, q PassengerRequest, now time.Time) bool {
	return sameDirection(r, q) &&
		r.Status == DriverRouteOpen &&
		!isRouteExpired(r, now) &&
		departureInWindow(r.DepartureAt, q.DesiredDepartureFrom, q.DesiredDepartureTo) &&
		r.SeatsRemaining >= q.PassengerCount &&
		r.DriverUserID != q.PassengerUserID
}

// directionFilter faqat from → to tomonga ketadigan yo'nalishlarni qoldiradigan
// funksiya qaytaradi. Qaytgan funksiya from va to'ni "eslab qoladi" — bu closure.
func directionFilter(from, to int64) func(DriverRoute) bool {
	return func(r DriverRoute) bool {
		return r.FromLocationID == from && r.ToLocationID == to
	}
}

// newIDGenerator har chaqirilganda keyingi raqamni beradigan funksiya qaytaradi:
// 1, 2, 3, … Har bir generator o'z hisobini alohida yuritadi.
func newIDGenerator() func() int64 {
	var last int64
	return func() int64 {
		last++
		return last
	}
}
