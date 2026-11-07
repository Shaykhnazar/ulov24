package ridematching

import (
	"errors"

	"github.com/Shaykhnazar/ulov24/internal/rides"
)

// Bu fayldagi turlar va ApplyAcceptance manapochta'dagi Ulov24 kodidan
// (internal/ridematching) aynan shu nomlar bilan olingan: 2-kursda ular
// haqiqiy ma'lumotlar bazasiga ulanganda nom o'zgarmaydi.

var (
	ErrUnauthorizedActor = errors.New("actor is not part of ride match")
	ErrTerminalRideMatch = errors.New("ride match is terminal")
	ErrRideCapacityFull  = rides.ErrRideCapacityFull // bitta qiymat ikkala paketda
)

// UserRole — moslikda kim harakat qilyapti.
type UserRole string

const (
	RoleDriver    UserRole = "driver"
	RolePassenger UserRole = "passenger"
)

// DualStatus — bir tomonning (haydovchi yoki yo'lovchi) javobi.
type DualStatus string

const (
	DualStatusPending   DualStatus = "pending"
	DualStatusAccepted  DualStatus = "accepted"
	DualStatusRejected  DualStatus = "rejected"
	DualStatusCancelled DualStatus = "cancelled"
)

// MatchStatus — moslikning umumiy holati.
type MatchStatus string

const (
	MatchPending            MatchStatus = "pending"
	MatchPartial            MatchStatus = "partial"
	MatchAccepted           MatchStatus = "accepted"
	MatchRejected           MatchStatus = "rejected"
	MatchCancelled          MatchStatus = "cancelled"
	MatchCompleted          MatchStatus = "completed"
	MatchClosedCapacityFull MatchStatus = "closed_capacity_full"
)

// RideMatchState — moslikning ikki tomondan javoblari va umumiy holati.
type RideMatchState struct {
	ID              int64       `json:"id"`
	DriverStatus    DualStatus  `json:"driver_status"`
	PassengerStatus DualStatus  `json:"passenger_status"`
	OverallStatus   MatchStatus `json:"overall_status"`
}

// ApplyAcceptance bir tomonning "qabul qildim" javobini qo'llaydi va yangi
// holatni qaytaradi. Ikkinchi qiymat true bo'lsa, ikkala tomon ham qabul qildi:
// endi joy band qilinishi kerak.
func ApplyAcceptance(state RideMatchState, role UserRole) (RideMatchState, bool, error) {
	if isTerminal(state.OverallStatus) {
		return RideMatchState{}, false, ErrTerminalRideMatch
	}

	next := state
	switch role {
	case RoleDriver:
		next.DriverStatus = DualStatusAccepted
	case RolePassenger:
		next.PassengerStatus = DualStatusAccepted
	default:
		return RideMatchState{}, false, ErrUnauthorizedActor
	}

	if next.DriverStatus == DualStatusAccepted && next.PassengerStatus == DualStatusAccepted {
		next.OverallStatus = MatchAccepted
		return next, true, nil
	}

	next.OverallStatus = MatchPartial
	return next, false, nil
}

func isTerminal(status MatchStatus) bool {
	switch status {
	case MatchAccepted, MatchRejected, MatchCancelled, MatchCompleted, MatchClosedCapacityFull:
		return true
	default:
		return false
	}
}

// RideMatch — bitta yo'nalish va bitta so'rov orasidagi moslik.
type RideMatch struct {
	RideMatchState
	DriverRouteID      int64 `json:"driver_route_id"`
	PassengerRequestID int64 `json:"passenger_request_id"`
	PassengerCount     int   `json:"passenger_count"`
}

// Accept bir tomonning roziligini yozadi. Ikkala tomon rozi bo'lganda
// yo'nalishdan joy band qilinadi. route pointer orqali keladi: aks holda biz
// yo'nalishning nusxasini o'zgartirgan bo'lardik, asl joylar soni esa
// o'zgarmay qolardi.
func (m *RideMatch) Accept(role UserRole, route *rides.DriverRoute) error {
	next, bothAccepted, err := ApplyAcceptance(m.RideMatchState, role)
	if err != nil {
		return err
	}
	if bothAccepted {
		if err := route.Reserve(m.PassengerCount); err != nil {
			m.OverallStatus = MatchClosedCapacityFull
			return err
		}
	}
	m.RideMatchState = next
	return nil
}
