// Ulov24 yo'lovchilarni o'sha tomonga ketayotgan haydovchilar bilan bog'laydi.
package main

import (
	"fmt"
	"os"
	"strconv"
)

// Version — dasturning qaysi bosqichda ekani (git tag bilan bir xil).
const Version = "c1-m3"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Ishlatish: ulov24 <buyruq>")
		fmt.Println("Buyruqlar: version, modes, price <son>, routes, match")
		return
	}

	switch os.Args[1] {
	case "routes":
		current := now()
		routes, _ := sampleData(current)
		for _, r := range routes {
			label := ""
			switch r.TripMode {
			case TripModeNow:
				label = "hozir"
			case TripModeToday:
				label = "bugun"
			case TripModePlanned:
				label = "rejali"
			}
			if current.After(r.DepartureAt) {
				label = "eskirgan"
			}
			fmt.Printf("#%d %s → %s · %s · %d joy · %s · %s\n",
				r.ID, cityName(r.FromLocationID), cityName(r.ToLocationID),
				r.DepartureAt.In(Tashkent).Format("02.01 15:04"),
				r.SeatsRemaining, FormatSum(r.PricePerSeat), label)
		}
	case "match":
		current := now()
		routes, requests := sampleData(current)
		for _, q := range requests {
			found := ""
			for _, r := range routes {
				if r.FromLocationID != q.FromLocationID || r.ToLocationID != q.ToLocationID {
					continue // boshqa yo'nalish
				}
				if r.Status != DriverRouteOpen || current.After(r.DepartureAt) {
					continue // yopilgan yoki eskirgan
				}
				if r.DepartureAt.Before(q.DesiredDepartureFrom) || r.DepartureAt.After(q.DesiredDepartureTo) {
					continue // vaqt to'g'ri kelmaydi
				}
				if r.SeatsRemaining < q.PassengerCount {
					continue // joy yetmaydi
				}
				if r.DriverUserID == q.PassengerUserID {
					continue // o'zi bilan o'zi
				}
				if found != "" {
					found += ", "
				}
				found += fmt.Sprintf("#%d", r.ID)
			}
			if found == "" {
				fmt.Printf("So'rov #%d (%d kishi): mos yo'nalish yo'q\n", q.ID, q.PassengerCount)
			} else {
				fmt.Printf("So'rov #%d (%d kishi): yo'nalish %s\n", q.ID, q.PassengerCount, found)
			}
		}
	case "version":
		fmt.Println("ulov24", Version)
	case "modes":
		fmt.Println(TripModeNow, "—", TripModeNow.Label())
		fmt.Println(TripModeToday, "—", TripModeToday.Label())
		fmt.Println(TripModePlanned, "—", TripModePlanned.Label())
	case "price":
		if len(os.Args) < 3 {
			fmt.Println("Narxni yozing: ulov24 price 150000")
			os.Exit(1)
		}
		amount, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Narx butun son bo'lishi kerak:", os.Args[2])
			os.Exit(1)
		}
		if amount < 0 {
			fmt.Println("Narx manfiy bo'lishi mumkin emas:", amount)
			os.Exit(1)
		}
		fmt.Println(FormatSum(amount))
	default:
		fmt.Println("Noma'lum buyruq:", os.Args[1])
		os.Exit(1)
	}
}
