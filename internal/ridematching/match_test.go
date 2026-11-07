package ridematching

import (
	"errors"
	"testing"

	"github.com/Shaykhnazar/ulov24/internal/rides"
)

func TestApplyAcceptance(t *testing.T) {
	pending := RideMatchState{ID: 1, DriverStatus: DualStatusPending, PassengerStatus: DualStatusPending, OverallStatus: MatchPending}

	afterPassenger, reserved, err := ApplyAcceptance(pending, RolePassenger)
	if err != nil || reserved || afterPassenger.OverallStatus != MatchPartial || afterPassenger.PassengerStatus != DualStatusAccepted {
		t.Fatalf("passenger first: %+v, reserved=%v, err=%v", afterPassenger, reserved, err)
	}

	afterDriver, reserved, err := ApplyAcceptance(afterPassenger, RoleDriver)
	if err != nil || !reserved || afterDriver.OverallStatus != MatchAccepted {
		t.Fatalf("then driver: %+v, reserved=%v, err=%v", afterDriver, reserved, err)
	}

	for _, role := range []UserRole{RoleDriver, RolePassenger} {
		if _, _, err := ApplyAcceptance(afterDriver, role); !errors.Is(err, ErrTerminalRideMatch) {
			t.Errorf("%s on terminal state: err = %v", role, err)
		}
	}

	if _, _, err := ApplyAcceptance(pending, UserRole("x")); !errors.Is(err, ErrUnauthorizedActor) {
		t.Errorf("unknown role: err = %v", err)
	}
}

func newTestMatch(count int) RideMatch {
	return RideMatch{
		RideMatchState: RideMatchState{ID: 1, DriverStatus: DualStatusPending, PassengerStatus: DualStatusPending, OverallStatus: MatchPending},
		DriverRouteID:  1, PassengerRequestID: 1, PassengerCount: count,
	}
}

func TestAcceptReservesOnce(t *testing.T) {
	route := rides.DriverRoute{ID: 1, SeatsTotal: 3, SeatsRemaining: 3, Status: rides.DriverRouteOpen}
	m := newTestMatch(2)

	if err := m.Accept(RolePassenger, &route); err != nil || route.SeatsRemaining != 3 {
		t.Fatalf("passenger accept: err=%v seats=%d", err, route.SeatsRemaining)
	}
	if err := m.Accept(RoleDriver, &route); err != nil || route.SeatsRemaining != 1 {
		t.Fatalf("driver accept: err=%v seats=%d", err, route.SeatsRemaining)
	}
	if err := m.Accept(RoleDriver, &route); !errors.Is(err, ErrTerminalRideMatch) || route.SeatsRemaining != 1 {
		t.Fatalf("second driver accept: err=%v seats=%d", err, route.SeatsRemaining)
	}
}

func TestAcceptCapacityFull(t *testing.T) {
	route := rides.DriverRoute{ID: 1, SeatsTotal: 1, SeatsRemaining: 1, Status: rides.DriverRouteOpen}
	m := newTestMatch(2)
	if err := m.Accept(RolePassenger, &route); err != nil {
		t.Fatal(err)
	}
	err := m.Accept(RoleDriver, &route)
	if !errors.Is(err, ErrRideCapacityFull) {
		t.Fatalf("err = %v, want ErrRideCapacityFull", err)
	}
	if m.OverallStatus != MatchClosedCapacityFull || route.SeatsRemaining != 1 {
		t.Fatalf("status=%s seats=%d", m.OverallStatus, route.SeatsRemaining)
	}
}

func TestReserveToZeroMarksFull(t *testing.T) {
	route := rides.DriverRoute{ID: 1, SeatsTotal: 2, SeatsRemaining: 2, Status: rides.DriverRouteOpen}
	if err := route.Reserve(2); err != nil {
		t.Fatal(err)
	}
	if route.SeatsRemaining != 0 || route.Status != rides.DriverRouteFull {
		t.Fatalf("seats=%d status=%s", route.SeatsRemaining, route.Status)
	}
	if err := route.Reserve(1); !errors.Is(err, ErrRideCapacityFull) || route.SeatsRemaining != 0 {
		t.Fatalf("reserve on full route: err=%v seats=%d", err, route.SeatsRemaining)
	}
}
