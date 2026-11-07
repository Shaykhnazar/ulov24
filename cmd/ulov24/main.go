// Ulov24 yo'lovchilarni o'sha tomonga ketayotgan haydovchilar bilan bog'laydi.
package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Shaykhnazar/ulov24/internal/ridematching"
	"github.com/Shaykhnazar/ulov24/internal/rides"
	"github.com/Shaykhnazar/ulov24/internal/store"
)

// Version — dasturning qaysi bosqichda ekani (git tag bilan bir xil).
const Version = "c1-m7"

const usage = `Ishlatish: ulov24 <buyruq>
Buyruqlar:
  version | modes | cities | price <son>
  route add <qayerdan> <qayerga> "<YYYY-MM-DD HH:MM>" <joylar> <rejim> [narx]
  request add <qayerdan> <qayerga> "<dan>" "<gacha>" <kishi> <rejim>
  feed [qayerdan qayerga]
  match
  accept <moslik> driver|passenger`

func main() {
	if len(os.Args) < 2 {
		fmt.Println(usage)
		return
	}
	args := os.Args[2:]

	switch os.Args[1] {
	case "version":
		fmt.Println("ulov24", Version)
	case "modes":
		for _, m := range []rides.TripMode{rides.TripModeNow, rides.TripModeToday, rides.TripModePlanned} {
			fmt.Println(m, "—", m.Label())
		}
	case "cities":
		for id := rides.CityTashkent; id <= rides.CityNamangan; id++ {
			fmt.Println(id, rides.CityNames[id])
		}
	case "price":
		if len(args) < 1 {
			fail("Narxni yozing: ulov24 price 150000")
		}
		fmt.Println(rides.FormatSum(parsePrice(args[0])))
	case "route":
		if len(args) < 6 || args[0] != "add" {
			fail(usage)
		}
		addRoute(args[1:])
	case "request":
		if len(args) < 7 || args[0] != "add" {
			fail(usage)
		}
		addRequest(args[1:])
	case "feed":
		showFeed(args)
	case "match":
		createMatches()
	case "accept":
		if len(args) < 2 {
			fail(usage)
		}
		accept(args[0], args[1])
	default:
		fail("Noma'lum buyruq: " + os.Args[1])
	}
}

// fail xabarni chiqaradi va dasturni 1 kodi bilan to'xtatadi.
func fail(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func dataPath() string {
	if p := os.Getenv("ULOV24_DATA"); p != "" {
		return p
	}
	return "ulov24.json"
}

func load() store.Snapshot {
	s, err := store.LoadFile(dataPath())
	if errors.Is(err, store.ErrCorruptData) {
		fail("Ma'lumot fayli buzilgan: " + dataPath())
	}
	if err != nil {
		fail("Ma'lumot faylini o'qib bo'lmadi: " + err.Error())
	}
	return s
}

func save(s store.Snapshot) {
	if err := store.SaveFile(dataPath(), s); err != nil {
		fail("Saqlab bo'lmadi: " + err.Error())
	}
}

func now() time.Time {
	t, err := rides.Now()
	if err != nil {
		fail("ULOV24_NOW noto'g'ri: " + os.Getenv("ULOV24_NOW"))
	}
	return t
}

// user — buyruqni kim beryapti. ULOV24_USER berilmasa: haydovchi 1, yo'lovchi 2.
func user(fallback int64) int64 {
	v := os.Getenv("ULOV24_USER")
	if v == "" {
		return fallback
	}
	id, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		fail("ULOV24_USER butun son bo'lishi kerak: " + v)
	}
	return id
}

func parseCity(name string) int64 {
	id, ok := rides.ParseCity(name)
	if !ok {
		fail("Noma'lum shahar: " + name)
	}
	return id
}

func parseDirection(fromName, toName string) (int64, int64) {
	from, to := parseCity(fromName), parseCity(toName)
	if from == to {
		fail("Qayerdan va qayerga bir xil bo'lmasin")
	}
	return from, to
}

func parseTime(s string) time.Time {
	t, err := rides.ParseLocalTime(s)
	if err != nil {
		fail("Vaqt formati: YYYY-MM-DD HH:MM")
	}
	return t
}

func parseCount(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		fail("Son bo'lishi kerak: " + s)
	}
	if err := rides.ValidateSeats(n); err != nil {
		fail(err.Error())
	}
	return n
}

func parsePrice(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		fail("Narx butun son bo'lishi kerak: " + s)
	}
	if n < 0 {
		fail(fmt.Sprintf("Narx manfiy bo'lishi mumkin emas: %d", n))
	}
	return n
}

// addRoute: <qayerdan> <qayerga> "<vaqt>" <joylar> <rejim> [narx]
func addRoute(args []string) {
	current := now()
	from, to := parseDirection(args[0], args[1])
	dep := parseTime(args[2])
	seats := parseCount(args[3])
	mode := rides.TripMode(args[4])
	if err := rides.ValidateTripMode(mode, dep, current); err != nil {
		fail(err.Error())
	}
	var price *int
	if len(args) > 5 {
		p := parsePrice(args[5])
		price = &p
	}

	s := load()
	r := rides.DriverRoute{
		ID: nextRouteID(s), DriverUserID: user(1), FromLocationID: from, ToLocationID: to,
		DepartureAt: dep, SeatsTotal: seats, SeatsRemaining: seats, PricePerSeat: price,
		TripMode: mode, Status: rides.DriverRouteOpen,
	}
	s.Routes = append(s.Routes, r)
	save(s)
	fmt.Printf("Yo'nalish #%d qo'shildi\n", r.ID)
}

// addRequest: <qayerdan> <qayerga> "<dan>" "<gacha>" <kishi> <rejim>
func addRequest(args []string) {
	current := now()
	from, to := parseDirection(args[0], args[1])
	windowFrom, windowTo := parseTime(args[2]), parseTime(args[3])
	if windowFrom.After(windowTo) {
		fail("Vaqt oralig'i noto'g'ri: boshlanishi oxiridan keyin")
	}
	count := parseCount(args[4])
	mode := rides.TripMode(args[5])
	if err := rides.ValidateTripMode(mode, windowTo, current); err != nil {
		fail(err.Error())
	}

	s := load()
	q := rides.PassengerRequest{
		ID: nextRequestID(s), PassengerUserID: user(2), FromLocationID: from, ToLocationID: to,
		DesiredDepartureFrom: windowFrom, DesiredDepartureTo: windowTo, PassengerCount: count,
		TripMode: mode, Status: rides.PassengerRequestOpen,
	}
	s.Requests = append(s.Requests, q)
	save(s)
	fmt.Printf("So'rov #%d qo'shildi\n", q.ID)
}

func showFeed(args []string) {
	current := now()
	s := load()
	routes := s.Routes
	if len(args) >= 2 {
		from, to := parseDirection(args[0], args[1])
		routes = ridematching.IndexByDirection(routes)[ridematching.Direction{From: from, To: to}]
	}
	for _, r := range ridematching.Feed(routes, current) {
		fmt.Println(routeLine(r, r.TripMode.Label()))
	}
}

// createMatches har bir ochiq so'rov uchun mos yo'nalishlarni topadi va hali
// faol mosligi yo'q juftliklar uchun yangi moslik yaratadi.
func createMatches() {
	current := now()
	s := load()
	byDirection := ridematching.IndexByDirection(s.Routes)
	created := 0
	for i := range s.Requests {
		q := &s.Requests[i]
		if q.Status != rides.PassengerRequestOpen && q.Status != rides.PassengerRequestMatched {
			continue
		}
		for _, r := range byDirection[ridematching.Direction{From: q.FromLocationID, To: q.ToLocationID}] {
			if ridematching.Compatible(r, *q, current) != nil || hasActiveMatch(s, r.ID, q.ID) {
				continue
			}
			m := ridematching.RideMatch{
				RideMatchState: ridematching.RideMatchState{
					ID: nextMatchID(s), DriverStatus: ridematching.DualStatusPending,
					PassengerStatus: ridematching.DualStatusPending, OverallStatus: ridematching.MatchPending,
				},
				DriverRouteID: r.ID, PassengerRequestID: q.ID, PassengerCount: q.PassengerCount,
			}
			s.Matches = append(s.Matches, m)
			q.Status = rides.PassengerRequestMatched
			created++
			fmt.Printf("Moslik #%d: so'rov #%d → yo'nalish #%d\n", m.ID, q.ID, r.ID)
		}
	}
	if created == 0 {
		fmt.Println("Yangi moslik yo'q")
		return
	}
	save(s)
}

func accept(matchArg, roleArg string) {
	role := ridematching.UserRole(roleArg)
	if role != ridematching.RoleDriver && role != ridematching.RolePassenger {
		fail("Rol driver yoki passenger bo'lishi kerak: " + roleArg)
	}
	id, err := strconv.ParseInt(matchArg, 10, 64)
	if err != nil {
		fail("Moslik raqami son bo'lishi kerak: " + matchArg)
	}

	s := load()
	m := findMatch(s, id)
	if m == nil {
		fail(fmt.Sprintf("Bunday moslik yo'q: %d", id))
	}
	route := findRoute(s, m.DriverRouteID)
	if route == nil {
		fail(fmt.Sprintf("Moslikdagi yo'nalish topilmadi: %d", m.DriverRouteID))
	}

	err = m.Accept(role, route)
	switch {
	case errors.Is(err, ridematching.ErrTerminalRideMatch):
		fail(fmt.Sprintf("Moslik #%d allaqachon yopilgan", m.ID))
	case errors.Is(err, ridematching.ErrRideCapacityFull):
		save(s) // moslik closed_capacity_full holatida saqlanadi
		fail(fmt.Sprintf("Moslik #%d: %s · joy yetmadi", m.ID, m.OverallStatus))
	case err != nil:
		fail("Xato: " + err.Error())
	}

	line := fmt.Sprintf("Moslik #%d: %s", m.ID, m.OverallStatus)
	if m.OverallStatus == ridematching.MatchAccepted {
		if q := findRequest(s, m.PassengerRequestID); q != nil {
			q.Status = rides.PassengerRequestAccepted
		}
		line += fmt.Sprintf(" · yo'nalishda %d joy qoldi", route.SeatsRemaining)
	}
	save(s)
	fmt.Println(line)
}
