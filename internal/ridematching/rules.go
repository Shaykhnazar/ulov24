package ridematching

import (
	"errors"
	"time"

	"github.com/Shaykhnazar/ulov24/internal/rides"
)

// sameDirection: yo'nalish va so'rov bir xil shahardan bir xil shaharga.
func sameDirection(r rides.DriverRoute, q rides.PassengerRequest) bool {
	return r.FromLocationID == q.FromLocationID && r.ToLocationID == q.ToLocationID
}

// departureInWindow: jo'nash vaqti [from, to] oralig'ida (chegaralar ham kiradi).
func departureInWindow(dep, from, to time.Time) bool {
	return !dep.Before(from) && !dep.After(to)
}

// Compatible nima uchun mos kelmasligini ham aytadi: nil — mos, aks holda
// quyidagi xatolardan biri. errors.Is bilan tekshirish mumkin.
var (
	ErrDirectionMismatch = errors.New("route goes another way")
	ErrRouteNotOpen      = errors.New("route is not open")
	ErrRouteExpired      = errors.New("route has already departed")
	ErrOutsideWindow     = errors.New("departure is outside the requested time window")
	ErrSameUser          = errors.New("driver and passenger are the same user")
)

// Compatible yo'nalish shu so'rovga to'g'ri keladimi, tekshiradi.
func Compatible(r rides.DriverRoute, q rides.PassengerRequest, now time.Time) error {
	switch {
	case !sameDirection(r, q):
		return ErrDirectionMismatch
	case r.Status != rides.DriverRouteOpen:
		return ErrRouteNotOpen
	case r.IsExpired(now):
		return ErrRouteExpired
	case !departureInWindow(r.DepartureAt, q.DesiredDepartureFrom, q.DesiredDepartureTo):
		return ErrOutsideWindow
	case r.SeatsRemaining < q.PassengerCount:
		return ErrRideCapacityFull
	case r.DriverUserID == q.PassengerUserID:
		return ErrSameUser
	}
	return nil
}
